package gqlcheck

import (
	"encoding/base64"
)

// WithBasicAuth is an alias to set basic auth in the request header.
func (tt *Tester) WithBasicAuth(user, pass string) *Tester {
	auth := user + ":" + pass
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	return tt.WithHeader("Authorization", "Basic "+encodedAuth)
}

// WithBearerAuth is an alias to set bearer auth in the request header.
func (tt *Tester) WithBearerAuth(token string) *Tester {
	return tt.WithHeader("Authorization", "Bearer "+token)
}
