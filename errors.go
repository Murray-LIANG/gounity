package gounity

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// IsUnityError returns true if err is an unity error.
func IsUnityError(err error) bool {
	_, ok := errors.Cause(err).(*unityError)
	return ok
}

// GetUnityErrorStatusCode returns the unity error code.
func GetUnityErrorStatusCode(err error) int {
	if e, ok := errors.Cause(err).(*unityError); ok {
		return e.HttpStatusCode
	}
	return -1
}

type message struct {
	message string
	fields  map[string]interface{}
}

func (m *message) String() string {
	if m.fields == nil {
		m.fields = map[string]interface{}{}
	}
	strs := []string{}
	if m.message != "" {
		strs = append(strs, m.message)
	}
	for k, v := range m.fields {
		strs = append(strs, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(strs, ",")
}

func (m *message) withField(key string, value interface{}) *message {
	return m.withFields(map[string]interface{}{key: value})
}

func (m *message) withFields(fields map[string]interface{}) *message {
	res := newMessage()
	res.fields = map[string]interface{}{}
	for k, v := range m.fields {
		res.fields[k] = v
	}
	for k, v := range fields {
		res.fields[k] = v
	}
	return res
}

func newMessage() *message {
	return &message{}
}
