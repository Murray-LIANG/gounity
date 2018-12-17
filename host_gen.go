package gounity

import (
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
)

type genHostOperator interface {
	NewHostById(id string) *Host

	NewHostByName(name string) *Host

	GetHostById(id string) (*Host, error)

	GetHostByName(name string) (*Host, error)

	GetHosts() ([]*Host, error)

	FillHosts(respEntries []*instanceResp) ([]*Host, error)

	FilterHosts(filter *filter) ([]*Host, error)
}

// NewHostById constructs a `Host` object with id.
func (u *Unity) NewHostById(id string) *Host {
	return &Host{
		Resource: Resource{
			typeName: typeNameHost, typeFields: typeFieldsHost, unity: u,
		},
		Id: id,
	}
}

// NewHostByName constructs a `Host` object with name.
func (u *Unity) NewHostByName(name string) *Host {
	return &Host{
		Resource: Resource{
			typeName: typeNameHost, typeFields: typeFieldsHost, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *Host) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *Host
		err    error
	)
	
	switch r.Id {
	case "":
		if latest, err = r.unity.GetHostByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetHostById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetHostById retrives the `Host` by given its id.
func (u *Unity) GetHostById(id string) (*Host, error) {
	res := u.NewHostById(id)

	if err := u.GetInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Host by id failed")
	}
	return res, nil
}

// GetHostByName retrives the `Host` by given its name.
func (u *Unity) GetHostByName(name string) (*Host, error) {
	res := u.NewHostByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Host by name failed")
	}
	return res, nil
}

// GetHosts retrives all `Host` objects.
func (u *Unity) GetHosts() ([]*Host, error) {

	return u.FilterHosts(nil)
}

// FilterHosts filters the `Host` objects by given filters.
func (u *Unity) FilterHosts(filter *filter) ([]*Host, error) {
	respEntries, err := u.GetCollection(typeNameHost, typeFieldsHost, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter Host failed")
	}
	res, err := u.FillHosts(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill Hosts failed")
	}
	return res, nil
}

// FillHosts generates the `Host` objects from collection query response.
func (u *Unity) FillHosts(respEntries []*instanceResp) ([]*Host, error) {
	resSlice := []*Host{}
	for _, entry := range respEntries {
		res := u.NewHostById("") // empty id for fake `Host` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Host` object using its id.
func (r *Host) Repr() *idRepresent {
	if r.Id == "" {
		log.Infof("refreshing %v from unity", r)
		err := r.Refresh()
		if err != nil {
			log.Errorf("refresh %v from unity failed, %+v", r, err)
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
