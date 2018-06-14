package gounity

import (
	"context"
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
	lunResp := &lunResp{}
	err := u.client.Get(context.Background(), buildInstanceQueryURL("lun", id, fieldsLUN),
		nil, lunResp)
	if err != nil {
		return nil, err
	}

	return &lunResp.Content, nil
}

func (u *Unity) GetLUNs() ([]*LUN, error) {
	lunsResp := &lunsResp{}
	err := u.client.Get(context.Background(),
		buildCollectionQueryURL("lun", fieldsLUN), nil, lunsResp)
	if err != nil {
		return nil, err
	}

	res := []*LUN{}
	for _, lunResp := range lunsResp.Entries {
		res = append(res, &(lunResp.Content))
	}

	return res, nil
}
