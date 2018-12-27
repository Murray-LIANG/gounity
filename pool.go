package gounity

import (
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func newCreateLunBody(p *Pool, opts ...Option) map[string]interface{} {

	o := NewOptions(opts...)
	defer o.WarnNotUsedOptions()

	body := map[string]interface{}{"lunParameters": o.NewLunParameters(p)}

	if name := o.PopName(); name != nil {
		body["name"] = name
	}

	if ha := o.PopHostAccess(); ha != nil {
		body["hostAccess"] = ha
	}
	return body
}

// CreateLun creates a new Lun on the pool.
func (p *Pool) CreateLun(opts ...Option) (*Lun, error) {
	body := newCreateLunBody(p, opts...)

	fields := map[string]interface{}{
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error
	log.Debug("creating lun")
	if createdId, err = p.unity.PostOnType(
		typeStorageResource, actionCreateLun, body,
	); err != nil {
		return nil, errors.Wrap(err, msg.withMessage("create lun failed").String())
	}

	log.WithField("createdLunId", createdId).Debug("lun created")

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

func newCreateFilesystemBody(
	p *Pool, nasServer *NasServer, opts ...Option,
) map[string]interface{} {

	o := NewOptions(opts...)
	defer o.WarnNotUsedOptions()

	fsParams := map[string]interface{}{
		"nasServer": nasServer.Repr(),
		"pool":      p.Repr(),
	}
	if size := o.PopSize(); size != nil {
		fsParams["size"] = size
	}

	body := map[string]interface{}{"fsParameters": fsParams}
	if name := o.PopName(); name != nil {
		body["name"] = name
	}
	return body
}

// CreateFilesystem creates a new filesystem on the pool.
func (p *Pool) CreateFilesystem(
	nasServer *NasServer, opts ...Option,
) (*Filesystem, error) {

	body := newCreateFilesystemBody(p, nasServer, opts...)

	fields := map[string]interface{}{
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error
	log.Debug("creating filesystem")
	if createdId, err = p.unity.PostOnType(
		typeStorageResource, actionCreateFilesystem, body,
	); err != nil {
		return nil, errors.Wrap(err, msg.withMessage("create filesystem failed").String())
	}

	log.WithField("createdFilesystemId", createdId).Debug("filesystem created")

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
