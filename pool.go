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

func (u *Unity) GetPoolByID(id string) (*Pool, error) {
	res := &Pool{}
	if err := u.getInstanceByID("pool", id, fieldsPool, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) GetPools() ([]*Pool, error) {
	collection, err := u.getCollection("pool", fieldsPool, nil, reflect.TypeOf(Pool{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*Pool)
	return res, nil
}

func (p *Pool) CreateLUN(name string, sizeGB uint64) (*LUN, error) {
	lunParams := map[string]interface{}{
		"pool": represent(p),
		"size": GBToBytes(sizeGB),
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

	if createdLUN, err := p.Unity.GetLUNByID(createdID); err != nil {
		logger.WithError(err).Error("failed to get the created lun")
		return nil, err
	} else {
		return createdLUN, nil
	}

}
