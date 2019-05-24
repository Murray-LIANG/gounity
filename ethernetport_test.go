package gounity_test

import (
	"testing"

	"github.com/factioninc/gounity/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetEthernetById(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	port, err := ctx.Unity.GetEthernetPortById("spa_eth2")
	assert.Nil(t, err)

	assert.Equal(t, "spa_eth2", port.Id)
	assert.Equal(t, 5, port.Health.Value)
	assert.Equal(t, "SP A Ethernet Port 2", port.Name)
	assert.Equal(t, true, port.IsLinkUp)
}

func TestGetEthernetPorts(t *testing.T) {
	ctx, err := testutil.NewTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.TearDown()

	ports, err := ctx.Unity.GetEthernetPorts()
	assert.Nil(t, err)

	assert.Equal(t, 8, len(ports))

	expectedId := []string{
		"spa_eth2", "spa_eth3", "spa_mgmt", "spa_srm",
		"spb_eth2", "spb_eth3", "spb_mgmt", "spb_srm",
	}
	expectedHeath := []int{5, 5, 5, 5, 5, 5, 5, 5}
	expectedName := []string{
		"SP A Ethernet Port 2", "SP A Ethernet Port 3", "SP A Management Port",
		"SP A Sync Replication Management Port",
		"SP B Ethernet Port 2", "SP B Ethernet Port 3", "SP B Management Port",
		"SP B Sync Replication Management Port",
	}
	expectedIsLinkUp := []bool{true, false, true, true, true, false, true, true}
	for i, port := range ports {
		assert.Equal(t, expectedId[i], port.Id)
		assert.Equal(t, expectedHeath[i], port.Health.Value)
		assert.Equal(t, expectedName[i], port.Name)
		assert.Equal(t, expectedIsLinkUp[i], port.IsLinkUp)
	}
}
