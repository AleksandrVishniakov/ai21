package ai21

import "fmt"

type APIError struct {
	Code    int
	Message string
	URL     string
	Method  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s request: %q: %d. %s", e.Method, e.URL, e.Code, e.Message)
}
