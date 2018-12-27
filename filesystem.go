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

	logger := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	var createdId string
	var err error

	logger.Debug("creating nfs share")
	if createdId, err = fs.unity.PostOnType(
		typeStorageResource, actionModifyFilesystem, body,
	); err != nil {
		return nil, errors.Wrap(err, msg.withMessage("create nfs share failed").String())
	}

	logger.WithField("createdNfsShareId", createdId).Debug("nfs share created")

	created, err := fs.unity.GetNfsShareById(createdId)
	if err != nil {
		return nil, errors.Wrap(
			err,
			msg.withField("createdNfsShareId",
				createdId).withMessage("get created nfs share failed").String(),
		)
	}
	return created, err
}
