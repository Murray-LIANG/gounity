package gounity_test

import (
	"testing"

	"github.com/Murray-LIANG/gounity"
	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetLunById(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	lun, err := ctx.Unity.GetLunById("sv_1")
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.Id)
	assert.Equal(t, "pool_1", lun.Pool.Id)
}

func TestGetLunByIdNotFound(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	_, err = ctx.Unity.GetLunById("sv_2")
	assert.NotNil(t, err)

	assert.True(t, gounity.IsUnityResourceNotFoundError(err))
}

func TestGetLuns(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	luns, err := ctx.Unity.GetLuns()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(luns))
	ids := []string{}
	for _, lun := range luns {
		ids = append(ids, lun.Id)
	}
	assert.EqualValues(t, []string{"sv_1", "sv_3", "sv_16", "sv_17"}, ids)
}

func TestLunDelete(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	lun := ctx.Unity.NewLunById("sv_1")
	assert.Nil(t, err)

	err = lun.Delete()
	assert.Nil(t, err)
}
