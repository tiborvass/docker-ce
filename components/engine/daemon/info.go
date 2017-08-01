package daemon

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-connections/sockets"
	"github.com/moby/moby-core/api"
	"github.com/moby/moby-core/api/types"
	"github.com/moby/moby-core/cli/debug"
	"github.com/moby/moby-core/daemon/logger"
	"github.com/moby/moby-core/dockerversion"
	"github.com/moby/moby-core/pkg/fileutils"
	"github.com/moby/moby-core/pkg/parsers/kernel"
	"github.com/moby/moby-core/pkg/parsers/operatingsystem"
	"github.com/moby/moby-core/pkg/platform"
	"github.com/moby/moby-core/pkg/sysinfo"
	"github.com/moby/moby-core/pkg/system"
	"github.com/moby/moby-core/registry"
	"github.com/moby/moby-core/volume/drivers"
)

// SystemInfo returns information about the host server the daemon is running on.
func (daemon *Daemon) SystemInfo() (*types.Info, error) {
	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}

	operatingSystem := "<unknown>"
	if s, err := operatingsystem.GetOperatingSystem(); err != nil {
		logrus.Warnf("Could not get operating system name: %v", err)
	} else {
		operatingSystem = s
	}

	// Don't do containerized check on Windows
	if runtime.GOOS != "windows" {
		if inContainer, err := operatingsystem.IsContainerized(); err != nil {
			logrus.Errorf("Could not determine if daemon is containerized: %v", err)
			operatingSystem += " (error determining if containerized)"
		} else if inContainer {
			operatingSystem += " (containerized)"
		}
	}

	meminfo, err := system.ReadMemInfo()
	if err != nil {
		logrus.Errorf("Could not read system memory info: %v", err)
		meminfo = &system.MemInfo{}
	}

	sysInfo := sysinfo.New(true)
	cRunning, cPaused, cStopped := stateCtr.get()

	securityOptions := []string{}
	if sysInfo.AppArmor {
		securityOptions = append(securityOptions, "name=apparmor")
	}
	if sysInfo.Seccomp && supportsSeccomp {
		profile := daemon.seccompProfilePath
		if profile == "" {
			profile = "default"
		}
		securityOptions = append(securityOptions, fmt.Sprintf("name=seccomp,profile=%s", profile))
	}
	if selinuxEnabled() {
		securityOptions = append(securityOptions, "name=selinux")
	}
	rootIDs := daemon.idMappings.RootPair()
	if rootIDs.UID != 0 || rootIDs.GID != 0 {
		securityOptions = append(securityOptions, "name=userns")
	}

	imageCount := 0
	drivers := ""
	for p, ds := range daemon.stores {
		imageCount += len(ds.imageStore.Map())
		drivers += daemon.GraphDriverName(p)
		if len(daemon.stores) > 1 {
			drivers += fmt.Sprintf(" (%s) ", p)
		}
	}

	// TODO @jhowardmsft LCOW support. For now, hard-code the platform shown for the driver status
	p := runtime.GOOS
	if system.LCOWSupported() {
		p = "linux"
	}

	drivers = strings.TrimSpace(drivers)
	v := &types.Info{
		ID:                 daemon.ID,
		Containers:         int(cRunning + cPaused + cStopped),
		ContainersRunning:  int(cRunning),
		ContainersPaused:   int(cPaused),
		ContainersStopped:  int(cStopped),
		Images:             imageCount,
		Driver:             drivers,
		DriverStatus:       daemon.stores[p].layerStore.DriverStatus(),
		Plugins:            daemon.showPluginsInfo(),
		IPv4Forwarding:     !sysInfo.IPv4ForwardingDisabled,
		BridgeNfIptables:   !sysInfo.BridgeNFCallIPTablesDisabled,
		BridgeNfIP6tables:  !sysInfo.BridgeNFCallIP6TablesDisabled,
		Debug:              debug.IsEnabled(),
		NFd:                fileutils.GetTotalUsedFds(),
		NGoroutines:        runtime.NumGoroutine(),
		SystemTime:         time.Now().Format(time.RFC3339Nano),
		LoggingDriver:      daemon.defaultLogConfig.Type,
		CgroupDriver:       daemon.getCgroupDriver(),
		NEventsListener:    daemon.EventsService.SubscribersCount(),
		KernelVersion:      kernelVersion,
		OperatingSystem:    operatingSystem,
		IndexServerAddress: registry.IndexServer,
		OSType:             platform.OSType,
		Architecture:       platform.Architecture,
		RegistryConfig:     daemon.RegistryService.ServiceConfig(),
		NCPU:               sysinfo.NumCPU(),
		MemTotal:           meminfo.MemTotal,
		DockerRootDir:      daemon.configStore.Root,
		Labels:             daemon.configStore.Labels,
		ExperimentalBuild:  daemon.configStore.Experimental,
		ServerVersion:      dockerversion.Version,
		ClusterStore:       daemon.configStore.ClusterStore,
		ClusterAdvertise:   daemon.configStore.ClusterAdvertise,
		HTTPProxy:          sockets.GetProxyEnv("http_proxy"),
		HTTPSProxy:         sockets.GetProxyEnv("https_proxy"),
		NoProxy:            sockets.GetProxyEnv("no_proxy"),
		LiveRestoreEnabled: daemon.configStore.LiveRestoreEnabled,
		SecurityOptions:    securityOptions,
		Isolation:          daemon.defaultIsolation,
	}

	// Retrieve platform specific info
	daemon.FillPlatformInfo(v, sysInfo)

	hostname := ""
	if hn, err := os.Hostname(); err != nil {
		logrus.Warnf("Could not get hostname: %v", err)
	} else {
		hostname = hn
	}
	v.Name = hostname

	return v, nil
}

// SystemVersion returns version information about the daemon.
func (daemon *Daemon) SystemVersion() types.Version {
	v := types.Version{
		Version:       dockerversion.Version,
		GitCommit:     dockerversion.GitCommit,
		MinAPIVersion: api.MinVersion,
		GoVersion:     runtime.Version(),
		Os:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		BuildTime:     dockerversion.BuildTime,
		Experimental:  daemon.configStore.Experimental,
	}

	kernelVersion := "<unknown>"
	if kv, err := kernel.GetKernelVersion(); err != nil {
		logrus.Warnf("Could not get kernel version: %v", err)
	} else {
		kernelVersion = kv.String()
	}
	v.KernelVersion = kernelVersion

	return v
}

func (daemon *Daemon) showPluginsInfo() types.PluginsInfo {
	var pluginsInfo types.PluginsInfo

	pluginsInfo.Volume = volumedrivers.GetDriverList()
	pluginsInfo.Network = daemon.GetNetworkDriverList()
	// The authorization plugins are returned in the order they are
	// used as they constitute a request/response modification chain.
	pluginsInfo.Authorization = daemon.configStore.AuthorizationPlugins
	pluginsInfo.Log = logger.ListDrivers()

	return pluginsInfo
}
