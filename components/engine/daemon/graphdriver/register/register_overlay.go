// +build !exclude_graphdriver_overlay,linux

package register

import (
	// register the overlay graphdriver
	_ "github.com/moby/moby-core/daemon/graphdriver/overlay"
	_ "github.com/moby/moby-core/daemon/graphdriver/overlay2"
)
