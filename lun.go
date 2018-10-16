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

// GetLunById retrives the Lun by given its Id.
func (u *Unity) GetLunById(id string) (*Lun, error) {
	res := &Lun{}
	if err := u.getInstanceByID("lun", id, fieldsLUN, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetLuns retrives all LUNs.
func (u *Unity) GetLuns() ([]*Lun, error) {
	collection, err := u.getCollection("lun", fieldsLUN, nil, reflect.TypeOf(Lun{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*Lun)
	return res, nil
}
