package gounity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPoolByID(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	pool, err := ctx.unity.GetPoolByID("pool_1")
	assert.Nil(t, err)

	assert.Equal(t, "pool_1", pool.ID)
}

func TestGetPools(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	pools, err := ctx.unity.GetPools()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(pools))
	ids := []string{}
	for _, pool := range pools {
		ids = append(ids, pool.ID)
	}
	assert.EqualValues(t, []string{"pool_1", "pool_2", "pool_7", "pool_9"}, ids)
}

func TestCreateLUN(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	pool, err := ctx.unity.GetPoolByID("pool_1")
	assert.Nil(t, err)

	lun, err := pool.CreateLUN("lun-gounity", 3)
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.ID)
}

func TestCreateLUNNameExist(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	pool, err := ctx.unity.GetPoolByID("pool_1")
	assert.Nil(t, err)

	_, err = pool.CreateLUN("lun-name-exist-gounity", 3)
	assert.NotNil(t, err)

	unityError, ok := err.(*UnityError)
	assert.True(t, ok)
	assert.Equal(t, UnityLunNameExistErrorCode, unityError.ErrorCode)
}
