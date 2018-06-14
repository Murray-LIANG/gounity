package gounity

type UnityErrorMessage struct {
	Message string `json:"en-US"`
}

type UnityError struct {
	ErrorCode      int64                `json:"errorCode"`
	HTTPStatusCode int                  `json:"httpStatusCode`
	Messages       []*UnityErrorMessage `json:"messages"`
	Message        string
}

type UnityErrorResp struct {
	Error *UnityError `json:"error"`
}

func (e UnityError) Error() string {
	return e.Message
}

type Health struct {
	Value          int      `json:"value"`
	DescriptionIds []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

type poolResp struct {
	Content Pool `json:"content"`
}

type Pool struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SizeFree    uint64 `json:"sizeFree"`
	SizeTotal   uint64 `json:"sizeTotal"`
	SizeUsed    uint64 `json:"sizeUsed"`
}

type lunResp struct {
	Content LUN `json:"content"`
}

type lunsResp struct {
	Entries []lunResp `json:"entries"`
}

type LUN struct {
	Description           string            `json:"description"`
	Health                Health            `json:"health"`
	HostAccess            []BlockHostAccess `json:"hostAccess"`
	ID                    string            `json:"id"`
	IsThinEnabled         bool              `json:"isThinEnabled"`
	MetadataSize          uint64            `json:"metadataSize"`
	MetadataSizeAllocated uint64            `json:"metadataSizeAllocated"`
	Name                  string            `json:"name"`
	Pool                  Pool              `json:"pool"`
	SizeAllocated         uint64            `json:"sizeAllocated"`
	SizeTotal             uint64            `json:"sizeTotal"`
	SizeUsed              uint64            `json:"sizeUsed"`
	SnapCount             uint32            `json:"snapCount"`
	SnapWwn               string            `json:"snapWwn"`
	SnapsSize             uint64            `json:"snapsSize"`
	SnapsSizeAllocated    uint64            `json:"snapsSizeAllocated"`
	Wwn                   string            `json:"wwn"`
}

type HostLUNAccessEnum int

const (
	NO_ACCESS HostLUNAccessEnum = iota
	PRODUCTION
	SNAPSHOT
	BOTH
	MIXED // TODO(ryan) Mixed = 0xffff
)

type Host struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Health      Health `json:"health"`
	Description string `json:"description"`
	OsType      string `json:"osType"`
	// TODO(ryan) add other attributes
}

type BlockHostAccess struct {
	Host       Host              `json:"host"`
	AccessMask HostLUNAccessEnum `json:"accessMask"`
	// TODO(ryan) add other attributes
}
