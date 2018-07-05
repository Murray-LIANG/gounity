package gounity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLUNByID(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	lun, err := ctx.unity.GetLUNByID("sv_1")
	assert.Nil(t, err)

	assert.Equal(t, "sv_1", lun.ID)
	assert.Equal(t, "pool_1", lun.Pool.ID)
}

func TestGetLUNByIDNotFound(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	_, err = ctx.unity.GetLUNByID("sv_2")
	assert.NotNil(t, err)

	unityError, ok := err.(*UnityError)
	assert.True(t, ok)
	assert.Equal(t, UnityResourceNotFoundErrorCode, unityError.ErrorCode)
}

func TestGetLUNs(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	luns, err := ctx.unity.GetLUNs()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(luns))
	ids := []string{}
	for _, lun := range luns {
		ids = append(ids, lun.ID)
	}
	assert.EqualValues(t, []string{"sv_1", "sv_3", "sv_16", "sv_17"}, ids)
}
