package gounity

import (
	"testing"
)

func TestGetPoolByID(t *testing.T) {
	u, _ := NewUnity("10.245.101.39", "admin", "Password123!", true)
	pool, _ := u.GetPoolByID("pool_1")
	t.Errorf("pool: %v", pool)
}
