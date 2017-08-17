// +build !windows

package main

import (
	"github.com/go-check/check"
	"github.com/moby/moby-core/integration-cli/checker"
)

func (s *DockerSuite) TestInfoSecurityOptions(c *check.C) {
	testRequires(c, SameHostDaemon, seccompEnabled, Apparmor, DaemonIsLinux)

	out, _ := dockerCmd(c, "info")
	c.Assert(out, checker.Contains, "Security Options:\n apparmor\n seccomp\n  Profile: default\n")
}
