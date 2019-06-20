package gounity

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type createJobResourceResp struct {
	Content struct {
		Id string `json:"id"`
	} `json:"content"`
}

func newCreateJobBody(jobTaskRequests []*JobTaskRequest, description string) map[string]interface{} {
	body := map[string]interface{}{
		"description": description,
		"tasks":       jobTaskRequests,
	}

	return body
}

func (j *Job) Create(jobTaskRequests []*JobTaskRequest, description string) error {
	body := newCreateJobBody(jobTaskRequests, description)

	fields := map[string]interface{}{
		"requestBody": body,
	}
	log := logrus.WithFields(fields)
	msg := newMessage().withFields(fields)

	log.Debug("creating job")
	resp := &createJobResourceResp{}
	if err := j.Unity.CreateOnType(typeNameJob, body, resp); err != nil {
		return errors.Wrapf(err, "create job failed: %s", err)
	}

	log.WithField("createdJobId", j.Id).Debug("job created")
	j.Id = resp.Content.Id

	err := j.Refresh()
	if err != nil {
		return errors.Wrapf(
			err,
			"could not retrieve job: %s", msg.withField("createdJobId", j.Id),
		)
	}

	return err
}

func (j *Job) Cancel() error {
	body := map[string]interface{}{}

	logrus.Debug("cancel job")
	err := j.Unity.PostOnInstance(
		typeNameJob, j.Id, "cancel", body, nil,
	)

	if err != nil {
		return errors.Wrapf(err, "canceling job failed: %s", err)
	}

	logrus.WithField("jobId", j.Id).Debug("Job successfully canceled")
	return err
}

func (j *Job) Delete() error {
	body := map[string]interface{}{}

	logrus.Debug("delete job")
	err := j.Unity.PostOnInstance(
		typeNameJob, j.Id, "delete", body, nil,
	)

	if err != nil {
		return errors.Wrapf(err, "deleting job failed: %s", err)
	}

	logrus.WithField("jobId", j.Id).Debug("Job successfully deleted")
	return err
}
