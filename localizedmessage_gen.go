// DO NOT EDIT.
// GENERATED by go:generate at 2019-06-06 09:02:53.366181 +0000 UTC.
package gounity

// LocalizedMessage defines `localizedMessage` type.
type LocalizedMessage struct {
	Resource

	Locale  string `json:"locale"`
	Message string `json:"message"`
}
