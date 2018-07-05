package gounity

import (
	"reflect"
	"strings"
)

var (
	fieldsLUN = strings.Join([]string{
		// "compressionPercent",
		// "compressionSizeSaved",
		// "currentNode",
		// "defaultNode",
		"description",
		"health",
		"hostAccess",
		"id",
		// "ioLimitPolicy.id",
		// "isCompressionEnabled",
		// "isReplicationDestination",
		// "isSnapSchedulePaused",
		"isThinEnabled",
		"metadataSize",
		"metadataSizeAllocated",
		"name",
		// "perTierSizeUsed",
		"pool.id",
		"sizeAllocated",
		"sizeTotal",
		"sizeUsed",
		"snapCount",
		// "snapSchedule.id",
		"snapWwn",
		"snapsSize",
		"snapsSizeAllocated",
		// "storageResource.id",
		// "tieringPolicy",
		// "type",
		"wwn",
	}, ",")
)

func (u *Unity) GetLUNByID(id string) (*LUN, error) {
	res := &LUN{}
	if err := u.getInstanceByID("lun", id, fieldsLUN, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) GetLUNs() ([]*LUN, error) {
	collection, err := u.getCollection("lun", fieldsLUN, nil, reflect.TypeOf(LUN{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*LUN)
	return res, nil
}
