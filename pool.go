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
	if createdId, err = p.Unity.PostOnType(
		typeStorageResource, actionCreateLun, body,
	); err != nil {
		return nil, errors.Wrapf(err, "create lun failed: %s", msg)
	}

	log.WithField("createdLunId", createdId).Debug("lun created")

	createdLun, err := p.Unity.GetLunById(createdId)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"get created lun failed: %s", msg.withField("createdLunId", createdId),
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

	log.Debug("creating filesystem")
	resId, err := p.Unity.PostOnType(typeStorageResource, actionCreateFilesystem, body)
	if err != nil {
		return nil, errors.Wrapf(err, "create filesystem failed: %s", msg)
	}

	storageResource, err := p.Unity.GetStorageResourceById(resId)
	if err != nil {
		return nil, errors.Wrapf(
			err, "get created storageResource failed: %s", msg.withField("resId", resId),
		)
	}

	fs := storageResource.Filesystem
	if fs == nil || fs.Id == "" {
		return nil, errors.Errorf(
			"get filesystem from storageResource failed: %s",
			msg.withField("storageResource", storageResource),
		)
	}
	log.WithField("filesystemId", fs.Id).Debug("filesystem created")

	err = fs.Refresh()
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"get created filesystem failed: %s", msg.withField("filesystemId", fs.Id),
		)
	}
	return fs, err
}
