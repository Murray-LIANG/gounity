package gounity

import (
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// CreateNfsShare exports the nfs share from this filesystem.
func (fs *Filesystem) CreateNfsShare(
	name string, defaultAccess NFSShareDefaultAccessEnum,
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

	fields := map[string]interface{}{
		"requestBody": body,
	}

	logger := log.WithFields(fields)
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
