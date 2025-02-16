package codefresh

import (
	"fmt"
	"strings"

)

type graphqlError struct {
	Message string
	Locations [] struct {
		Line int
		Column int
	}
	Extensions struct {
		Code string
		Exception struct {
			Stacktrace []string
		}
	}
}

type graphqlErrorResponse struct {
	errors []graphqlError 
	concatenatedErrors string
}
	

func (e graphqlErrorResponse) Error() string {

	if e.concatenatedErrors != "" {
		return e.concatenatedErrors
	}
	var sb strings.Builder
	for _, err := range e.errors {
		sb.WriteString(fmt.Sprintln(err.Message))
	}
	e.concatenatedErrors = sb.String()
	return e.concatenatedErrors
}