package gounity

import (
	"context"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	fieldsPool = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
		"sizeFree",
		"sizeTotal",
		"sizeUsed",
	}, ",")
)

// GetPoolById retrives the pool by given its Id.
func (u *Unity) GetPoolById(id string) (*Pool, error) {
	res := &Pool{}
	if err := u.getInstanceById("pool", id, fieldsPool, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetPoolByName retrives the pool by given its name.
func (u *Unity) GetPoolByName(name string) (*Pool, error) {
	res := &Pool{}
	if err := u.getInstanceByName("pool", name, fieldsPool, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetPools retrives all pools.
func (u *Unity) GetPools() ([]*Pool, error) {
	collection, err := u.getCollection("pool", fieldsPool, nil, reflect.TypeOf(Pool{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*Pool)
	return res, nil
}

// CreateLun creates a new Lun on the pool.
func (p *Pool) CreateLun(name string, sizeGB uint64) (*Lun, error) {
	lunParams := map[string]interface{}{
		"pool": represent(p),
		"size": gbToBytes(sizeGB),
	}
	body := map[string]interface{}{
		"name":          name,
		"lunParameters": lunParams,
	}
	logger := log.WithField("requestBody", body)
	logger.Debug("creating lun")

	resp := &storageResourceCreateResp{}
	if err := p.Unity.client.Post(context.Background(),
		postCollectionUrl("storageResource", "createLun"), nil, body, resp); err != nil {

		logger.WithError(err).Error("failed to create lun")
		return nil, err
	}

	createdId := resp.Content.StorageResource.Id
	logger.WithField("createdLunId", createdId).Debug("lun created")

	createdLun, err := p.Unity.GetLunById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created lun")
	}
	return createdLun, err
}

// CreateFilesystem creates a new filesystem on the pool.
func (p *Pool) CreateFilesystem(
	nas_server *NasServer, name string, sizeGB uint64,
) (*Filesystem, error) {

	fsParams := map[string]interface{}{
		"nasServer": represent(nas_server),
		"pool":      represent(p),
		"size":      gbToBytes(sizeGB),
	}
	body := map[string]interface{}{
		"name":          name,
		"lunParameters": fsParams,
	}
	logger := log.WithField("requestBody", body)
	logger.Debug("creating filesystem")

	resp := &storageResourceCreateResp{}
	if err := p.Unity.client.Post(context.Background(),
		postCollectionUrl("storageResource", "createFilesystem"),
		nil, body, resp); err != nil {

		logger.WithError(err).Error("failed to create filesystem")
		return nil, err
	}

	createdId := resp.Content.StorageResource.Id
	logger.WithField("createdFilesystemId", createdId).Debug("filesystem created")

	createdFs, err := p.Unity.GetFilesystemById(createdId)
	if err != nil {
		logger.WithError(err).Error("failed to get the created filesystem")
	}
	return createdFs, err
}
