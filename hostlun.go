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

// FilterHostLUN filters the `HostLUN` by given its host ID and LUN ID.
func (u *Unity) FilterHostLUN(hostID, lunID string) (*HostLUN, error) {
	filter := newFilter(fmt.Sprintf(`host eq "%s"`, hostID)).and(
		fmt.Sprintf(`lun eq "%s"`, lunID))
	collection, err := u.getCollection("hostLUN", fieldsHostLUN, filter,
		reflect.TypeOf(HostLUN{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*HostLUN)
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
