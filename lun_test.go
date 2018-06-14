package gounity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestLUN() *Unity {
	u, _ := NewUnity("10.245.101.39", "admin", "Password123!", true)
	return u
}

func TestGetLUNByID(t *testing.T) {
	u := setupTestLUN()
	lun, _ := u.GetLUNByID("sv_1")

	assert.Equal(t, "sv_1", lun.ID)
	assert.Equal(t, "pool_1", lun.Pool.ID)
}

func TestGetLUNs(t *testing.T) {
	u := setupTestLUN()
	luns, _ := u.GetLUNs()

	assert.Equal(t, 15, len(luns))
	ids := []string{}
	for _, lun := range luns {
		ids = append(ids, lun.ID)
	}
	// assert.Equal(t, "sv_1", lun.ID)
	// assert.Equal(t, "pool_1", lun.Pool.ID)
}
