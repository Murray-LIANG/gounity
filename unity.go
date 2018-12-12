package gounity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	traceHttp, _ = strconv.ParseBool(os.Getenv("GOUNITY_TRACEHTTP"))
)

// Unity defines the connection to Unity system.
type Unity struct {
	client RestClient
}

// Storage defines a Unity system.
type Storage interface {
	GetPools() ([]*Pool, error)
	GetPoolById(id string) (*Pool, error)
	GetPoolByName(name string) (*Pool, error)

	GetLuns() ([]*Lun, error)
	GetLunById(id string) (*Lun, error)
	GetLunByName(name string) (*Lun, error)
}

// Health defines Unity corresponding `health` type.
type Health struct {
	Value          int      `json:"value"`
	DescriptionIds []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

// NewUnity creates a connection to a Unity system.
func NewUnity(mgmtIp, username, password string, insecure bool) (*Unity, error) {

	fields := map[string]interface{}{
		"mgmtIp":    mgmtIp,
		"insecure":  insecure,
		"traceHttp": traceHttp,
	}
	logger := log.WithFields(fields)

	logger.Debug("gounity connection initializing")

	if mgmtIp == "" {
		logger.Error("mgmtIp is required")
		return nil, newGounityError("mgmtIp is required").withFields(fields)
	}

	opts := RestClientOptions{
		Insecure:  insecure,
		TraceHttp: traceHttp,
	}

	host := fmt.Sprintf("%s://%s", "https", mgmtIp)

	restClient, err := NewRestClient(context.Background(), host, username, password, opts)
	if err != nil {
		logger.Error("failed to create rest client")
		return nil, newGounityError(
			"failed to create rest client").withFields(fields).withError(err)
	}

	unity := &Unity{client: restClient}
	return unity, nil
}

type instanceResp struct {
	Content json.RawMessage `json:"content"`
}

func (u *Unity) getInstanceById(resType, id, fields string, instance interface{}) error {
	resp := &instanceResp{}
	if err := u.client.Get(context.Background(), queryInstanceUrl(resType, id, fields),
		nil, resp); err != nil {
		return err
	}
	if err := json.Unmarshal(resp.Content, instance); err != nil {
		return err
	}
	return nil
}

func (u *Unity) getInstanceByName(
	resType, name, fields string, instance interface{},
) error {

	return u.getInstanceById(resType, "name:"+name, fields, instance)
}

type filter []string

func NewFilter(f string) *filter {
	return &filter{f}
}

func (f *filter) And(andFilter string) *filter {
	newFilter := append(*f, "and")
	newFilter = append(newFilter, andFilter)
	return &newFilter
}

func (f *filter) String() string {
	return strings.Join(*f, " ")
}

type collectionResp struct {
	Entries []*instanceResp `json:"entries"`
}

func (u *Unity) getCollection(
	resType, fields string, filter *filter,
) ([]*instanceResp, error) {

	resp := &collectionResp{}
	if err := u.client.Get(
		context.Background(),
		queryCollectionUrl(resType, fields, filter),
		nil,
		resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
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

func (u *Unity) postOnType(
	typeName, action string, body map[string]interface{},
) (string, error) {

	resp := &storageResourceCreateResp{}
	if err := u.client.Post(context.Background(), postTypeUrl(typeName, action),
		nil, body, resp); err != nil {
		return "", err
	}

	return resp.Content.StorageResource.Id, nil
}

func (u *Unity) postOnInstance(
	typeName, resId, action string, body map[string]interface{},
) error {

	if err := u.client.Post(context.Background(),
		postInstanceUrl(typeName, resId, action), nil, body, nil); err != nil {
		return err
	}
	return nil
}
