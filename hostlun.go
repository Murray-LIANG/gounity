package gounity

import (
	"strings"

	"github.com/pkg/errors"
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

type HostLunOperator interface {
	genHostLunOperator

	FilterHostLunByHostAndLun(hostId, lunId string) (*HostLun, error)
}

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
	filter := NewFilterf(`host eq "%s"`, hostId).Andf(`lun eq "%s"`, lunId)

	fileds := map[string]interface{}{
		"hostId": hostId,
		"lunId":  lunId,
		"filter": filter,
	}
	msg := newMessage().withFields(fileds)

	hostLuns, err := u.FilterHostLuns(filter)
	if err != nil {
		return nil, errors.Wrap(err, msg.withMessage("filter hostluns failed").String())
	}
	if len(hostLuns) == 0 {
		return nil, errors.Wrap(
			err, msg.withMessage("filter returned 0 hostlun").String(),
		)
	}
	if len(hostLuns) > 1 {
		return nil, errors.Wrap(
			err, msg.withField("numOfHostLun",
				len(hostLuns)).withMessage(
				"filter returned more than one hostluns").String(),
		)
	}
	return hostLuns[0], nil
}
