// +build integration

package gounity

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mgmtIp   = flag.String("mgmtIp", "", "Mgmt IP of Unity.")
	username = flag.String("username", "", "Username of Unity.")
	password = flag.String("password", "", "Password of Unity.")
)

func TestE2EGetPools(t *testing.T) {
	u, err := NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	pools, err := u.GetPools()
	assert.Nil(t, err, "failed to get all pools")

	assert.NotEmpty(t, pools, "no pools on unity")

	for _, pool := range pools {
		assert.NotEmpty(t, pool.Id, "no Id for pool")
		assert.Truef(t, strings.HasPrefix(pool.Id, "pool_"),
			"pool Id: %v is not with format `pool_x`", pool.Id)
	}
}
