package gounity

import (
	"strings"
)

var (
	typeNameNasServer   = "nasServer"
	typeFieldsNasServer = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
	}, ",")
)

type NasServerOperator interface {
	genNasServerOperator
}

// NasServer defines Unity corresponding `NasServer` type.
type NasServer struct {
	Resource
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
}

//go:generate ./gen_resource.sh resource_tmpl.go nasserver_gen.go NasServer
