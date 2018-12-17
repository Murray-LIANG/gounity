package gounity

import (
	"strings"
)

var (
	typeNameNfsShare   = "nfsShare"
	typeFieldsNfsShare = strings.Join([]string{
		"description",
		"id",
		"name",
		"exportPaths",
	}, ",")
)

type NfsShareOperator interface {
	genNfsShareOperator
}

// NfsShare defines Unity corresponding `NfsShare` type.
type NfsShare struct {
	Resource
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ExportPaths []string `json:"exportPaths"`
}

//go:generate ./gen_resource.sh resource_tmpl.go nfsshare_gen.go NfsShare
