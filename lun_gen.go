package gounity

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type genLunOperator interface {
	NewLunById(id string) *Lun

	NewLunByName(name string) *Lun

	GetLunById(id string) (*Lun, error)

	GetLunByName(name string) (*Lun, error)

	GetLuns() ([]*Lun, error)

	FillLuns(respEntries []*instanceResp) ([]*Lun, error)

	FilterLuns(filter *filter) ([]*Lun, error)
}

// NewLunById constructs a `Lun` object with id.
func (u *Unity) NewLunById(id string) *Lun {
	return &Lun{
		Resource: Resource{
			typeName: typeNameLun, typeFields: typeFieldsLun, unity: u,
		},
		Id: id,
	}
}

// NewLunByName constructs a `Lun` object with name.
func (u *Unity) NewLunByName(name string) *Lun {
	return &Lun{
		Resource: Resource{
			typeName: typeNameLun, typeFields: typeFieldsLun, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *Lun) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *Lun
		err    error
	)

	switch r.Id {
	case "":
		if latest, err = r.unity.GetLunByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetLunById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetLunById retrives the `Lun` by given its id.
func (u *Unity) GetLunById(id string) (*Lun, error) {
	res := u.NewLunById(id)

	err := u.GetInstanceById(res.typeName, id, res.typeFields, res)
	if err != nil {
		if IsUnityError(err) {
			return nil, err
		}
		return nil, errors.Wrap(err, "get Lun by id failed")
	}
	return res, nil
}

// GetLunByName retrives the `Lun` by given its name.
func (u *Unity) GetLunByName(name string) (*Lun, error) {
	res := u.NewLunByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Lun by name failed")
	}
	return res, nil
}

// GetLuns retrives all `Lun` objects.
func (u *Unity) GetLuns() ([]*Lun, error) {

	return u.FilterLuns(nil)
}

// FilterLuns filters the `Lun` objects by given filters.
func (u *Unity) FilterLuns(filter *filter) ([]*Lun, error) {
	respEntries, err := u.GetCollection(typeNameLun, typeFieldsLun, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter Lun failed")
	}
	res, err := u.FillLuns(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill Luns failed")
	}
	return res, nil
}

// FillLuns generates the `Lun` objects from collection query response.
func (u *Unity) FillLuns(respEntries []*instanceResp) ([]*Lun, error) {
	resSlice := []*Lun{}
	for _, entry := range respEntries {
		res := u.NewLunById("") // empty id for fake `Lun` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Lun` object using its id.
func (r *Lun) Repr() *idRepresent {
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
