package gounity

import (
	"strings"
)

var (
	typeNameLun   = "lun"
	typeFieldsLun = strings.Join([]string{
		// "compressionPercent",
		// "compressionSizeSaved",
		// "currentNode",
		// "defaultNode",
		"description",
		"health",
		"hostAccess",
		"id",
		// "ioLimitPolicy.id",
		// "isCompressionEnabled",
		// "isReplicationDestination",
		// "isSnapSchedulePaused",
		"isThinEnabled",
		"metadataSize",
		"metadataSizeAllocated",
		"name",
		// "perTierSizeUsed",
		"pool.id",
		"sizeAllocated",
		"sizeTotal",
		"sizeUsed",
		"snapCount",
		// "snapSchedule.id",
		"snapWwn",
		"snapsSize",
		"snapsSizeAllocated",
		// "storageResource.id",
		// "tieringPolicy",
		// "type",
		"wwn",
	}, ",")
)

// Lun defines Unity corresponding `lun` type.
type Lun struct {
	Resource
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

// BlockHostAccess defines Unity corresponding `blockHostAccess` type.
type BlockHostAccess struct {
	Host       *Host             `json:"host,omitempty"`
	AccessMask HostLunAccessEnum `json:"accessMask"`
}

// HostLunAccessEnum defines Unity corresponding `HostLunAccessEnum` enumeration.
type HostLunAccessEnum int

const (
	// HostLunAccessNoAccess defines `NoAccess` value of HostLunAccessEnum.
	HostLunAccessNoAccess HostLunAccessEnum = iota

	// HostLunAccessProduction defines `Production` value of HostLunAccessEnum.
	HostLunAccessProduction

	// HostLunAccessSnapshot defines `Snapshot` value of HostLunAccessEnum.
	HostLunAccessSnapshot

	// HostLunAccessBoth defines `Both` value of HostLunAccessEnum.
	HostLunAccessBoth

	// HostLunAccessMixed defines `Mixed` value of HostLunAccessEnum.
	HostLunAccessMixed // TODO(ryan) Mixed = 0xffff
)

//go:generate ./gen_resource.sh resource_tmpl.go lun_gen.go Lun
