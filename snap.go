package gounity

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)
type createSnapshotResourceResp struct {
	Content struct {
		Id string `json:"id"`
	} `json:"content"`
}

func newCreateSnapshotBody(s *Snap, sr *StorageResource) map[string]interface{} {
	body := map[string]interface{}{
		"name": s.Name,
		"storageResource": *sr.Repr(),
	}

	return body
}

func (s *Snap) Create(sr *StorageResource) (error) {
	body := newCreateSnapshotBody(s, sr)

	fields := map[string]interface{}{
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("creating snapshot")
	resp := &createSnapshotResourceResp{}
	if err := s.Unity.CreateOnType(typeNameSnap, body, resp); err != nil {
		return errors.Wrapf(err, "create snapshot failed: %s", err)
	}

	createdId := resp.Content.Id

	log.WithField("createdSnapshotId", createdId).Debug("snapshot created")

	snap, err := s.Unity.GetSnapById(createdId)
	if err != nil {
		return errors.Wrapf(
			err,
			"could not retrieve snapshot: %s", msg.withField("createdSnapshotId", createdId),
		)
	}
	*s = *snap
	return err
}

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
