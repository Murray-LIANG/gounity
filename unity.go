package gounity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	traceHTTP, _ = strconv.ParseBool(os.Getenv("GOUNITY_TRACEHTTP"))
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

	GetLUNs() ([]*LUN, error)
	GetLUNById(id string) (*LUN, error)
	GetLUNByName(name string) (*LUN, error)
}

// NewUnity creates a connection to a Unity system.
func NewUnity(mgmtIP, username, password string, insecure bool) (*Unity, error) {

	fields := map[string]interface{}{
		"mgmtIp":    mgmtIP,
		"insecure":  insecure,
		"traceHTTP": traceHTTP,
	}
	logger := log.WithFields(fields)

	logger.Debug("gounity connection initializing")

	if mgmtIP == "" {
		logger.Error("mgmtIP is required")
		return nil, newGounityError("mgmtIP is required").withFields(fields)
	}

	opts := RestClientOptions{
		Insecure:  insecure,
		TraceHTTP: traceHTTP,
	}

	host := fmt.Sprintf("%s://%s", "https", mgmtIP)

	restClient, err := NewRestClient(context.Background(), host, username, password, opts)
	if err != nil {
		logger.Error("failed to create rest client")
		return nil, newGounityError(
			"failed to create rest client").withFields(fields).withError(err)
	}

	unity := &Unity{client: restClient}
	return unity, nil
}

func setUnity(instancePtr reflect.Value, unity *Unity) {
	if instancePtr.Kind() != reflect.Ptr {
		log.WithField("instancePtr",
			instancePtr).Debug("`instancePtr` is not a pointer, skip setting field unity")
		return
	}
	instStruct := instancePtr.Elem()
	if instStruct.Kind() != reflect.Struct {
		log.WithField("instStruct",
			instancePtr).Debug("`instStruct` is not a struct, skip setting field unity")
		return
	}
	field := instStruct.FieldByName("Unity")
	if !field.IsValid() {
		log.Debug("field `Unity` is not found, skip setting field unity")
		return
	}
	if field.Kind() != reflect.TypeOf(unity).Kind() {
		log.WithField("fieldKind", field.Kind()).Debug("field `Unity` is type `*Unity`")
		return
	}
	field.Set(reflect.ValueOf(unity))
}

func (u *Unity) getInstanceByID(resType, id, fields string, instance interface{}) error {
	resp := &instanceResp{}
	if err := u.client.Get(context.Background(), queryInstanceURL(resType, id, fields),
		nil, resp); err != nil {
		return err
	}
	if err := json.Unmarshal(resp.Content, instance); err != nil {
		return err
	}
	setUnity(reflect.ValueOf(instance), u)
	return nil
}

type filter []string

func newFilter(f string) *filter {
	return &filter{f}
}

func (f *filter) and(andFilter string) *filter {
	newFilter := append(*f, "and")
	newFilter = append(newFilter, andFilter)
	return &newFilter
}

func (f *filter) string() string {
	return strings.Join(*f, " ")
}

func (u *Unity) getCollection(resType, fields string, filter *filter,
	instanceType reflect.Type) (interface{}, error) {
	resp := &collectionResp{}
	if err := u.client.Get(context.Background(),
		queryCollectionURL(resType, fields, filter), nil,
		resp); err != nil {
		return nil, err
	}

	collection := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(instanceType)), 0,
		len(resp.Entries))
	for _, entry := range resp.Entries {
		instance := reflect.New(instanceType)
		if err := json.Unmarshal(entry.Content, instance.Interface()); err != nil {
			return nil, err
		}
		setUnity(instance, u)
		collection = reflect.Append(collection, instance)
	}
	return collection.Interface(), nil
}
