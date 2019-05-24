// +build integration

package gounity_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/factioninc/gounity"
	"github.com/stretchr/testify/assert"
)

var (
	mgmtIp   = flag.String("mgmtIp", "", "Mgmt IP of Unity.")
	username = flag.String("username", "", "Username of Unity.")
	password = flag.String("password", "", "Password of Unity.")
)

func TestE2EGetPools(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	pools, err := u.GetPools()
	assert.Nil(t, err, "get all pools failed")

	assert.NotEmpty(t, pools, "no pools on unity")

	for _, pool := range pools {
		assert.NotEmpty(t, pool.Id, "no Id for pool")
		assert.Truef(t, strings.HasPrefix(pool.Id, "pool_"),
			"pool Id: %v is not with format `pool_x`", pool.Id)
	}
}

func TestE2ECreateLun(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	// pool := u.NewPoolByName("Manila_Pool")
	pool := u.NewPoolById("pool_1")

	newLun, err := pool.CreateLun("lun-gounity", 3)
	assert.Nil(t, err, "create lun failed")

	assert.Equal(t, "lun-gounity", newLun.Name)
}

func TestE2EAttach(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	host := u.NewHostByName("windows10-liangr")
	lun := u.NewLunByName("lun-gounity")

	hlu, err := host.Attach(lun)
	assert.Nil(t, err, "attach lun failed")

	assert.Equal(t, uint16(0), hlu)
}

func TestE2EDetach(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	host := u.NewHostByName("windows10-liangr")
	lun := u.NewLunByName("lun-gounity")

	err = host.Detach(lun)
	assert.Nil(t, err, "detach lun failed")

	hostLun, err := u.FilterHostLunByHostAndLun(host.Id, lun.Id)
	assert.Nil(t, hostLun, "hostLun still exists")
	assert.NotNil(t, err, "hostLun still exists")
}

func TestE2EDeleteLun(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	lun := u.NewLunByName("lun-gounity")

	err = lun.Delete()
	assert.Nil(t, err, "delete lun failed")

	_, err = u.GetLunByName("lun-gounity")
	assert.True(t, gounity.IsUnityResourceNotFoundError(err))
}

func TestE2ECreateNfsShare(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	pool := u.NewPoolByName("Manila_Pool")
	nasServer := u.NewNasServerByName("nas-ht")

	nfs, err := pool.CreateNfsShare(
		nasServer, "nfs-gounity", 3,
		gounity.DefaultAccessOpt(gounity.NFSShareDefaultAccessReadWrite),
	)
	assert.Nil(t, err, "create nfs share failed")

	assert.Equal(t, []string{"10.245.47.97:/nfs-gounity"}, nfs.ExportPaths)
}

func TestE2EDeleteNfsShare(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	nfs := u.NewNfsShareByName("nfs-gounity")

	err = nfs.Delete()
	assert.Nil(t, err, "delete nfs share failed")

	_, err = u.GetNfsShareByName("nfs-gounity")
	assert.True(t, gounity.IsUnityResourceNotFoundError(err))
}

func TestE2EDeleteFilesystem(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	fs := u.NewFilesystemByName("nfs-gounity")

	err = fs.Delete()
	assert.Nil(t, err, "delete filesystem failed")

	_, err = u.GetFilesystemByName("nfs-gounity")
	assert.True(t, gounity.IsUnityResourceNotFoundError(err))
}

func TestE2EGetIscsiPortals(t *testing.T) {
	u, err := gounity.NewUnity(*mgmtIp, *username, *password, true)
	assert.Nilf(t, err, "cannot connect to unity: %v", *mgmtIp)

	portals, err := u.GetIscsiPortals()
	assert.Nil(t, err, "get all iscsi portals failed")

	portal := portals[0]
	assert.Equal(t, "if_3", portal.Id)
	assert.Equal(t, "10.245.47.95", portal.IpAddress)
	assert.Equal(t, "iqn.1992-04.com.emc:cx.fnm00150600267.a0", portal.IscsiNode.Name)
}
