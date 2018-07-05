package gounity

import "encoding/json"

type UnityErrorMessage struct {
	Message string `json:"en-US"`
}

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

type StorageResource struct {
	ID string `json:"id"`
}

type instanceResp struct {
	Content json.RawMessage `json:"content"`
}

type collectionResp struct {
	Entries []*instanceResp `json:"entries"`
}

type Pool struct {
	Unity       *Unity `json:"-"`
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SizeFree    uint64 `json:"sizeFree,omitempty"`
	SizeTotal   uint64 `json:"sizeTotal,omitempty"`
	SizeUsed    uint64 `json:"sizeUsed,omitempty"`
}

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

type Health struct {
	Value          int      `json:"value"`
	DescriptionIds []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

type HostLUNAccessEnum int

const (
	HostLUNAccess_NO_ACCESS HostLUNAccessEnum = iota
	HostLUNAccess_PRODUCTION
	HostLUNAccess_SNAPSHOT
	HostLUNAccess_BOTH
	HostLUNAccess_MIXED // TODO(ryan) Mixed = 0xffff
)

type Host struct {
	Unity       *Unity
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Health      *Health `json:"health,omitempty"`
	Description string  `json:"description"`
	OsType      string  `json:"osType"`
}

type BlockHostAccess struct {
	Host       *Host             `json:"host,omitempty"`
	AccessMask HostLUNAccessEnum `json:"accessMask"`
}

type HostLUNTypeEnum int

const (
	HostLUNType_UNKNOWN HostLUNTypeEnum = iota
	HostLUNType_LUN
	HostLUNType_Snap
)

type HostLUN struct {
	ID            string          `json:"id"`
	Host          *Host           `json:"host"`
	Type          HostLUNTypeEnum `json:"type"`
	Hlu           uint16          `json:"hlu"`
	LUN           *LUN            `json:"lun"`
	IsReadOnly    bool            `json:"isReadOnly"`
	IsDefaultSnap bool            `json:"isDefaultSnap"`
}
