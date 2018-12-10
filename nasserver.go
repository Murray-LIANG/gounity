package gounity

import (
	"strings"
)

var (
	fieldsNasServer = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
	}, ",")
)

// GetNasServerById retrives the nas server by given its Id.
func (u *Unity) GetNasServerById(id string) (*NasServer, error) {
	res := &NasServer{}
	if err := u.getInstanceById("nasServer", id, fieldsNasServer, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNasServerByName retrives the nas server by given its name.
func (u *Unity) GetNasServerByName(name string) (*NasServer, error) {
	res := &NasServer{}
	if err := u.getInstanceByName("nasServer", name, fieldsNasServer, res); err != nil {
		return nil, err
	}
	return res, nil
}