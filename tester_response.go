package gqlcheck

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Response is a struct that represents a GraphQL response.
type Response struct {
	Data   map[string]any   `json:"data,omitempty"`
	Errors []map[string]any `json:"errors,omitempty"`
}

// String returns the string representation of the response.
func (r Response) String() string {
	b, _ := json.MarshalIndent(r, "", "  ") //nolint:errchkjson
	return string(b)
}

// Cb set a callback function to evaluate the response.
func (tt *Tester) Cb(callback func(*http.Response)) {
	if !tt.requireCheck() {
		return
	}
	tt.tester.Cb(callback)
}

// Response sets the response to the provided variable.
func (tt *Tester) Response(out any) {
	if !tt.requireCheck() {
		return
	}
	t := tt.tester.T()
	tt.tester.Cb(func(resp *http.Response) {
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.NoError(t, resp.Body.Close())
		require.NoError(t, json.Unmarshal(b, &out))
	})
}

// HasErrors checks if the response contains errors.
func (tt *Tester) HasErrors() *Tester {
	if !tt.requireCheck() {
		return tt
	}
	newTester := *tt // copy
	newTester.tester = tt.tester.MatchesJSONQuery(`.errors`)
	return &newTester
}

// HasNoErrors checks if the response does not contain errors.
func (tt *Tester) HasNoErrors() *Tester {
	if !tt.requireCheck() {
		return tt
	}
	newTester := *tt // copy
	newTester.tester = tt.tester.NotMatchesJSONQuery(`.errors`)
	return &newTester
}

// HasNoError is an alias for HasNoErrors.
func (tt *Tester) HasNoError() *Tester {
	return tt.HasNoErrors()
}

// HasJSON checks if the response body has the provided value.
func (tt *Tester) HasJSON(expected any) *Tester {
	if !tt.requireCheck() {
		return tt
	}
	newTester := *tt // copy
	newTester.tester = tt.tester.HasJSON(expected)
	return &newTester
}

// HasData checks if the response body has the provided data.
func (tt *Tester) HasData(expected any) *Tester {
	return tt.HasJSON(map[string]any{"data": expected})
}

// ContainsString checks if the response body contains the provided string.
func (tt *Tester) ContainsString(s string) *Tester {
	if !tt.requireCheck() {
		return tt
	}
	newTester := *tt // copy
	newTester.tester = tt.tester.ContainsString(s)
	return &newTester
}
