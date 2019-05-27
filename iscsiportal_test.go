package gounity_test

import (
	"testing"

	"github.com/Murray-LIANG/gounity"
	"github.com/Murray-LIANG/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetIscsiPortalById(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	portal, err := ctx.Unity.GetIscsiPortalById("if_3")
	assert.Nil(t, err)

	assert.Equal(t, "spa_eth2", portal.EthernetPort.Id)
	assert.Equal(t, "10.245.47.1", portal.Gateway)
	assert.Equal(t, "if_3", portal.Id)
	assert.Equal(t, "10.245.47.95", portal.IpAddress)
	assert.Equal(t, gounity.IpProtocolVersionIPv4, portal.IpProtocolVersion)
	assert.Equal(t, "iscsinode_spa_eth2", portal.IscsiNode.Id)
	assert.Equal(t, "iqn.1992-04.com.emc:cx.fnm00150600267.a0", portal.IscsiNode.Name)
	assert.Equal(t, "255.255.255.0", portal.Netmask)
	assert.Equal(t, uint32(0), portal.V6PrefixLength)
	assert.Equal(t, uint16(0), portal.VlanId)
}

func TestGetIscsiPortals(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	portals, err := ctx.Unity.GetIscsiPortals()
	assert.Nil(t, err)

	assert.Equal(t, 2, len(portals))

	expectedId := []string{"if_3", "if_4"}
	expectedIp := []string{"10.245.47.95", "10.245.47.96"}
	expectedIqn := []string{
		"iqn.1992-04.com.emc:cx.fnm00150600267.a0",
		"iqn.1992-04.com.emc:cx.fnm00150600267.b0",
	}
	for i, portal := range portals {
		assert.Equal(t, expectedId[i], portal.Id)
		assert.Equal(t, expectedIp[i], portal.IpAddress)
		assert.Equal(t, expectedIqn[i], portal.IscsiNode.Name)
	}
}
