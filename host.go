package gounity

import (
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

// Attach attaches the Lun to the host.
func (h *Host) Attach(lun *Lun) (uint16, error) {
	hostAccess := []interface{}{
		map[string]interface{}{
			"host":       h.Repr(),
			"accessMask": HostLUNAccessProduction,
		},
	}

	// Refresh lun from unity to get latest its hostAccess.
	if err := lun.Refresh(); err != nil {
		return 0, errors.Wrapf(err, "refresh lun failed")
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

	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("attaching lun to host")
	if err := h.Unity.PostOnInstance(
		typeStorageResource, lun.Id, actionModifyLun, body,
	); err != nil {
		return 0, errors.Wrapf(err, "attach lun to host failed: %s", msg)
	}

	hostLun, err := h.Unity.FilterHostLunByHostAndLun(h.Id, lun.Id)
	if err != nil {
		return 0, errors.Wrapf(err, "filter hostlun by host and lun failed: %s", msg)
	}
	return hostLun.Hlu, nil
}

// Detach detaches the Lun from the host.
func (h *Host) Detach(lun *Lun) error {

	// Refresh lun from unity to get latest its hostAccess.
	if err := lun.Refresh(); err != nil {
		return errors.Wrapf(err, "refresh lun failed")
	}
	hostAccess := []interface{}{}

	if h.Id == "" {
		if err := h.Refresh(); err != nil {
			return errors.Wrapf(err, "refresh host failed")
		}
	}
	for _, exist := range lun.HostAccess {
		if exist.Host.Id == h.Id {
			continue
		}
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

	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("detaching lun from host")
	if err := h.Unity.PostOnInstance(
		typeStorageResource, lun.Id, actionModifyLun, body,
	); err != nil {
		return errors.Wrapf(err, "detach lun from host failed: %s", msg)
	}
	return nil
}
