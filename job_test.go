package gounity_test

import (
	"github.com/Murray-LIANG/gounity"
	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJob_CreateJob(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	tasks := []*gounity.JobTaskRequest{
		{
			Name:         "task",
			Description:  "",
			SubmitTime:   "",
			StartTime:    "",
			Object:       "",
			Action:       "",
			Dependencies: []string{},
		},
	}

	job := &gounity.Job{}
	job.Unity = ctx.Unity

	err = job.Create(tasks, "description")
	assert.Nil(t, err)
	assert.Equal(t, "36674839992", job.Id)
}

func TestJob_CancelJob(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	job, err := ctx.Unity.GetJobById("7455367723")
	assert.Nil(t, err)
	assert.NotNil(t, job)

	err = job.Cancel()
	assert.Nil(t, err)
}

func TestJob_Delete(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	job, err := ctx.Unity.GetJobById("7455362334")
	assert.Nil(t, err)
	assert.NotNil(t, job)

	err = job.Delete()
	assert.Nil(t, err)
}
