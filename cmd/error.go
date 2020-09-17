package cmd

import (
	"fmt"
	"strings"
)

type Errors []error

func (errList Errors) Error() string {
	if len(errList) < 1 {
		return ""
	}

	out := make([]string, len(errList))
	for i := range errList {
		out[i] = errList[i].Error()
	}
	return strings.Join(out, ", ")
}

// StatusError reports an unsuccessful exit by a command.
type StatusError struct {
	Status     string
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("Status: %s, Code: %d", e.Status, e.StatusCode)
}
