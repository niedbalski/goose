package client

import (
	"time"

	"gopkg.in/niedbalski/goose.v3/identity"
)

type AuthCleanup func()

func SetAuthenticationTimeout(timeout time.Duration) AuthCleanup {
	origTimeout := authenticationTimeout
	authenticationTimeout = timeout
	return func() {
		authenticationTimeout = origTimeout
	}
}

func SetAuthenticator(client AuthenticatingClient, auth identity.Authenticator) {
	client.(*authenticatingClient).authMode = auth
}
