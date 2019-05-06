package gounity_test

import (
	"testing"

	"github.com/Murray-LIANG/gounity"
	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestStorageResource_CreateSnapshot(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	sr, err := ctx.Unity.GetStorageResourceById("sv_1")
	assert.Nil(t, err)
	snap := ctx.Unity.NewSnapByName("new_snap")
	assert.Nil(t, err)
	assert.NotNil(t, snap)
	err = snap.Create(sr)
	assert.Nil(t, err)
	assert.Equal(t, "38654714905", snap.Id)
	assert.Equal(t, "new_snap", snap.Name)
}


func TestSnap_GetFilter(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	filter := gounity.NewFilter(`lun.id eq "sv_1"`)
	snap, err := ctx.Unity.FilterSnaps(filter)
	assert.Nil(t, err)

	lun := snap[0].Lun
	assert.Equal(t, "38654714567", snap[0].Id)
	assert.Equal(t, "sv_1", lun.Id)
	err = lun.Refresh()
	assert.Nil(t, err)
	assert.Equal(t, "Lun1", lun.Name)
}

func TestSnap_CopySnap(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()
	snap, err := ctx.Unity.GetSnapById("38654714770")
	assert.Nil(t, err)

	snapCopy, err := snap.Copy("new_snap")
	assert.Nil(t, err)
	assert.Equal(t, "new_snap", snapCopy.Name)
	assert.Equal(t, "38654714905", snapCopy.Id)
}

func TestSnap_AttachToHost(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()
	snap, err := ctx.Unity.GetSnapById("38654714770")
	assert.Nil(t, err)

	host := ctx.Unity.NewHostById("Host_1")
	assert.Nil(t, err)
	err = snap.AttachToHost(host, gounity.SnapAccessLevelReadWrite)
	assert.Nil(t, err)
}

func TestSnap_DetachFromHost(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()
	snap, err := ctx.Unity.GetSnapById("38654714770")
	assert.Nil(t, err)
	err = snap.DetachFromHost()
	assert.Nil(t, err)
}
