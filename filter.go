package gounity

import (
	"fmt"
	"strings"
)

type filter []string

// NewFilter returns a filter used for filtering collection.
func NewFilter(f string) *filter {
	return &filter{f}
}

// NewFilterf returns a filter used for filtering collection.
func NewFilterf(format string, args ...interface{}) *filter {
	return NewFilter(fmt.Sprintf(format, args...))
}

func (f *filter) And(andFilter string) *filter {
	newFilter := append(*f, "and")
	newFilter = append(newFilter, andFilter)
	return &newFilter
}

func (f *filter) Andf(format string, args ...interface{}) *filter {
	return f.And(fmt.Sprintf(format, args...))
}

func (f *filter) String() string {
	return strings.Join(*f, " ")
}
