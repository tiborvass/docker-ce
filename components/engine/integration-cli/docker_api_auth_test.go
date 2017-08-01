package main

import (
	"net/http"

	"github.com/go-check/check"
	"github.com/moby/moby-core/api/types"
	"github.com/moby/moby-core/integration-cli/checker"
	"github.com/moby/moby-core/integration-cli/request"
)

// Test case for #22244
func (s *DockerSuite) TestAuthAPI(c *check.C) {
	testRequires(c, Network)
	config := types.AuthConfig{
		Username: "no-user",
		Password: "no-password",
	}

	expected := "Get https://registry-1.docker.io/v2/: unauthorized: incorrect username or password"
	status, body, err := request.SockRequest("POST", "/auth", config, daemonHost())
	c.Assert(err, check.IsNil)
	c.Assert(status, check.Equals, http.StatusUnauthorized)
	msg := getErrorMessage(c, body)
	c.Assert(msg, checker.Contains, expected, check.Commentf("Expected: %v, got: %v", expected, msg))
}
