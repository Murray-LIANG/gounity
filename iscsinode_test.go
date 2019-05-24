package gounity_test

import (
	"testing"

	"github.com/factioninc/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetIscsiNodeById(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	node, err := ctx.Unity.GetIscsiNodeById("iscsinode_spa_eth2")
	assert.Nil(t, err)

	assert.Equal(t, "iscsinode_spa_eth2", node.Id)
	assert.Equal(t, "0267.a0", node.Alias)
	assert.Equal(t, "iqn.1992-04.com.emc:cx.fnm00150600267.a0", node.Name)
	assert.Equal(t, "spa_eth2", node.EthernetPort.Id)
}

func TestGetIscsiNodes(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	nodes, err := ctx.Unity.GetIscsiNodes()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(nodes))

	expectedId := []string{
		"iscsinode_spa_eth2", "iscsinode_spa_eth3",
		"iscsinode_spb_eth2", "iscsinode_spb_eth3",
	}
	expectedAlias := []string{"0267.a0", "0267.a1", "0267.b0", "0267.b1"}
	expectedName := []string{
		"iqn.1992-04.com.emc:cx.fnm00150600267.a0",
		"iqn.1992-04.com.emc:cx.fnm00150600267.a1",
		"iqn.1992-04.com.emc:cx.fnm00150600267.b0",
		"iqn.1992-04.com.emc:cx.fnm00150600267.b1",
	}
	expectedPortId := []string{"spa_eth2", "spa_eth3", "spb_eth2", "spb_eth3"}
	for i, node := range nodes {
		assert.Equal(t, expectedId[i], node.Id)
		assert.Equal(t, expectedAlias[i], node.Alias)
		assert.Equal(t, expectedName[i], node.Name)
		assert.Equal(t, expectedPortId[i], node.EthernetPort.Id)
	}
}
