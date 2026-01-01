package gqlcheck

import (
	"net/http"

	"github.com/ikawaha/httpcheck"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
}

// Tester represents the GraphQL tester.
type Tester struct {
	checker *Checker
	tester  *httpcheck.Tester

	t           TestingT
	method      string
	path        string            // URL path with query parameters
	body        any               // request body
	contentType string            // Content-Type header
	headers     map[string]string // additional headers
}

// TestOption represents an option for the Test method.
type TestOption func(*Tester)

// WithMethod sets the HTTP method for the request.
func WithMethod(method string) TestOption {
	return func(t *Tester) {
		t.method = method
	}
}

// Test starts a new test with the given *testing.T.
// The default HTTP method is POST. You can specify a different method using WithMethod option.
func (c *Checker) Test(t TestingT, opts ...TestOption) *Tester {
	tester := &Tester{
		checker: c,
		t:       t,
		method:  http.MethodPost, // default is POST
		path:    "",
		headers: map[string]string{},
	}
	for _, opt := range opts {
		opt(tester)
	}

	return tester
}

// requireCheck verifies that Check() has been called before using assertion methods.
// If Check() has not been called, it logs an error and returns false.
func (tt *Tester) requireCheck() bool {
	if tt.tester == nil {
		tt.t.Errorf("Check() must be called before using assertion methods")
		return false
	}
	return true
}

// Check makes request to built request object.
// After request is made, it saves response object for future assertions.
func (tt *Tester) Check() *Tester {
	client := tt.checker.client.Test(tt.t, tt.method, tt.path)

	if tt.contentType != "" && tt.headers["Content-Type"] == "" {
		client = client.WithHeader("Content-Type", tt.contentType)
	}
	for k, v := range tt.headers {
		client = client.WithHeader(k, v)
	}

	if tt.body != nil {
		switch v := tt.body.(type) {
		case string:
			client = client.WithString(v)
		default:
			client = client.WithJSON(v)
		}
	}

	client = client.Check()

	return &Tester{
		checker:     tt.checker,
		t:           tt.t,
		method:      tt.method,
		path:        tt.path,
		body:        tt.body,
		contentType: tt.contentType,
		headers:     tt.headers,
		tester:      client,
	}
}
