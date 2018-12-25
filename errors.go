package gounity

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func IsUnityError(err error) bool {
	_, ok := errors.Cause(err).(*unityError)
	return ok
}

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
	if m.fields == nil {
		m.fields = map[string]interface{}{}
	}
	m.fields[key] = value
	return m
}

func (m *message) withFields(fields map[string]interface{}) *message {
	for k, v := range fields {
		m = m.withField(k, v)
	}
	return m
}

func (m *message) withMessage(msg string) *message {
	m.message = msg
	return m
}

func (m *message) withMessagef(format string, args ...interface{}) *message {
	m.message = fmt.Sprintf(format, args...)
	return m
}

func newMessage() *message {
	return &message{}
}
