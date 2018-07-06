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

// GetPoolByID retrives the pool by given its ID.
func (u *Unity) GetPoolByID(id string) (*Pool, error) {
	res := &Pool{}
	if err := u.getInstanceByID("pool", id, fieldsPool, res); err != nil {
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

// CreateLUN creates a new LUN on the pool.
func (p *Pool) CreateLUN(name string, sizeGB uint64) (*LUN, error) {
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
		postCollectionURL("storageResource", "createLun"), nil, body, resp); err != nil {

		logger.WithError(err).Error("failed to create lun")
		return nil, err
	}

	createdID := resp.Content.StorageResource.ID
	logger.WithField("createdLunID", createdID).Debug("lun created")

	createdLUN, err := p.Unity.GetLUNByID(createdID)
	if err != nil {
		logger.WithError(err).Error("failed to get the created lun")
	}
	return createdLUN, err
}
