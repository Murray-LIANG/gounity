package gounity

import (
	"strings"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	typeNameHost   = "host"
	typeFieldsHost = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
		"osType",
	}, ",")
)

type HostOperator interface {
	genHostOperator
}

// Host defines Unity corresponding `host` type.
type Host struct {
	Resource
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
	OsType      string  `json:"osType"`
}

//go:generate ./gen_resource.sh resource_tmpl.go host_gen.go Host

// Attach attaches the Lun to the host.
func (h *Host) Attach(lun *Lun) (uint16, error) {
	hostAccess := []interface{}{
		map[string]interface{}{
			"host":       h.Repr(),
			"accessMask": HostLunAccessProduction,
		},
	}
	for _, exist := range lun.HostAccess {
		hostAccess = append(hostAccess,
			map[string]interface{}{
				"host":       exist.Host.Repr(),
				"accessMask": exist.AccessMask,
			},
		)
	}

	body := map[string]interface{}{
		"lunParameters": map[string]interface{}{"hostAccess": hostAccess},
	}

	fields := map[string]interface{}{
		"host":        h,
		"lun":         lun,
		"requestBody": body,
	}

	logger := log.WithFields(fields)
	msg := newMessage().withFields(fields)

	logger.Debug("attaching lun to host")
	if err := h.unity.PostOnInstance(
		typeStorageResource, lun.Id, actionModifyLun, body,
	); err != nil {
		return 0, errors.Wrap(err, msg.withMessage("attach lun to host failed").String())
	}

	hostLun, err := h.unity.FilterHostLunByHostAndLun(h.Id, lun.Id)
	if err != nil {
		return 0, errors.Wrap(
			err, msg.withMessage("filter hostlun by host and lun failed").String(),
		)
	}
	return hostLun.Hlu, nil
}
