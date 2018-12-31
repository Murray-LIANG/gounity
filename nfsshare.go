package gounity

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *NfsShare) Delete() error {

	body := map[string]interface{}{
		"nfsShareDelete": []interface{}{
			map[string]interface{}{
				"nfsShare": r.Repr(),
			},
		},
	}

	fields := map[string]interface{}{
		"nfsShare":    r,
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	if r.Filesystem == nil || r.Filesystem.StorageResource == nil {
		log.Info("refreshing nfsShare from unity")
		err := r.Refresh()
		if err != nil {
			return errors.Wrapf(err, "refresh nfsShare from unity failed: %s", msg)
		}
	}

	log.Debug("deleting nfs share")
	err := r.Unity.PostOnInstance(
		typeStorageResource, r.Filesystem.StorageResource.Id, actionModifyFilesystem, body,
	)
	if err != nil {
		return errors.Wrapf(err, "delete nfs share failed: %s", msg)
	}
	log.Debug("nfs share deleted")
	return nil
}
