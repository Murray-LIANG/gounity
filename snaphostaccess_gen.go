// DO NOT EDIT.
// GENERATED by go:generate at 2019-05-20 13:26:16.605589 +0000 UTC.
package gounity

// SnapHostAccess defines `snapHostAccess` type.
type SnapHostAccess struct {
	Resource

	Host       *Host               `json:"host"`
	AccessMask SnapAccessLevelEnum `json:"accessMask"`
}
