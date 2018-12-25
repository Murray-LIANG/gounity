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
		return nil, errors.Wrap(err, msg.withMessage("filter hostluns failed").String())
	}
	if len(hostLuns) == 0 {
		return nil, errors.Wrap(
			err, msg.withMessage("filter returned 0 hostlun").String(),
		)
	}
	if len(hostLuns) > 1 {
		return nil, errors.Wrap(
			err, msg.withField("numOfHostLun",
				len(hostLuns)).withMessage(
				"filter returned more than one hostluns").String(),
		)
	}
	return hostLuns[0], nil
}
