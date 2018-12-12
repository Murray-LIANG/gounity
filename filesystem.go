package gounity

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	typeNameFilesystem   = "filesystem"
	typeFieldsFilesystem = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
	}, ",")
)

// Filesystem defines Unity corresponding `Filesystem` type.
type Filesystem struct {
	Resource
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
}

// NfsShareDefaultAccessEnum defines Unity corresponding `NFSShareDefaultAccessEnum`
// enumeration.
type NfsShareDefaultAccessEnum int

const (
	// NoAccess defines `NoAccess` value of NfsShareDefaultAccessEnum.
	NoAccess NfsShareDefaultAccessEnum = iota

	// ReadOnly defines `ReadOnly` value of NfsShareDefaultAccessEnum.
	ReadOnly

	// ReadWrite defines `ReadWrite` value of NfsShareDefaultAccessEnum.
	ReadWrite

	// Root defines `Root` value of NfsShareDefaultAccessEnum.
	Root
)

//go:generate ./gen_resource.sh resource_tmpl.go filesystem_gen.go Filesystem

// CreateNfsShare exports the nfs share from this filesystem.
func (fs *Filesystem) CreateNfsShare(
	name string, defaultAccess NfsShareDefaultAccessEnum,
) (*NfsShare, error) {

	shareParams := map[string]interface{}{
		"defaultAccess": defaultAccess,
	}
	shareCreate := map[string]interface{}{
		"path":               "/",
		"name":               name,
		"nfsShareParameters": shareParams,
	}
	body := map[string]interface{}{
		"nfsShareCreate": []interface{}{shareCreate},
	}
	logger := log.WithField("requestBody", body)
	logger.Debug("creating nfs share")

	var createdId string
	var err error
	if createdId, err = fs.unity.postOnType(typeStorageResource, actionModifyFilesystem,
		body); err != nil {

		logger.WithError(err).Error("failed to create nfs share")
		return nil, err
	}

	logger.WithField("createdNfsShareId", createdId).Debug("nfs share created")

	created, err := fs.unity.GetNfsShareById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created nfs share")
	}
	return created, err
}
