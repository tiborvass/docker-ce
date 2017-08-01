package daemon

import (
	"github.com/moby/moby-core/api/types"
	"github.com/moby/moby-core/api/types/backend"
	"github.com/moby/moby-core/api/types/versions/v1p19"
	"github.com/moby/moby-core/container"
	"github.com/moby/moby-core/daemon/exec"
)

// This sets platform-specific fields
func setPlatformSpecificContainerFields(container *container.Container, contJSONBase *types.ContainerJSONBase) *types.ContainerJSONBase {
	return contJSONBase
}

// containerInspectPre120 get containers for pre 1.20 APIs.
func (daemon *Daemon) containerInspectPre120(name string) (*v1p19.ContainerJSON, error) {
	return &v1p19.ContainerJSON{}, nil
}

func inspectExecProcessConfig(e *exec.Config) *backend.ExecProcessConfig {
	return &backend.ExecProcessConfig{
		Tty:        e.Tty,
		Entrypoint: e.Entrypoint,
		Arguments:  e.Args,
	}
}
