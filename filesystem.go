package gounity

import (
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func newCreateNfsShareBody(
	fs *Filesystem, name string, opts ...Option,
) map[string]interface{} {
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
	return map[string]interface{}{
		"nfsShareCreate": []interface{}{shareCreate},
	}
}

// CreateNfsShare exports the nfs share from this filesystem.
func (fs *Filesystem) CreateNfsShare(
	name string, opts ...Option,
) (*NfsShare, error) {

	body := newCreateNfsShareBody(fs, name, opts...)

	fields := map[string]interface{}{
		"requestBody": body,
	}

	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error

	log.Debug("creating nfs share")
	if createdId, err = fs.unity.PostOnType(
		typeStorageResource, actionModifyFilesystem, body,
	); err != nil {
		return nil, errors.Wrapf(err, "create nfs share failed: %s", msg)
	}

	log.WithField("createdNfsShareId", createdId).Debug("nfs share created")

	created, err := fs.unity.GetNfsShareById(createdId)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"get created nfs share failed: %s",
			msg.withField("createdNfsShareId", createdId),
		)
	}
	return created, err
}
