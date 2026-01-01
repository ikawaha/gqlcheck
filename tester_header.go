package gqlcheck

// WithHeader set header in the request.
func (tt *Tester) WithHeader(key, value string) *Tester {
	newTester := *tt // copy
	newTester.headers = make(map[string]string, len(tt.headers)+1)
	for k, v := range tt.headers {
		newTester.headers[k] = v
	}
	newTester.headers[key] = value

	// If Content-Type is set via WithHeader, update contentType field
	if key == "Content-Type" {
		newTester.contentType = value
	}

	return &newTester
}

// WithHeaders sets header in the request.
func (tt *Tester) WithHeaders(headers map[string]string) *Tester {
	newTester := *tt // copy
	newTester.headers = make(map[string]string, len(tt.headers)+len(headers))
	for k, v := range tt.headers {
		newTester.headers[k] = v
	}
	for k, v := range headers {
		newTester.headers[k] = v
	}
	return &newTester
}
