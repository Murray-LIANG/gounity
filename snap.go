package gounity

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (s *Snap) Copy(copyName string) (*Snap, error) {
	body := map[string]interface{}{
		"copyName": copyName,
	}

	fields := map[string]interface{}{
		"requestBody": body,
	}

	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("copying snapshot")
	err := s.Unity.PostOnInstance(
		typeNameSnap, s.Id, "copy", body,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "copying snapshot failed: %s", msg)
	}

	snap := s.Unity.NewSnapByName(copyName)
	if err = snap.Refresh(); err != nil {
		return nil, errors.Wrapf(err, "get snapshot failed: %s", msg)
	}

	log.WithField("copySnapId", snap.Id).Debug("Snapshot successfully copied")
	return snap, err
}

func (s *Snap) AttachToHost(host *Host, access SnapAccessLevelEnum) (error) {
	hostAccess := []interface{}{
		map[string]interface{}{
			"host":       *host.Repr(),
			"allowedAccess": access,
		},
	}

	fields := map[string]interface{}{
		"requestBody": hostAccess,
	}

	body := map[string]interface{}{"hostAccess": hostAccess}

	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("attaching snapshot")
	err := s.Unity.PostOnInstance(
		typeNameSnap, s.Id, "attach", body,
	)
	if err != nil {
		return errors.Wrapf(err, "attaching snapshot failed: %s", msg)
	}

	log.WithField("copySnapId", s.Id).Debug("Snapshot successfully attached")
	return err
}

func (s *Snap) DetachFromHost() (error) {
	body := map[string]interface{}{}

	logrus.Debug("detaching snapshot")
	err := s.Unity.PostOnInstance(
		typeNameSnap, s.Id, "detach", body,
	)
	if err != nil {
		return errors.Wrapf(err, "detaching snapshot failed: %s", err)
	}

	logrus.WithField("snapId", s.Id).Debug("Snapshot successfully attached")
	return err
}
