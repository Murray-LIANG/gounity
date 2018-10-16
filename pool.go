package gounity

import (
	"context"
	"errors"
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

// GetPoolByID retrives the pool by given its Id.
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

	createdID := resp.Content.StorageResource.Id
	logger.WithField("createdLunID", createdID).Debug("lun created")

	createdLUN, err := p.Unity.GetLunById(createdID)
	if err != nil {
		logger.WithError(err).Error("failed to get the created lun")
	}
	return createdLUN, err
}

func (u *Unity) GetPoolById(id string) (*Pool, error) {
	return nil, errors.New("Not Implemented.")
}

func (u *Unity) GetPoolByName(id string) (*Pool, error) {
	return nil, errors.New("Not Implemented.")
}
