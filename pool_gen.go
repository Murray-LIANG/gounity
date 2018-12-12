package gounity

import (
	"encoding/json"
)

// NewPoolById constructs a Pool object with id.
func (u *Unity) NewPoolById(id string) *Pool {
	return &Pool{
		Resource: Resource{
			typeName: typeNamePool, typeFields: typeFieldsPool, unity: u,
		},
		Id: id,
	}
}

// NewPoolByName constructs a Pool object with name.
func (u *Unity) NewPoolByName(name string) *Pool {
	return &Pool{
		Resource: Resource{
			typeName: typeNamePool, typeFields: typeFieldsPool, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from unity.
func (r *Pool) Refresh() (*Pool, error) {

	if res, err := r.unity.GetPoolById(r.Id); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// GetPoolById retrives the Pool by given its id.
func (u *Unity) GetPoolById(id string) (*Pool, error) {
	res := u.NewPoolById(id)

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetPoolByName retrives the Pool by given its name.
func (u *Unity) GetPoolByName(name string) (*Pool, error) {
	res := u.NewPoolByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetPools retrives all Pools.
func (u *Unity) GetPools() ([]*Pool, error) {

	respEntries, err := u.getCollection(typeNamePool, typeFieldsPool, nil)
	if err != nil {
		return nil, err
	}
	res, err := u.fillPools(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillPools(respEntries []*instanceResp) ([]*Pool, error) {
	resSlice := []*Pool{}
	for _, entry := range respEntries {
		res := u.NewPoolById("") // empty id for fake Pool object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a Pool object using its id.
func (r *Pool) Repr() *idRepresent {
	id := r.Id
	if id == "" {
		if r, err := r.Refresh(); err != nil {
			// TODO (ryan) Add log here
			return nil
		} else {
			id = r.Id
		}
	}
	return &idRepresent{Id: id}
}
