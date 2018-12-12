package gounity

import (
	"encoding/json"
)

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
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
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

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetHostByName retrives the `Host` by given its name.
func (u *Unity) GetHostByName(name string) (*Host, error) {
	res := u.NewHostByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetHosts retrives all `Host` objects.
func (u *Unity) GetHosts() ([]*Host, error) {

	return u.FilterHosts(nil)
}

// FilterHosts filters the `Host` objects by given filters.
func (u *Unity) FilterHosts(filter *filter) ([]*Host, error) {
	respEntries, err := u.getCollection(typeNameHost, typeFieldsHost, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillHosts(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillHosts(respEntries []*instanceResp) ([]*Host, error) {
	resSlice := []*Host{}
	for _, entry := range respEntries {
		res := u.NewHostById("") // empty id for fake `Host` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Host` object using its id.
func (r *Host) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
