package gounity_test

import (
	"github.com/Murray-LIANG/gounity"
	"testing"

	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreateNfsShare(t *testing.T) {

	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	fs, err := ctx.Unity.GetFilesystemById("fs_1")
	assert.Nil(t, err)

	nfs, err := fs.CreateNfsShare("nfs_1")
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}

func TestCreateNfsShareWithDefaultAccess(t *testing.T) {

	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	fs, err := ctx.Unity.GetFilesystemById("fs_1")
	assert.Nil(t, err)

	nfs, err := fs.CreateNfsShare(
		"nfs_1", gounity.DefaultAccessOpt(gounity.NFSShareDefaultAccessReadWrite),
	)
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}