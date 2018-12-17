package gounity

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type genDUMMYOperator interface {
	NewDUMMYById(id string) *DUMMY

	NewDUMMYByName(name string) *DUMMY

	GetDUMMYById(id string) (*DUMMY, error)

	GetDUMMYByName(name string) (*DUMMY, error)

	GetDUMMYs() ([]*DUMMY, error)

	FillDUMMYs(respEntries []*instanceResp) ([]*DUMMY, error)

	FilterDUMMYs(filter *filter) ([]*DUMMY, error)
}

// NewDUMMYById constructs a `DUMMY` object with id.
func (u *Unity) NewDUMMYById(id string) *DUMMY {
	return &DUMMY{
		Resource: Resource{
			typeName: typeNameDUMMY, typeFields: typeFieldsDUMMY, unity: u,
		},
		Id: id,
	}
}

// NewDUMMYByName constructs a `DUMMY` object with name.
func (u *Unity) NewDUMMYByName(name string) *DUMMY {
	return &DUMMY{
		Resource: Resource{
			typeName: typeNameDUMMY, typeFields: typeFieldsDUMMY, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *DUMMY) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
	}

	var (
		latest *DUMMY
		err    error
	)

	switch r.Id {
	case "":
		if latest, err = r.unity.GetDUMMYByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetDUMMYById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetDUMMYById retrives the `DUMMY` by given its id.
func (u *Unity) GetDUMMYById(id string) (*DUMMY, error) {
	res := u.NewDUMMYById(id)

	err := u.GetInstanceById(res.typeName, id, res.typeFields, res)
	if err != nil {
		if IsUnityError(err) {
			return nil, err
		}
		return nil, errors.Wrap(err, "get DUMMY by id failed")
	}
	return res, nil
}

// GetDUMMYByName retrives the `DUMMY` by given its name.
func (u *Unity) GetDUMMYByName(name string) (*DUMMY, error) {
	res := u.NewDUMMYByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get DUMMY by name failed")
	}
	return res, nil
}

// GetDUMMYs retrives all `DUMMY` objects.
func (u *Unity) GetDUMMYs() ([]*DUMMY, error) {

	return u.FilterDUMMYs(nil)
}

// FilterDUMMYs filters the `DUMMY` objects by given filters.
func (u *Unity) FilterDUMMYs(filter *filter) ([]*DUMMY, error) {
	respEntries, err := u.GetCollection(typeNameDUMMY, typeFieldsDUMMY, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter DUMMY failed")
	}
	res, err := u.FillDUMMYs(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill DUMMYs failed")
	}
	return res, nil
}

// FillDUMMYs generates the `DUMMY` objects from collection query response.
func (u *Unity) FillDUMMYs(respEntries []*instanceResp) ([]*DUMMY, error) {
	resSlice := []*DUMMY{}
	for _, entry := range respEntries {
		res := u.NewDUMMYById("") // empty id for fake `DUMMY` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `DUMMY` object using its id.
func (r *DUMMY) Repr() *idRepresent {
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
