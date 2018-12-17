package gounity

import (
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
)

type genNasServerOperator interface {
	NewNasServerById(id string) *NasServer

	NewNasServerByName(name string) *NasServer

	GetNasServerById(id string) (*NasServer, error)

	GetNasServerByName(name string) (*NasServer, error)

	GetNasServers() ([]*NasServer, error)

	FillNasServers(respEntries []*instanceResp) ([]*NasServer, error)

	FilterNasServers(filter *filter) ([]*NasServer, error)
}

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
		return fmt.Errorf(
			"cannot refresh on resource without Id nor Name, resource:%v", r,
		)
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

	if err := u.GetInstanceById(res.typeName, id, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get NasServer by id failed")
	}
	return res, nil
}

// GetNasServerByName retrives the `NasServer` by given its name.
func (u *Unity) GetNasServerByName(name string) (*NasServer, error) {
	res := u.NewNasServerByName(name)
	if err := u.GetInstanceByName(res.typeName, name, res.typeFields, res); err != nil {
		return nil, errors.Wrap(err, "get NasServer by name failed")
	}
	return res, nil
}

// GetNasServers retrives all `NasServer` objects.
func (u *Unity) GetNasServers() ([]*NasServer, error) {

	return u.FilterNasServers(nil)
}

// FilterNasServers filters the `NasServer` objects by given filters.
func (u *Unity) FilterNasServers(filter *filter) ([]*NasServer, error) {
	respEntries, err := u.GetCollection(typeNameNasServer, typeFieldsNasServer, filter)
	if err != nil {
		return nil, errors.Wrap(err, "filter NasServer failed")
	}
	res, err := u.FillNasServers(respEntries)
	if err != nil {
		return nil, errors.Wrap(err, "fill NasServers failed")
	}
	return res, nil
}

// FillNasServers generates the `NasServer` objects from collection query response.
func (u *Unity) FillNasServers(respEntries []*instanceResp) ([]*NasServer, error) {
	resSlice := []*NasServer{}
	for _, entry := range respEntries {
		res := u.NewNasServerById("") // empty id for fake `NasServer` object
		if err := json.Unmarshal(entry.Content, res); err != nil {
			return nil, errors.Wrapf(err, "decode to %v failed", res)
		}
		resSlice = append(resSlice, res)
	}
	return resSlice, nil
}

// Repr represents a `NasServer` object using its id.
func (r *NasServer) Repr() *idRepresent {
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
