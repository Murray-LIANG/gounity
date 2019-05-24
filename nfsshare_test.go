package gounity_test

import (
	"testing"

	"github.com/factioninc/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNfsShareDelete(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	nfs := ctx.Unity.NewNfsShareById("NFSShare_1")
	assert.Nil(t, err)

	err = nfs.Delete()
	assert.Nil(t, err)

}
