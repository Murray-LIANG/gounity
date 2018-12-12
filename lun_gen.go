package gounity

import (
	"encoding/json"
)

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
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
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

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetLunByName retrives the `Lun` by given its name.
func (u *Unity) GetLunByName(name string) (*Lun, error) {
	res := u.NewLunByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetLuns retrives all `Lun` objects.
func (u *Unity) GetLuns() ([]*Lun, error) {

	return u.FilterLuns(nil)
}

// FilterLuns filters the `Lun` objects by given filters.
func (u *Unity) FilterLuns(filter *filter) ([]*Lun, error) {
	respEntries, err := u.getCollection(typeNameLun, typeFieldsLun, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillLuns(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillLuns(respEntries []*instanceResp) ([]*Lun, error) {
	resSlice := []*Lun{}
	for _, entry := range respEntries {
		res := u.NewLunById("") // empty id for fake `Lun` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `Lun` object using its id.
func (r *Lun) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
