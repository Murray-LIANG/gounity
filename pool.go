package gounity

import (
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func newCreateLunBody(
	p *Pool, name string, sizeGB uint64, opts ...Option,
) map[string]interface{} {

	o := NewOptions(opts...)
	defer o.WarnNotUsedOptions()

	body := map[string]interface{}{
		"name": name,
		"lunParameters": map[string]interface{}{
			"pool": p.Repr(),
			"size": gbToBytes(sizeGB),
		},
	}

	if ha := o.PopHostAccess(); ha != nil {
		body["hostAccess"] = ha
	}
	return body
}

// CreateLun creates a new Lun on the pool.
// Parameter - name and sizeGB are required.
// HostAccessOpt is optional.
func (p *Pool) CreateLun(
	name string, sizeGB uint64, opts ...Option,
) (*Lun, error) {

	body := newCreateLunBody(p, name, sizeGB, opts...)

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

func newFsParameters(
	p *Pool, nasServer *NasServer, sizeGB uint64, o *Options,
) map[string]interface{} {
	res := map[string]interface{}{
		"nasServer": nasServer.Repr(),
		"pool":      p.Repr(),
		"size":      gbToBytes(sizeGB),
	}
	if protocol := o.PopSupportedProtocols(); protocol != nil {
		res["supportedProtocols"] = protocol
	}
	return res
}

func newCreateFilesystemBody(
	p *Pool, nasServer *NasServer, name string, sizeGB uint64, opts ...Option,
) map[string]interface{} {

	o := NewOptions(opts...)
	defer o.WarnNotUsedOptions()

	return map[string]interface{}{
		"name":         name,
		"fsParameters": newFsParameters(p, nasServer, sizeGB, o),
	}
}

// CreateFilesystem creates a new filesystem on the pool.
// Parameters - nasServer, name and sizeGB are required.
// SupportedProtocolsOpt is optional.
func (p *Pool) CreateFilesystem(
	nasServer *NasServer, name string, sizeGB uint64, opts ...Option,
) (*Filesystem, error) {

	body := newCreateFilesystemBody(p, nasServer, name, sizeGB, opts...)

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

func newCreateNfsShareBody(
	p *Pool, nasServer *NasServer, name string, sizeGB uint64, opts ...Option,
) map[string]interface{} {

	// Add supportedProtocols to filesystem create.
	opts = append(opts, SupportedProtocolsOpt(FSSupportedProtocolNFS))

	o := NewOptions(opts...)
	defer o.WarnNotUsedOptions()

	shareCreate := map[string]interface{}{
		"path": "/",
		"name": name,
	}

	if da := o.PopDefaultAccess(); da != nil {
		shareCreate["nfsShareParameters"] = map[string]interface{}{
			"defaultAccess": da,
		}
	}

	body := map[string]interface{}{
		"name":           name,
		"fsParameters":   newFsParameters(p, nasServer, sizeGB, o),
		"nfsShareCreate": []interface{}{shareCreate},
	}
	return body
}

// CreateNfsShare creates a new filesystem on the pool then exports a nfs share from it.
// Parameters - nasServer, name and sizeGB are required.
// DefaultAccessOpt is optional.
func (p *Pool) CreateNfsShare(
	nasServer *NasServer, name string, sizeGB uint64, opts ...Option,
) (*NfsShare, error) {

	body := newCreateNfsShareBody(p, nasServer, name, sizeGB, opts...)

	fields := map[string]interface{}{
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("creating filesystem and nfs share")
	// Ignore the returned resource id, use share name to get the nfs share.
	_, err := p.Unity.PostOnType(typeStorageResource, actionCreateFilesystem, body)
	if err != nil {
		return nil, errors.Wrapf(err, "create filesystem and nfs share failed: %s", msg)
	}

	nfs := p.Unity.NewNfsShareByName(name)
	if err = nfs.Refresh(); err != nil {
		return nil, errors.Wrapf(err, "get created nfs share failed: %s", msg)
	}

	log.WithField("createdNfsShareId", nfs.Id).Debug("nfs share created")
	return nfs, err
}
