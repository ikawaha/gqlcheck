package gqlcheck

import "net/http"

// HasStatus checks if the response status code is the given status code.
func (tt *Tester) HasStatus(status int) *Tester {
	if !tt.requireCheck() {
		return tt
	}
	newTester := *tt // copy
	newTester.tester = tt.tester.HasStatus(status)
	return &newTester
}

// HasStatusOK checks if the response status code is 200.
func (tt *Tester) HasStatusOK() *Tester {
	return tt.HasStatus(http.StatusOK)
}
