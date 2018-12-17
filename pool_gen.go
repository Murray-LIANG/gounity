package gounity

import (
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
)

type genPoolOperator interface {
	NewPoolById(id string) *Pool

	NewPoolByName(name string) *Pool

	GetPoolById(id string) (*Pool, error)

	GetPoolByName(name string) (*Pool, error)

	GetPools() ([]*Pool, error)

	FillPools(respEntries []*instanceResp) ([]*Pool, error)

	FilterPools(filter *filter) ([]*Pool, error)
}

// NewPoolById constructs a `Pool` object with id.
func (u *Unity) NewPoolById(id string) *Pool {
	return &Pool{
		Resource: Resource{
			typeName: typeNamePool, typeFields: typeFieldsPool, unity: u,
		},
		Id: id,
	}
}

// NewPoolByName constructs a `Pool` object with name.
func (u *Unity) NewPoolByName(name string) *Pool {
	return &Pool{
		Resource: Resource{
			typeName: typeNamePool, typeFields: typeFieldsPool, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *Pool) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *Pool
		err    error
	)
	
	switch r.Id {
	case "":
		if latest, err = r.unity.GetPoolByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetPoolById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetPoolById retrives the `Pool` by given its id.
func (u *Unity) GetPoolById(id string) (*Pool, error) {
	res := u.NewPoolById(id)

	if err := u.GetInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Pool by id failed")
	}
	return res, nil
}

// GetPoolByName retrives the `Pool` by given its name.
func (u *Unity) GetPoolByName(name string) (*Pool, error) {
	res := u.NewPoolByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get Pool by name failed")
	}
	return res, nil
}

// GetPools retrives all `Pool` objects.
func (u *Unity) GetPools() ([]*Pool, error) {

	return u.FilterPools(nil)
}

// FilterPools filters the `Pool` objects by given filters.
func (u *Unity) FilterPools(filter *filter) ([]*Pool, error) {
	respEntries, err := u.GetCollection(typeNamePool, typeFieldsPool, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter Pool failed")
	}
	res, err := u.FillPools(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill Pools failed")
	}
	return res, nil
}

// FillPools generates the `Pool` objects from collection query response.
func (u *Unity) FillPools(respEntries []*instanceResp) ([]*Pool, error) {
	resSlice := []*Pool{}
	for _, entry := range respEntries {
		res := u.NewPoolById("") // empty id for fake `Pool` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Pool` object using its id.
func (r *Pool) Repr() *idRepresent {
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
