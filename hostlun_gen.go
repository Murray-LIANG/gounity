package gounity

import (
	"encoding/json"
)

// NewHostLunById constructs a `HostLun` object with id.
func (u *Unity) NewHostLunById(id string) *HostLun {
	return &HostLun{
		Resource: Resource{
			typeName: typeNameHostLun, typeFields: typeFieldsHostLun, unity: u,
		},
		Id: id,
	}
}

// NewHostLunByName constructs a `HostLun` object with name.
func (u *Unity) NewHostLunByName(name string) *HostLun {
	return &HostLun{
		Resource: Resource{
			typeName: typeNameHostLun, typeFields: typeFieldsHostLun, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *HostLun) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
	}

	var (
		latest *HostLun
		err    error
	)
	switch r.Id {
	case "":
		if latest, err = r.unity.GetHostLunByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetHostLunById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetHostLunById retrives the `HostLun` by given its id.
func (u *Unity) GetHostLunById(id string) (*HostLun, error) {
	res := u.NewHostLunById(id)

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetHostLunByName retrives the `HostLun` by given its name.
func (u *Unity) GetHostLunByName(name string) (*HostLun, error) {
	res := u.NewHostLunByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetHostLuns retrives all `HostLun` objects.
func (u *Unity) GetHostLuns() ([]*HostLun, error) {

	return u.FilterHostLuns(nil)
}

// FilterHostLuns filters the `HostLun` objects by given filters.
func (u *Unity) FilterHostLuns(filter *filter) ([]*HostLun, error) {
	respEntries, err := u.getCollection(typeNameHostLun, typeFieldsHostLun, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillHostLuns(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillHostLuns(respEntries []*instanceResp) ([]*HostLun, error) {
	resSlice := []*HostLun{}
	for _, entry := range respEntries {
		res := u.NewHostLunById("") // empty id for fake `HostLun` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `HostLun` object using its id.
func (r *HostLun) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
