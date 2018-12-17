package gounity

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	// UnityResourceNotFoundErrorCode is error code for resource not found.
	UnityResourceNotFoundErrorCode = 131149829
	// UnityLunNameExistErrorCode is error code for Lun name existing.
	UnityLunNameExistErrorCode = 108007744
)

func IsUnityError(err error) bool {
	_, ok := errors.Cause(err).(*unityError)
	return ok
}

func IsUnityResourceNotFoundError(err error) bool {
	e, ok := errors.Cause(err).(*unityError)
	return ok && e.ErrorCode == UnityResourceNotFoundErrorCode
}

func IsUnityLunNameExistError(err error) bool {
	e, ok := errors.Cause(err).(*unityError)
	return ok && e.ErrorCode == UnityLunNameExistErrorCode
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
