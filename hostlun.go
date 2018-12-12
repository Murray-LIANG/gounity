package gounity

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	typeNameHostLun   = "hostLUN"
	typeFieldsHostLun = strings.Join([]string{
		"hlu",
		"host.id",
		"id",
		"isDefaultSnap",
		"isReadOnly",
		"lun.id",
		"type",
	}, ",")
)

// HostLun defines Unity corresponding `HostLun` type.
type HostLun struct {
	Resource
	Name          string          `json:"-"`
	Id            string          `json:"id"`
	Host          *Host           `json:"host"`
	Type          HostLunTypeEnum `json:"type"`
	Hlu           uint16          `json:"hlu"`
	Lun           *Lun            `json:"lun"`
	IsReadOnly    bool            `json:"isReadOnly"`
	IsDefaultSnap bool            `json:"isDefaultSnap"`
}

// HostLunTypeEnum defines Unity corresponding `HostLunTypeEnum` enumeration.
type HostLunTypeEnum int

const (
	// HostLunTypeUnknown defines `Unknown` value of HostLunTypeEnum.
	HostLunTypeUnknown HostLunTypeEnum = iota

	// HostLunTypeLun defines `Lun` value of HostLunTypeEnum.
	HostLunTypeLun

	// HostLunTypeSnap defines `Snap` value of HostLunTypeEnum.
	HostLunTypeSnap
)

//go:generate ./gen_resource.sh resource_tmpl.go hostlun_gen.go HostLun

// FilterHostLunByHostAndLun filters the `HostLun` by given its host Id and Lun Id.
func (u *Unity) FilterHostLunByHostAndLun(hostId, lunId string) (*HostLun, error) {
	filter := NewFilter(fmt.Sprintf(`host eq "%s"`, hostId)).And(
		fmt.Sprintf(`lun eq "%s"`, lunId))
	hostLuns, err := u.FilterHostLuns(filter)
	if err != nil {
		return nil, err
	}
	if len(hostLuns) == 0 {
		log.WithField("hostId", hostId).WithField("lunId",
			lunId).Info("filter returns 0 hostLun")
		return nil, nil
	}
	if len(hostLuns) > 1 {
		log.WithField("hostId", hostId).WithField("lunId", lunId).WithField("resultCount",
			len(hostLuns)).Info("filter returns more one hostLuns")
	}
	return hostLuns[0], nil
}
