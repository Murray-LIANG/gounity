package gounity

import (
	"encoding/json"
)

// NewNfsShareById constructs a `NfsShare` object with id.
func (u *Unity) NewNfsShareById(id string) *NfsShare {
	return &NfsShare{
		Resource: Resource{
			typeName: typeNameNfsShare, typeFields: typeFieldsNfsShare, unity: u,
		},
		Id: id,
	}
}

// NewNfsShareByName constructs a `NfsShare` object with name.
func (u *Unity) NewNfsShareByName(name string) *NfsShare {
	return &NfsShare{
		Resource: Resource{
			typeName: typeNameNfsShare, typeFields: typeFieldsNfsShare, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *NfsShare) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
	}

	var (
		latest *NfsShare
		err    error
	)
	switch r.Id {
	case "":
		if latest, err = r.unity.GetNfsShareByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetNfsShareById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetNfsShareById retrives the `NfsShare` by given its id.
func (u *Unity) GetNfsShareById(id string) (*NfsShare, error) {
	res := u.NewNfsShareById(id)

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNfsShareByName retrives the `NfsShare` by given its name.
func (u *Unity) GetNfsShareByName(name string) (*NfsShare, error) {
	res := u.NewNfsShareByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNfsShares retrives all `NfsShare` objects.
func (u *Unity) GetNfsShares() ([]*NfsShare, error) {

	return u.FilterNfsShares(nil)
}

// FilterNfsShares filters the `NfsShare` objects by given filters.
func (u *Unity) FilterNfsShares(filter *filter) ([]*NfsShare, error) {
	respEntries, err := u.getCollection(typeNameNfsShare, typeFieldsNfsShare, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillNfsShares(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillNfsShares(respEntries []*instanceResp) ([]*NfsShare, error) {
	resSlice := []*NfsShare{}
	for _, entry := range respEntries {
		res := u.NewNfsShareById("") // empty id for fake `NfsShare` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `NfsShare` object using its id.
func (r *NfsShare) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
