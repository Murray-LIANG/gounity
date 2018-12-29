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

type filterOn string

func NewFilterOn(typeName string) filterOn {
	return filterOn(typeName)
}

func (fo filterOn) Eq(id string) *filter {
	return NewFilterf(`%s eq "%s"`, fo, id)
}

func (f *filter) And(otherFilter *filter) *filter {
	newFilter := append(*f, "and")
	newFilter = append(newFilter, (*otherFilter)...)
	return &newFilter
}

func (f *filter) String() string {
	return strings.Join(*f, " ")
}
