package gounity

import (
	"github.com/pkg/errors"
)

type HostLUNOperator interface {
	HostLUNOperatorGen

	FilterHostLunByHostAndLun(hostId, lunId string) (*HostLUN, error)
}

// FilterHostLunByHostAndLun filters the `HostLun` by given its host Id and Lun Id.
func (u *Unity) FilterHostLunByHostAndLun(hostId, lunId string) (*HostLUN, error) {
	filter := NewFilterf(`host eq "%s"`, hostId).Andf(`lun eq "%s"`, lunId)

	fields := map[string]interface{}{
		"hostId": hostId,
		"lunId":  lunId,
		"filter": filter,
	}
	msg := newMessage().withFields(fields)

	hostLuns, err := u.FilterHostLUNs(filter)
	if err != nil {
		return nil, errors.Wrapf(err, "filter hostluns failed: %s", msg)
	}
	if len(hostLuns) == 0 {
		return nil, errors.Wrapf(err, "filter returned 0 hostlun: %s", msg)
	}
	if len(hostLuns) > 1 {
		return nil, errors.Wrapf(
			err,
			"filter returned more than one hostluns: %s",
			msg.withField("numOfHostLun", len(hostLuns)),
		)
	}
	return hostLuns[0], nil
}
