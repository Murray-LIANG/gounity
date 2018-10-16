package gounity

import (
	"fmt"
	"strings"
)

const (
	// UnityResourceNotFoundErrorCode is error code for resource not found.
	UnityResourceNotFoundErrorCode = 131149829
	// UnityLunNameExistErrorCode is error code for Lun name existing.
	UnityLunNameExistErrorCode = 108007744
)

type gounityError struct {
	message    string
	fields     map[string]interface{}
	innerError error
}

func (e *gounityError) Error() string {
	if e.fields == nil {
		e.fields = map[string]interface{}{}
	}
	if e.innerError != nil {
		e.fields["innerError"] = e.innerError
	}
	kvStrs := []string{}
	for k, v := range e.fields {
		kvStrs = append(kvStrs, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("%s %s", e.message, strings.Join(kvStrs, ","))
}

func (e *gounityError) withError(err error) *gounityError {
	e.innerError = err
	return e
}

func (e *gounityError) withField(key string, value interface{}) *gounityError {
	if e.fields == nil {
		e.fields = map[string]interface{}{}
	}
	e.fields[key] = value
	return e
}

func (e *gounityError) withFields(fields map[string]interface{}) *gounityError {
	for k, v := range fields {
		e = e.withField(k, v)
	}
	return e
}

func newGounityError(message string) *gounityError {
	return &gounityError{message: message}
}
