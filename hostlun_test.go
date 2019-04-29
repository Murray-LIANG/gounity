package gounity_test

import (
	"testing"

	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFilterHostLun(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	hostLun, err := ctx.Unity.FilterHostLunByHostAndLun("Host_1", "sv_1")

	assert.Nil(t, err)

	assert.Equal(t, uint16(1), hostLun.Hlu)
}
