package gounity

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	typeNamePool   = "pool"
	typeFieldsPool = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
		"sizeFree",
		"sizeTotal",
		"sizeUsed",
	}, ",")
)

// Pool defines Unity corresponding `pool` type.
type Pool struct {
	Resource
	Id          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SizeFree    uint64 `json:"sizeFree,omitempty"`
	SizeTotal   uint64 `json:"sizeTotal,omitempty"`
	SizeUsed    uint64 `json:"sizeUsed,omitempty"`
}

// go:generate ./gen_resource.sh resource_tmpl.go pool_gen.go Pool

// CreateLun creates a new Lun on the pool.
func (p *Pool) CreateLun(name string, sizeGB uint64) (*Lun, error) {
	lunParams := map[string]interface{}{
		"pool": p.Repr(),
		"size": gbToBytes(sizeGB),
	}
	body := map[string]interface{}{
		"name":          name,
		"lunParameters": lunParams,
	}
	logger := log.WithField("requestBody", body)
	logger.Debug("creating lun")

	var createdId string
	var err error
	if createdId, err = p.unity.postOnType(typeStorageResource, actionCreateLun, body); err != nil {

		logger.WithError(err).Error("failed to create lun")
		return nil, err
	}

	logger.WithField("createdLunId", createdId).Debug("lun created")

	createdLun, err := p.unity.GetLunById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created lun")
	}
	return createdLun, err
}

// CreateFilesystem creates a new filesystem on the pool.
func (p *Pool) CreateFilesystem(
	nasServer *NasServer, name string, sizeGB uint64,
) (*Filesystem, error) {

	fsParams := map[string]interface{}{
		"nasServer": nasServer.Repr(),
		"pool":      p.Repr(),
		"size":      gbToBytes(sizeGB),
	}
	body := map[string]interface{}{
		"name":          name,
		"lunParameters": fsParams,
	}
	logger := log.WithField("requestBody", body)
	logger.Debug("creating filesystem")

	var createdId string
	var err error
	if createdId, err = p.unity.postOnType(typeStorageResource, actionCreateFilesystem,
		body); err != nil {

		logger.WithError(err).Error("failed to create filesystem")
		return nil, err
	}

	logger.WithField("createdFilesystemId", createdId).Debug("filesystem created")

	createdFs, err := p.unity.GetFilesystemById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created filesystem")
	}
	return createdFs, err
}
