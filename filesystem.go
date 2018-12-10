package gounity

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	fieldsFs = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
	}, ",")
)

// GetFilesystemById retrives the filesystem by given its Id.
func (u *Unity) GetFilesystemById(id string) (*Filesystem, error) {
	res := &Filesystem{}
	if err := u.getInstanceById("filesystem", id, fieldsFs, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetFilesystemByName retrives the filesystem by given its name.
func (u *Unity) GetFilesystemByName(name string) (*Filesystem, error) {
	res := &Filesystem{}
	if err := u.getInstanceByName("filesystem", name, fieldsFs, res); err != nil {
		return nil, err
	}
	return res, nil
}

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

	resp := &storageResourceCreateResp{}
	if err := fs.Unity.client.Post(context.Background(),
		postCollectionUrl("storageResource", "modifyFilesystem"),
		nil, body, resp); err != nil {

		logger.WithError(err).Error("failed to create nfs share")
		return nil, err
	}

	createdId := resp.Content.StorageResource.Id
	logger.WithField("createdNfsShareId", createdId).Debug("nfs share created")

	created, err := fs.Unity.GetNfsShareById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created nfs share")
	}
	return created, err
}
