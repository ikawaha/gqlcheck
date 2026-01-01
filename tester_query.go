package gqlcheck

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Query is a struct to represent a query.
type Query struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// String returns the string representation of the query.
func (q Query) String() string {
	b, _ := json.MarshalIndent(q, "", "  ") // nolint:errchkjson
	return string(b)
}

// Request sets the query and variables to the request.
func (tt *Tester) Request(q Query) *Tester {
	if len(q.Variables) > 0 {
		return tt.QueryWithVariables(q.Query, q.Variables)
	}
	return tt.Query(q.Query)
}

// Query sets the query to the request.
func (tt *Tester) Query(q string) *Tester {
	newTester := *tt // copy

	if tt.method == http.MethodGet {
		// GET: use query parameter, no Content-Type header
		params := url.Values{}
		params.Set("query", q)
		newTester.path = "/?" + params.Encode()
		return &newTester
	}

	// POST: Set Content-Type to application/graphql only if not already set
	if newTester.contentType == "" {
		newTester.contentType = "application/graphql"
	}

	// If Content-Type is application/json, use JSON format
	if newTester.contentType == "application/json" {
		newTester.body = map[string]any{"query": q}
	} else {
		// Otherwise, use query string as body
		newTester.body = q
	}
	return &newTester
}

// QueryWithVariables sets the query and variables to the request.
func (tt *Tester) QueryWithVariables(q string, variables map[string]any) *Tester {
	newTester := *tt // copy

	if tt.method == http.MethodGet {
		// GET: use query parameters, no Content-Type header
		params := url.Values{}
		params.Set("query", q)
		if len(variables) > 0 {
			varsJSON, _ := json.Marshal(variables) // nolint:errchkjson
			params.Set("variables", string(varsJSON))
		}
		newTester.path = "/?" + params.Encode()
		return &newTester
	}

	// POST: Set Content-Type to application/json only if not already set
	// (variables require JSON format)
	if newTester.contentType == "" {
		newTester.contentType = "application/json"
	}

	newTester.body = map[string]any{
		"query":     q,
		"variables": variables,
	}
	return &newTester
}
