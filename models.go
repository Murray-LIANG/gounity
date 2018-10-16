package gounity

import "encoding/json"

// UnityErrorMessage defines the error message struct returned by Unity.
type UnityErrorMessage struct {
	Message string `json:"en-US"`
}

// UnityError defines the error struct returned by Unity.
type UnityError struct {
	ErrorCode      int                 `json:"errorCode"`
	HttpStatusCode int                 `json:"httpStatusCode"`
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

// StorageResource defines Unity corresponding storage resource(like pool, Lun .etc).
type StorageResource struct {
	Id string `json:"id"`
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
	Id          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SizeFree    uint64 `json:"sizeFree,omitempty"`
	SizeTotal   uint64 `json:"sizeTotal,omitempty"`
	SizeUsed    uint64 `json:"sizeUsed,omitempty"`
}

// Lun defines Unity corresponding `lun` type.
type Lun struct {
	Unity                 *Unity             `json:"-"`
	Description           string             `json:"description"`
	Health                *Health            `json:"health,omitempty"`
	HostAccess            []*BlockHostAccess `json:"hostAccess,omitempty"`
	Id                    string             `json:"id"`
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

// HostLunAccessEnum defines Unity corresponding `HostLunAccessEnum` enumeration.
type HostLunAccessEnum int

const (
	// HostLUNAccessNoAccess defines `NoAccess` value of HostLunAccessEnum.
	HostLUNAccessNoAccess HostLunAccessEnum = iota

	// HostLunAccessProduction defines `Production` value of HostLunAccessEnum.
	HostLunAccessProduction

	// HostLUNAccessSnapshot defines `Snapshot` value of HostLunAccessEnum.
	HostLUNAccessSnapshot

	// HostLUNAccessBoth defines `Both` value of HostLunAccessEnum.
	HostLUNAccessBoth

	// HostLUNAccessMixed defines `Mixed` value of HostLunAccessEnum.
	HostLUNAccessMixed // TODO(ryan) Mixed = 0xffff
)

// Host defines Unity corresponding `host` type.
type Host struct {
	Unity       *Unity
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
	OsType      string  `json:"osType"`
}

// BlockHostAccess defines Unity corresponding `blockHostAccess` type.
type BlockHostAccess struct {
	Host       *Host             `json:"host,omitempty"`
	AccessMask HostLunAccessEnum `json:"accessMask"`
}

// HostLunTypeEnum defines Unity corresponding `HostLunTypeEnum` enumeration.
type HostLunTypeEnum int

const (
	// HostLUNTypeUnknown defines `Unknown` value of HostLunTypeEnum.
	HostLUNTypeUnknown HostLunTypeEnum = iota

	// HostLUNTypeLUN defines `Lun` value of HostLunTypeEnum.
	HostLUNTypeLUN

	// HostLUNTypeSnap defines `Snap` value of HostLunTypeEnum.
	HostLUNTypeSnap
)

// HostLun defines Unity corresponding `HostLun` type.
type HostLun struct {
	Id            string          `json:"id"`
	Host          *Host           `json:"host"`
	Type          HostLunTypeEnum `json:"type"`
	Hlu           uint16          `json:"hlu"`
	Lun           *Lun            `json:"lun"`
	IsReadOnly    bool            `json:"isReadOnly"`
	IsDefaultSnap bool            `json:"isDefaultSnap"`
}
