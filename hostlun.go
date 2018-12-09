package gounity

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	fieldsHostLun = strings.Join([]string{
		"hlu",
		"host.id",
		"id",
		"isDefaultSnap",
		"isReadOnly",
		"lun.id",
		"type",
	}, ",")
)

// FilterHostLun filters the `HostLun` by given its host Id and Lun Id.
func (u *Unity) FilterHostLun(hostId, lunId string) (*HostLun, error) {
	filter := newFilter(fmt.Sprintf(`host eq "%s"`, hostId)).and(
		fmt.Sprintf(`lun eq "%s"`, lunId))
	collection, err := u.getCollection("hostLun", fieldsHostLun, filter,
		reflect.TypeOf(HostLun{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*HostLun)
	if len(res) == 0 {
		log.WithField("hostId", hostId).WithField("lunId",
			lunId).Info("filter returns 0 hostLun")
		return nil, nil
	}
	if len(res) > 1 {
		log.WithField("hostId", hostId).WithField("lunId", lunId).WithField("resultCount",
			len(res)).Info("filter returns more one hostLuns")
	}
	return res[0], nil
}
