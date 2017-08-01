package main

import (
	"github.com/go-check/check"
	"github.com/moby/moby-core/integration-cli/daemon"
)

func (s *DockerSwarmSuite) getDaemon(c *check.C, nodeID string) *daemon.Swarm {
	s.daemonsLock.Lock()
	defer s.daemonsLock.Unlock()
	for _, d := range s.daemons {
		if d.NodeID == nodeID {
			return d
		}
	}
	c.Fatalf("could not find node with id: %s", nodeID)
	return nil
}

// nodeCmd executes a command on a given node via the normal docker socket
func (s *DockerSwarmSuite) nodeCmd(c *check.C, id string, args ...string) (string, error) {
	return s.getDaemon(c, id).Cmd(args...)
}
