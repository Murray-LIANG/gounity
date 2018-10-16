package gounity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterHostLun(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	hostLun, err := ctx.unity.FilterHostLUN("Host_1", "sv_1")
	assert.Nil(t, err)

	assert.Equal(t, uint16(1), hostLun.Hlu)
}
