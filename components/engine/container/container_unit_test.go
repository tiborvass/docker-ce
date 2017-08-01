package container

import (
	"path/filepath"
	"testing"

	"github.com/moby/moby-core/api/types/container"
	swarmtypes "github.com/moby/moby-core/api/types/swarm"
	"github.com/moby/moby-core/pkg/signal"
)

func TestContainerStopSignal(t *testing.T) {
	c := &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{},
		},
	}

	def, err := signal.ParseSignal(signal.DefaultStopSignal)
	if err != nil {
		t.Fatal(err)
	}

	s := c.StopSignal()
	if s != int(def) {
		t.Fatalf("Expected %v, got %v", def, s)
	}

	c = &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{StopSignal: "SIGKILL"},
		},
	}
	s = c.StopSignal()
	if s != 9 {
		t.Fatalf("Expected 9, got %v", s)
	}
}

func TestContainerStopTimeout(t *testing.T) {
	c := &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{},
		},
	}

	s := c.StopTimeout()
	if s != DefaultStopTimeout {
		t.Fatalf("Expected %v, got %v", DefaultStopTimeout, s)
	}

	stopTimeout := 15
	c = &Container{
		CommonContainer: CommonContainer{
			Config: &container.Config{StopTimeout: &stopTimeout},
		},
	}
	s = c.StopSignal()
	if s != 15 {
		t.Fatalf("Expected 15, got %v", s)
	}
}

func TestContainerSecretReferenceDestTarget(t *testing.T) {
	ref := &swarmtypes.SecretReference{
		File: &swarmtypes.SecretReferenceFileTarget{
			Name: "app",
		},
	}

	d := getSecretTargetPath(ref)
	expected := filepath.Join(containerSecretMountPath, "app")
	if d != expected {
		t.Fatalf("expected secret dest %q; received %q", expected, d)
	}
}
