package gounity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

//go:generate go run cmd/modelgen/modelgen.go -w cmd/modelgen
//go:generate go fmt .

var (
	// Reading from env variable makes it easy to trace http without modifying
	// the source code.
	traceHttp, _ = strconv.ParseBool(os.Getenv("GOUNITY_TRACEHTTP"))
)

// RestClient acts as a REST client.
type RestClient interface {
	Do(
		ctx context.Context,
		method, path string, body, resp interface{},
	) error

	Get(
		ctx context.Context,
		path string, headers map[string]string, resp interface{},
	) error

	Post(
		ctx context.Context,
		path string, headers map[string]string, body, resp interface{},
	) error

	Delete(
		ctx context.Context,
		path string, headers map[string]string, body, resp interface{},
	) error
}

// UnityConnector defines the interface to storage system.
type UnityConnector interface {
	GetInstanceById(
		resType, id, fields string, instance interface{},
	) error

	GetInstanceByName(
		resType, name, fields string, instance interface{},
	) error

	GetCollection(
		resType, fields string, filter *filter,
	) ([]*instanceResp, error)

	PostOnType(
		typeName, action string, body map[string]interface{},
	) (string, error)

	PostOnInstance(
		typeName, resId, action string, body map[string]interface{},
	) error

	PoolOperatorGen

	LunOperatorGen

	HostOperatorGen

	HostLUNOperator

	NasServerOperatorGen

	FilesystemOperatorGen

	NfsShareOperatorGen
}

// Unity defines the connection to Unity system.
type Unity struct {
	client RestClient
}

// NewUnity creates a connection to a Unity system.
func NewUnity(
	mgmtIp, username, password string, insecure bool,
) (*Unity, error) {

	fields := map[string]interface{}{
		"mgmtIp":    mgmtIp,
		"insecure":  insecure,
		"traceHttp": traceHttp,
	}
	logger := log.WithFields(fields)
	msg := newMessage().withFields(fields)

	logger.Debug("gounity connection initializing")

	if mgmtIp == "" {
		return nil, errors.Errorf("mgmtIp is required: %s", msg)
	}

	opts := NewRestClientOptions(insecure, traceHttp)

	host := fmt.Sprintf("%s://%s", "https", mgmtIp)

	restClient, err := NewRestClient(context.Background(), host, username, password, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create rest client: %s", msg)
	}

	unity := &Unity{client: restClient}
	return unity, nil
}

type instanceResp struct {
	Content json.RawMessage `json:"content"`
}

// GetInstanceById queries instance via id.
func (u *Unity) GetInstanceById(
	resType, id, fields string, instance interface{},
) error {
	msg := newMessage().withFields(
		map[string]interface{}{
			"resourceType": resType,
			"id":           id,
			"fields":       fields,
		},
	)
	resp := &instanceResp{}
	if err := u.client.Get(
		context.Background(), queryInstanceUrl(resType, id, fields), nil, resp,
	); err != nil {
		return errors.Wrapf(err, "query instance failed: %s", msg)
	}
	if err := json.Unmarshal(resp.Content, instance); err != nil {
		return errors.Wrapf( err, "decode to %v failed: %s", instance, msg)
	}
	return nil
}

// GetInstanceByName queries instance via name.
func (u *Unity) GetInstanceByName(
	resType, name, fields string, instance interface{},
) error {
	return u.GetInstanceById(resType, "name:"+name, fields, instance)
}

type collectionResp struct {
	Entries []*instanceResp `json:"entries"`
}

// GetCollection queries instance collection.
func (u *Unity) GetCollection(
	resType, fields string, filter *filter,
) ([]*instanceResp, error) {

	msg := newMessage().withFields(
		map[string]interface{}{
			"resourceType": resType,
			"fields":       fields,
			"filter":       filter,
		},
	)
	resp := &collectionResp{}
	if err := u.client.Get(
		context.Background(), queryCollectionUrl(resType, fields, filter), nil, resp,
	); err != nil {
		return nil, errors.Wrapf(err, "query collection failed: %s", msg)
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

// PostOnType sends POST request on resource type.
func (u *Unity) PostOnType(
	typeName, action string, body map[string]interface{},
) (string, error) {

	msg := newMessage().withFields(
		map[string]interface{}{
			"typeName": typeName,
			"action":   action,
			"body":     body,
		},
	)
	resp := &storageResourceCreateResp{}
	if err := u.client.Post(
		context.Background(), postTypeUrl(typeName, action), nil, body, resp,
	); err != nil {
		return "", errors.Wrapf(err, "post on type failed: %s", msg)
	}

	return resp.Content.StorageResource.Id, nil
}

// PostOnInstance sends POST request on resource instance.
func (u *Unity) PostOnInstance(
	typeName, resId, action string, body map[string]interface{},
) error {

	msg := newMessage().withFields(
		map[string]interface{}{
			"typeName": typeName,
			"resId":    resId,
			"action":   action,
			"body":     body,
		},
	)
	if err := u.client.Post(
		context.Background(), postInstanceUrl(typeName, resId, action), nil, body, nil,
	); err != nil {
		return errors.Wrapf(err, "post on instance failed: %s", msg)
	}
	return nil
}
