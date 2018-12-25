package gounity

import (
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

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

	fields := map[string]interface{}{
		"requestBody": body,
	}
	logger := log.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error
	logger.Debug("creating lun")
	if createdId, err = p.unity.PostOnType(
		typeStorageResource, actionCreateLun, body,
	); err != nil {
		return nil, errors.Wrap(err, msg.withMessage("create lun failed").String())
	}

	logger.WithField("createdLunId", createdId).Debug("lun created")

	createdLun, err := p.unity.GetLunById(createdId)
	if err != nil {
		return nil, errors.Wrap(
			err,
			msg.withField("createdLunId",
				createdId).withMessage("get created lun failed").String(),
		)
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

	fields := map[string]interface{}{
		"requestBody": body,
	}
	logger := log.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error
	logger.Debug("creating filesystem")
	if createdId, err = p.unity.PostOnType(
		typeStorageResource, actionCreateFilesystem, body,
	); err != nil {
		return nil, errors.Wrap(err, msg.withMessage("create filesystem failed").String())
	}

	logger.WithField("createdFilesystemId", createdId).Debug("filesystem created")

	createdFs, err := p.unity.GetFilesystemById(createdId)
	if err != nil {
		return nil, errors.Wrap(
			err,
			msg.withField("createdFilesystemId",
				createdId).withMessage("get created filesystem failed").String(),
		)
	}
	return createdFs, err
}
