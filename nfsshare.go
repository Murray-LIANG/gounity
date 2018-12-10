package gounity

import (
	"strings"
)

var (
	fieldsNfsShare = strings.Join([]string{
		"description",
		"id",
		"name",
		"exportPaths",
	}, ",")
)

// GetNfsShareById retrives the nfs share by given its Id.
func (u *Unity) GetNfsShareById(id string) (*NfsShare, error) {
	res := &NfsShare{}
	if err := u.getInstanceById("nfsShare", id, fieldsNfsShare, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNfsShareByName retrives the nfs share by given its name.
func (u *Unity) GetNfsShareByName(name string) (*NfsShare, error) {
	res := &NfsShare{}
	if err := u.getInstanceByName("nfsShare", name, fieldsNfsShare, res); err != nil {
		return nil, err
	}
	return res, nil
}
