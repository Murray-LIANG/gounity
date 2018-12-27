package gounity_test

import (
	"testing"

	"github.com/Murray-LIANG/gounity"
	"github.com/Murray-LIANG/gounity/testutil"
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

func TestCreateLun(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	lun, err := pool.CreateLun(gounity.NameOpt("lun-gounity"), gounity.SizeGBOpt(3))
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.Id)
}

func TestCreateLunNameExist(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	pool, err := ctx.Unity.GetPoolById("pool_1")
	assert.Nil(t, err)

	_, err = pool.CreateLun(gounity.NameOpt("lun-name-exist-gounity"), gounity.SizeGBOpt(3))
	assert.NotNil(t, err)

	assert.True(t, gounity.IsUnityLunNameExistError(err))
}
