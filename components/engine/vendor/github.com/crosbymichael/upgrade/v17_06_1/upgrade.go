package v17_06_1

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func Upgrade(runcState, containerdConfig, containerdProcess string) error {
	newRuncState, err := UpgradeState(runcState)
	if err != nil {
		return err
	}
	newContainerdConfig, err := UpgradeConfig(containerdConfig)
	if err != nil {
		return err
	}
	newContainerdProcess, err := UpgradeProcessState(containerdProcess)
	if err != nil {
		return err
	}
	// Now that all JSON have been successfully rewritten,
	// replace original.
	// Try renaming all, and return error at the end.
	var errs []string
	for _, s := range []struct{ old, new string }{
		{runcState, newRuncState},
		{containerdConfig, newContainerdConfig},
		{containerdProcess, newContainerdProcess},
	} {
		if err := os.Rename(s.old, s.new); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

func UpgradeState(filename string) (string, error) {
	var x State
	return reencode(filename, "upgrade-runc-state", &x)
}

func UpgradeConfig(filename string) (string, error) {
	var x Spec
	return reencode(filename, "upgrade-containerd-config", &x)
}

func UpgradeProcessState(filename string) (string, error) {
	var x ProcessState
	return reencode(filename, "upgrade-containerd-processstate", &x)
}

func reencode(filename, prefix string, x interface{}) (string, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(filename, os.O_RDWR, fi.Mode())
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(x); err != nil {
		return "", err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}
	return filename, json.NewEncoder(f).Encode(x)
}
