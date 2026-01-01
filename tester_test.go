package gqlcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockTestingT is a mock implementation of TestingT for testing error handling
type mockTestingT struct {
	errorfCalled  bool
	failNowCalled bool
	errorMsg      string
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errorfCalled = true
	m.errorMsg = format
}

func (m *mockTestingT) FailNow() {
	m.failNowCalled = true
	// Note: We don't panic here since requireCheck() no longer calls FailNow()
}

func TestRequireCheck_HasStatusOK(t *testing.T) {
	// Setup a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"todos":[]}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := NewExternal(server.URL)
	mockT := &mockTestingT{}

	// Call HasStatusOK() without calling Check() first
	checker.Test(mockT).
		Query(`query {todos {text}}`).
		HasStatusOK()

	// Verify that Errorf was called (FailNow is not called anymore)
	if !mockT.errorfCalled {
		t.Error("Expected Errorf to be called, but it wasn't")
	}
	if mockT.errorMsg != "Check() must be called before using assertion methods" {
		t.Errorf("Expected error message 'Check() must be called before using assertion methods', got '%s'", mockT.errorMsg)
	}
}

func TestRequireCheck_HasNoError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"todos":[]}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := NewExternal(server.URL)
	mockT := &mockTestingT{}

	// Call HasNoError() without calling Check() first
	checker.Test(mockT).
		Query(`query {todos {text}}`).
		HasNoError()

	// Verify that Errorf was called
	if !mockT.errorfCalled {
		t.Error("Expected Errorf to be called, but it wasn't")
	}
}

func TestRequireCheck_Cb(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"todos":[]}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := NewExternal(server.URL)
	mockT := &mockTestingT{}

	// Call Cb() without calling Check() first
	checker.Test(mockT).
		Query(`query {todos {text}}`).
		Cb(func(resp *http.Response) {
			// This should never be called
			t.Error("Callback should not be executed")
		})

	// Verify that Errorf was called
	if !mockT.errorfCalled {
		t.Error("Expected Errorf to be called, but it wasn't")
	}
}

func TestCheck_ProperUsage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"todos":[]}}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	checker := NewExternal(server.URL)

	// Proper usage: Check() is called before assertion methods
	checker.Test(t).
		WithHeader("Content-Type", "application/json").
		Query(`query {todos {text}}`).
		Check().
		HasStatusOK().
		HasNoError().
		HasData(map[string]any{
			"todos": []any{},
		})
}
