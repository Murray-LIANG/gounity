package gounity

import (
	"strings"

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
			"host": h.Repr(),
			"accessMask": HostLunAccessProduction,
		},
	}
	for _, exist := range lun.HostAccess {
		hostAccess = append(hostAccess,
			map[string]interface{}{
				"host": exist.Host.Repr(),
				 "accessMask": exist.AccessMask,
			},
		)
	}

	body := map[string]interface{}{
		"lunParameters": map[string]interface{}{"hostAccess": hostAccess},
	}

	logger := log.WithField("host", h).WithField("lun", lun).WithField(
		"requestBody", body)
	logger.Debug("attaching lun to host")

	if err := h.unity.postOnInstance(typeStorageResource, lun.Id, actionModifyLun, body); err != nil {

		logger.WithError(err).Error("failed to attach lun to host")
		return 0, err
	}

	hostLun, err := h.unity.FilterHostLunByHostAndLun(h.Id, lun.Id)
	if err != nil {
		return 0, err
	}
	return hostLun.Hlu, nil
}
