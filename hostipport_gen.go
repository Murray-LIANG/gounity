// DO NOT EDIT.
// GENERATED by go:generate at 2019-06-06 09:02:53.356857 +0000 UTC.
package gounity

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// HostIPPort defines `hostIPPort` type.
type HostIPPort struct {
	Resource

	Id      string `json:"id"`
	Address string `json:"address"`
	Host    *Host  `json:"host"`
}

var (
	typeNameHostIPPort   = "hostIPPort"
	typeFieldsHostIPPort = strings.Join([]string{
		"id",
		"address",
		"host",
	}, ",")
)

type HostIPPortOperatorGen interface {
	NewHostIPPortById(id string) *HostIPPort

	GetHostIPPortById(id string) (*HostIPPort, error)

	GetHostIPPorts() ([]*HostIPPort, error)

	FillHostIPPorts(respEntries []*instanceResp) ([]*HostIPPort, error)

	FilterHostIPPorts(filter *filter) ([]*HostIPPort, error)
}

// NewHostIPPortById constructs a `HostIPPort` object with id.
func (u *Unity) NewHostIPPortById(
	id string,
) *HostIPPort {

	return &HostIPPort{
		Resource: Resource{
			typeName: typeNameHostIPPort, typeFields: typeFieldsHostIPPort, Unity: u,
		},
		Id: id,
	}
}

// Refresh updates the info from Unity.
func (r *HostIPPort) Refresh() error {

	if r.Id == "" {
		return fmt.Errorf(
			"cannot refresh on hostIPPort without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *HostIPPort
		err    error
	)

	switch r.Id {

	default:
		if latest, err = r.Unity.GetHostIPPortById(r.Id); err != nil {
			return err
		}
		*r = *latest
	}
	return nil
}

// GetHostIPPortById retrives the `hostIPPort` by given its id.
func (u *Unity) GetHostIPPortById(
	id string,
) (*HostIPPort, error) {

	res := u.NewHostIPPortById(id)
	err := u.GetInstanceById(res.typeName, id, res.typeFields, res)
	if err != nil {
		if IsUnityError(err) {
			return nil, err
		}
		return nil, errors.Wrap(err, "get hostIPPort by id failed")
	}
	return res, nil
}

// GetHostIPPorts retrives all `hostIPPort` objects.
func (u *Unity) GetHostIPPorts() ([]*HostIPPort, error) {

	return u.FilterHostIPPorts(nil)
}

// FilterHostIPPorts filters the `hostIPPort` objects by given filters.
func (u *Unity) FilterHostIPPorts(
	filter *filter,
) ([]*HostIPPort, error) {

	respEntries, err := u.GetCollection(typeNameHostIPPort, typeFieldsHostIPPort, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter hostIPPort failed")
	}
	res, err := u.FillHostIPPorts(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill hostIPPorts failed")
	}
	return res, nil
}

// FillHostIPPorts generates the `hostIPPort` objects from collection query response.
func (u *Unity) FillHostIPPorts(
	respEntries []*instanceResp,
) ([]*HostIPPort, error) {

	resSlice := []*HostIPPort{}
	for _, entry := range respEntries {
		res := u.NewHostIPPortById("") // empty id for fake `HostIPPort` object
		if err := u.unmarshalResource(entry.Content, res); err != nil {
			return nil, errors.Wrap(err, "decode to HostIPPort failed")
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `hostIPPort` object using its id.
func (r *HostIPPort) Repr() *idRepresent {

	log := logrus.WithField("hostIPPort", r)
	if r.Id == "" {
		log.Info("refreshing hostIPPort from unity")
		err := r.Refresh()
		if err != nil {
			log.WithError(err).Error("refresh hostIPPort from unity failed")
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
