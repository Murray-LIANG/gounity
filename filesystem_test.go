package gounity_test

import (
	"github.com/factioninc/gounity"
	"testing"

	"github.com/factioninc/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestExportNfsShare(t *testing.T) {

	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	fs, err := ctx.Unity.GetFilesystemById("fs_1")
	assert.Nil(t, err)

	nfs, err := fs.ExportNfsShare("nfs_1")
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}

func TestExportNfsShareWithOpt(t *testing.T) {

	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	fs, err := ctx.Unity.GetFilesystemById("fs_1")
	assert.Nil(t, err)

	nfs, err := fs.ExportNfsShare(
		"nfs_1", gounity.DefaultAccessOpt(gounity.NFSShareDefaultAccessReadWrite),
	)
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}

func TestFilesystemDelete(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	fs := ctx.Unity.NewFilesystemById("fs_1")
	assert.Nil(t, err)

	err = fs.Delete()
	assert.Nil(t, err)
}
