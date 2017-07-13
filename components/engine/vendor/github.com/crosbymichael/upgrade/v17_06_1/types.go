//+build ignore

package v17_06_1

import (
	"github.com/containerd/containerd/runtime"
	"github.com/opencontainers/runc/libcontainer"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

//go:generate -command rewrite go run ../gen/rewrite-structs.go

//go:generate rewrite .Process.Capabilities->linuxCapabilities .Linux.Resources.Memory.Swappiness->memorySwappiness .Linux.Seccomp.Syscalls->linuxSyscalls
type Spec specs.Spec

//go:generate rewrite .Capabilities->linuxCapabilities
type ProcessState runtime.ProcessState

//go:generate rewrite .Config.Capabilities->linuxCapabilities .Config.Cgroups.MemorySwappiness->memorySwappiness
type State libcontainer.State
