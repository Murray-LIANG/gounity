package gounity

import (
	"context"
	"strings"
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
	poolResp := &poolResp{}
	err := u.client.Get(context.Background(),
		buildInstanceQueryURL("pool", id, fieldsPool), nil, poolResp)
	if err != nil {
		return nil, err
	}

	return &poolResp.Content, nil
}
