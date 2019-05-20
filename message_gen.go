// DO NOT EDIT.
// GENERATED by go:generate at 2019-05-20 13:26:16.598163 +0000 UTC.
package gounity

// Message defines `message` type.
type Message struct {
	Resource

	Severity       SeverityEnum        `json:"severity"`
	ErrorCode      uint32              `json:"errorCode"`
	Created        string              `json:"created"`
	HttpStatusCode uint32              `json:"httpStatusCode"`
	Messages       []*LocalizedMessage `json:"messages"`
	MessageArgs    []string            `json:"messageArgs"`
}
