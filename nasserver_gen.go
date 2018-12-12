package gounity

import (
	"encoding/json"
)

// NewNasServerById constructs a `NasServer` object with id.
func (u *Unity) NewNasServerById(id string) *NasServer {
	return &NasServer{
		Resource: Resource{
			typeName: typeNameNasServer, typeFields: typeFieldsNasServer, unity: u,
		},
		Id: id,
	}
}

// NewNasServerByName constructs a `NasServer` object with name.
func (u *Unity) NewNasServerByName(name string) *NasServer {
	return &NasServer{
		Resource: Resource{
			typeName: typeNameNasServer, typeFields: typeFieldsNasServer, unity: u,
		},
		Name: name,
	}
}

// Refresh updates the info from Unity.
func (r *NasServer) Refresh() error {
	if r.Id == "" && r.Name == "" {
		return newGounityError(
			"cannot refresh on resource without Id nor Name").withField("resource", r)
	}

	var (
		latest *NasServer
		err    error
	)
	switch r.Id {
	case "":
		if latest, err = r.unity.GetNasServerByName(r.Name); err != nil {
			return err
		}
		r = latest
	default:
		if latest, err = r.unity.GetNasServerById(r.Id); err != nil {
			return err
		}
		r = latest
	}
	return nil
}

// GetNasServerById retrives the `NasServer` by given its id.
func (u *Unity) GetNasServerById(id string) (*NasServer, error) {
	res := u.NewNasServerById(id)

	if err := u.getInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNasServerByName retrives the `NasServer` by given its name.
func (u *Unity) GetNasServerByName(name string) (*NasServer, error) {
	res := u.NewNasServerByName(name)
	if err := u.getInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetNasServers retrives all `NasServer` objects.
func (u *Unity) GetNasServers() ([]*NasServer, error) {

	return u.FilterNasServers(nil)
}

// FilterNasServers filters the `NasServer` objects by given filters.
func (u *Unity) FilterNasServers(filter *filter) ([]*NasServer, error) {
	respEntries, err := u.getCollection(typeNameNasServer, typeFieldsNasServer, filter)
	if err != nil {
		return nil, err
	}
	res, err := u.fillNasServers(respEntries)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *Unity) fillNasServers(respEntries []*instanceResp) ([]*NasServer, error) {
	resSlice := []*NasServer{}
	for _, entry := range respEntries {
		res := u.NewNasServerById("") // empty id for fake `NasServer` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, err
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `NasServer` object using its id.
func (r *NasServer) Repr() *idRepresent {
	if r.Id == "" {
		err := r.Refresh()
		if err != nil {
			// TODO (ryan) Add log here
			return nil
		}
	}
	return &idRepresent{Id: r.Id}
}
