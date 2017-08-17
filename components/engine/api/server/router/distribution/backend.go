package distribution

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/moby/moby-core/api/types"
	"golang.org/x/net/context"
)

// Backend is all the methods that need to be implemented
// to provide image specific functionality.
type Backend interface {
	GetRepository(context.Context, reference.Named, *types.AuthConfig) (distribution.Repository, bool, error)
}
