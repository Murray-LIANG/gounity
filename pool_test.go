package gounity_test

import (
	"testing"

	"github.com/factioninc/gounity"
	"github.com/factioninc/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetPoolById(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	assert.Equal(t, "pool_1", pool.Id)
}

func TestGetPools(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pools, err := ctx.Unity.GetPools()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(pools))
	ids := []string{}
	for _, pool := range pools {
		ids = append(ids, pool.Id)
	}
	assert.EqualValues(t, []string{"pool_1", "pool_2", "pool_7", "pool_9"}, ids)
}

func TestNewPoolByIdThenRefresh(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool := ctx.Unity.NewPoolById("pool_1")
	assert.Equal(t, "pool_1", pool.Id)
	assert.Empty(t, pool.Name)

	err = pool.Refresh()
	assert.Nil(t, err)
	assert.Equal(t, "pool_1", pool.Id)
	assert.Equal(t, "Manila_Pool", pool.Name)
}

func TestNewPoolByNameThenRefresh(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool := ctx.Unity.NewPoolByName("Manila_Pool")
	assert.Empty(t, pool.Id)
	assert.Equal(t, "Manila_Pool", pool.Name)

	err = pool.Refresh()
	assert.Nil(t, err)
	assert.Equal(t, "pool_1", pool.Id)
	assert.Equal(t, "Manila_Pool", pool.Name)
}

func TestCreateLun(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	lun, err := pool.CreateLun("lun-gounity", 3)
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.Id)
}

func TestCreateLunWithOpt(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	host := ctx.Unity.NewHostById("Host_1")
	lun, err := pool.CreateLun(
		"lun-gounity", 3,
		gounity.HostAccessOpt(host, gounity.HostLUNAccessProduction),
	)
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.Id)
}

func TestCreateLunNameExist(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	_, err = pool.CreateLun("lun-name-exist-gounity", 3)
	assert.NotNil(t, err)

	assert.True(t, gounity.IsUnityLunNameExistError(err))
}

func TestCreateFilesystem(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	nas, err := ctx.Unity.GetNasServerById("nas_1")
	assert.Nil(t, err)

	fs, err := pool.CreateFilesystem(nas, "fs-name", 3)
	assert.Nil(t, err)

	assert.Equal(t, "fs_1", fs.Id)
}

func TestCreateFilesystemWithOpt(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	nas, err := ctx.Unity.GetNasServerById("nas_1")
	assert.Nil(t, err)

	fs, err := pool.CreateFilesystem(
		nas, "fs-name", 3,
		gounity.SupportedProtocolsOpt(gounity.FSSupportedProtocolNFS),
	)
	assert.Nil(t, err)

	assert.Equal(t, "fs_1", fs.Id)
}

func TestCreateNfsShare(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	nas, err := ctx.Unity.GetNasServerById("nas_1")
	assert.Nil(t, err)

	nfs, err := pool.CreateNfsShare(nas, "nfsshare-name", 3)
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}

func TestCreateNfsShareWithOpt(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	nas, err := ctx.Unity.GetNasServerById("nas_1")
	assert.Nil(t, err)

	nfs, err := pool.CreateNfsShare(
		nas, "nfsshare-name", 3,
		gounity.DefaultAccessOpt(gounity.NFSShareDefaultAccessReadWrite),
	)
	assert.Nil(t, err)

	assert.Equal(t, "NFSShare_1", nfs.Id)
}
