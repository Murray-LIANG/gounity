package gounity

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type genFilesystemOperator interface {
	NewFilesystemById(id string) *Filesystem

	NewFilesystemByName(name string) *Filesystem

	GetFilesystemById(id string) (*Filesystem, error)

	GetFilesystemByName(name string) (*Filesystem, error)

	GetFilesystems() ([]*Filesystem, error)

	FillFilesystems(respEntries []*instanceResp) ([]*Filesystem, error)

	FilterFilesystems(filter *filter) ([]*Filesystem, error)
}

// NewFilesystemById constructs a `Filesystem` object with id.
func (u *Unity) NewFilesystemById(id string) *Filesystem {
	return &Filesystem{
		Resource: Resource{
			typeName: typeNameFilesystem, typeFields: typeFieldsFilesystem, unity: u,
		},
		Id: id,
	}
}

// NewFilesystemByName constructs a `Filesystem` object with name.
func (u *Unity) NewFilesystemByName(name string) *Filesystem {
	return &Filesystem{
		Resource: Resource{
			typeName: typeNameFilesystem, typeFields: typeFieldsFilesystem, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *Filesystem) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *Filesystem
		err    error
	)

	switch r.Id {
	case "":
		if latest, err = r.unity.GetFilesystemByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetFilesystemById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetFilesystemById retrives the `Filesystem` by given its id.
func (u *Unity) GetFilesystemById(id string) (*Filesystem, error) {
	res := u.NewFilesystemById(id)

	if err := u.GetInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Filesystem by id failed")
	}
	return res, nil
}

// GetFilesystemByName retrives the `Filesystem` by given its name.
func (u *Unity) GetFilesystemByName(name string) (*Filesystem, error) {
	res := u.NewFilesystemByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Filesystem by name failed")
	}
	return res, nil
}

// GetFilesystems retrives all `Filesystem` objects.
func (u *Unity) GetFilesystems() ([]*Filesystem, error) {

	return u.FilterFilesystems(nil)
}

// FilterFilesystems filters the `Filesystem` objects by given filters.
func (u *Unity) FilterFilesystems(filter *filter) ([]*Filesystem, error) {
	respEntries, err := u.GetCollection(typeNameFilesystem, typeFieldsFilesystem, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter Filesystem failed")
	}
	res, err := u.FillFilesystems(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill Filesystems failed")
	}
	return res, nil
}

// FillFilesystems generates the `Filesystem` objects from collection query response.
func (u *Unity) FillFilesystems(respEntries []*instanceResp) ([]*Filesystem, error) {
	resSlice := []*Filesystem{}
	for _, entry := range respEntries {
		res := u.NewFilesystemById("") // empty id for fake `Filesystem` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Filesystem` object using its id.
func (r *Filesystem) Repr() *idRepresent {
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
