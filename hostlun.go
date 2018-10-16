package gounity

import (
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	fieldsHostLUN = strings.Join([]string{
		"hlu",
		"host.id",
		"id",
		"isDefaultSnap",
		"isReadOnly",
		"lun.id",
		"type",
	}, ",")
)

// FilterHostLUN filters the `HostLun` by given its host Id and Lun Id.
func (u *Unity) FilterHostLUN(hostID, lunID string) (*HostLun, error) {
	filter := newFilter(fmt.Sprintf(`host eq "%s"`, hostID)).and(
		fmt.Sprintf(`lun eq "%s"`, lunID))
	collection, err := u.getCollection("hostLUN", fieldsHostLUN, filter,
		reflect.TypeOf(HostLun{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*HostLun)
	if len(res) == 0 {
		log.WithField("hostID", hostID).WithField("lunID",
			lunID).Info("filter returns 0 hostLUN")
		return nil, nil
	}
	if len(res) > 1 {
		log.WithField("hostID", hostID).WithField("lunID", lunID).WithField("resultCount",
			len(res)).Info("filter returns more one hostLUNs")
	}
	return res[0], nil
}
