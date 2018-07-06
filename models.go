package gounity

import "encoding/json"

// UnityErrorMessage defines the error message struct returned by Unity.
type UnityErrorMessage struct {
	Message string `json:"en-US"`
}

// UnityError defines the error struct returned by Unity.
type UnityError struct {
	ErrorCode      int                 `json:"errorCode"`
	HTTPStatusCode int                 `json:"httpStatusCode"`
	Messages       []UnityErrorMessage `json:"messages"`
	Message        string
}

type unityErrorResp struct {
	Error *UnityError `json:"error,omitempty"`
}

func (e *UnityError) Error() string {
	return e.Message
}

type storageResourceCreateResp struct {
	Content struct {
		StorageResource *StorageResource `json:"storageResource,omitempty"`
	} `json:"content"`
}

// StorageResource defines Unity corresponding storage resource(like pool, LUN .etc).
type StorageResource struct {
	ID string `json:"id"`
}

type instanceResp struct {
	Content json.RawMessage `json:"content"`
}

type collectionResp struct {
	Entries []*instanceResp `json:"entries"`
}

// Pool defines Unity corresponding `pool` type.
type Pool struct {
	Unity       *Unity `json:"-"`
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SizeFree    uint64 `json:"sizeFree,omitempty"`
	SizeTotal   uint64 `json:"sizeTotal,omitempty"`
	SizeUsed    uint64 `json:"sizeUsed,omitempty"`
}

// LUN defines Unity corresponding `lun` type.
type LUN struct {
	Unity                 *Unity             `json:"-"`
	Description           string             `json:"description"`
	Health                *Health            `json:"health,omitempty"`
	HostAccess            []*BlockHostAccess `json:"hostAccess,omitempty"`
	ID                    string             `json:"id"`
	IsThinEnabled         bool               `json:"isThinEnabled"`
	MetadataSize          uint64             `json:"metadataSize"`
	MetadataSizeAllocated uint64             `json:"metadataSizeAllocated"`
	Name                  string             `json:"name"`
	Pool                  *Pool              `json:"pool,omitempty"`
	SizeAllocated         uint64             `json:"sizeAllocated"`
	SizeTotal             uint64             `json:"sizeTotal"`
	SizeUsed              uint64             `json:"sizeUsed"`
	SnapCount             uint32             `json:"snapCount"`
	SnapWwn               string             `json:"snapWwn"`
	SnapsSize             uint64             `json:"snapsSize"`
	SnapsSizeAllocated    uint64             `json:"snapsSizeAllocated"`
	Wwn                   string             `json:"wwn"`
}

// Health defines Unity corresponding `health` type.
type Health struct {
	Value          int      `json:"value"`
	DescriptionIds []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

// HostLUNAccessEnum defines Unity corresponding `HostLUNAccessEnum` enumeration.
type HostLUNAccessEnum int

const (
	// HostLUNAccessNoAccess defines `NoAccess` value of HostLUNAccessEnum.
	HostLUNAccessNoAccess HostLUNAccessEnum = iota

	// HostLUNAccessProduction defines `Production` value of HostLUNAccessEnum.
	HostLUNAccessProduction

	// HostLUNAccessSnapshot defines `Snapshot` value of HostLUNAccessEnum.
	HostLUNAccessSnapshot

	// HostLUNAccessBoth defines `Both` value of HostLUNAccessEnum.
	HostLUNAccessBoth

	// HostLUNAccessMixed defines `Mixed` value of HostLUNAccessEnum.
	HostLUNAccessMixed // TODO(ryan) Mixed = 0xffff
)

// Host defines Unity corresponding `host` type.
type Host struct {
	Unity       *Unity
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
	OsType      string  `json:"osType"`
}

// BlockHostAccess defines Unity corresponding `blockHostAccess` type.
type BlockHostAccess struct {
	Host       *Host             `json:"host,omitempty"`
	AccessMask HostLUNAccessEnum `json:"accessMask"`
}

// HostLUNTypeEnum defines Unity corresponding `HostLUNTypeEnum` enumeration.
type HostLUNTypeEnum int

const (
	// HostLUNTypeUnknown defines `Unknown` value of HostLUNTypeEnum.
	HostLUNTypeUnknown HostLUNTypeEnum = iota

	// HostLUNTypeLUN defines `LUN` value of HostLUNTypeEnum.
	HostLUNTypeLUN

	// HostLUNTypeSnap defines `Snap` value of HostLUNTypeEnum.
	HostLUNTypeSnap
)

// HostLUN defines Unity corresponding `HostLUN` type.
type HostLUN struct {
	ID            string          `json:"id"`
	Host          *Host           `json:"host"`
	Type          HostLUNTypeEnum `json:"type"`
	Hlu           uint16          `json:"hlu"`
	LUN           *LUN            `json:"lun"`
	IsReadOnly    bool            `json:"isReadOnly"`
	IsDefaultSnap bool            `json:"isDefaultSnap"`
}
