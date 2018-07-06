package gounity

import (
	"context"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	fieldsHost = strings.Join([]string{
		"description",
		"health",
		"id",
		"name",
		"osType",
	}, ",")
)

// GetHostByID retrives the host by given its ID.
func (u *Unity) GetHostByID(id string) (*Host, error) {
	res := &Host{}
	if err := u.getInstanceByID("host", id, fieldsHost, res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetHosts retrives all the hosts.
func (u *Unity) GetHosts() ([]*Host, error) {
	collection, err := u.getCollection("host", fieldsHost, nil, reflect.TypeOf(Host{}))
	if err != nil {
		return nil, err
	}
	res := collection.([]*Host)
	return res, nil
}

// Attach attaches the LUN to the host.
func (h *Host) Attach(lun *LUN) (uint16, error) {
	hostAccess := []interface{}{
		map[string]interface{}{"host": represent(h),
			"accessMask": HostLUNAccessProduction},
	}
	for _, exist := range lun.HostAccess {
		hostAccess = append(hostAccess,
			map[string]interface{}{
				"host": represent(exist.Host), "accessMask": exist.AccessMask})
	}

	body := map[string]interface{}{
		"lunParameters": map[string]interface{}{"hostAccess": hostAccess},
	}

	logger := log.WithField("host", h).WithField("lun", lun).WithField(
		"requestBody", body)
	logger.Debug("attacthing lun to host")

	if err := h.Unity.client.Post(context.Background(),
		postInstanceURL("storageResource", lun.ID, "modifyLun"), nil, body,
		nil); err != nil {

		logger.WithError(err).Error("failed to attach lun to host")
		return 0, err
	}

	hostLun, err := h.Unity.FilterHostLUN(h.ID, lun.ID)
	if err != nil {
		return 0, err
	}
	return hostLun.Hlu, nil
}
