package gounity

import (
	"encoding/json"
)

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
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
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

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetDUMMYByName retrives the `DUMMY` by given its name.
func (u *Unity) GetDUMMYByName(name string) (*DUMMY, error) {
	res := u.NewDUMMYByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetDUMMYs retrives all `DUMMY` objects.
func (u *Unity) GetDUMMYs() ([]*DUMMY, error) {

	return u.FilterDUMMYs(nil)
}

// FilterDUMMYs filters the `DUMMY` objects by given filters.
func (u *Unity) FilterDUMMYs(filter *filter) ([]*DUMMY, error) {
	respEntries, err := u.getCollection(typeNameDUMMY, typeFieldsDUMMY, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillDUMMYs(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillDUMMYs(respEntries []*instanceResp) ([]*DUMMY, error) {
	resSlice := []*DUMMY{}
	for _, entry := range respEntries {
		res := u.NewDUMMYById("") // empty id for fake `DUMMY` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `DUMMY` object using its id.
func (r *DUMMY) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
