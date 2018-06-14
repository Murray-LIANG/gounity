package gounity

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	traceHTTP, _ = strconv.ParseBool(os.Getenv("GOUNITY_TRACEHTTP"))
)

type ConnectionInfo struct {
	MgmtIP   string
	Username string
	Password string
}

type Unity struct {
	client         RestClient
	connectionInfo *ConnectionInfo
}

// UnitySystem represents a Unity storage
type Storage interface {
	GetPools() ([]*Pool, error)
	GetPoolById(id string) (*Pool, error)
	GetPoolByName(name string) (*Pool, error)

	GetLUNs() ([]*LUN, error)
	GetLUNById(id string) (*LUN, error)
	GetLUNByName(name string) (*LUN, error)
}

func newErrorWithFields(fields map[string]interface{}, message string,
	inner error) error {

	if fields == nil {
		fields = map[string]interface{}{}
	}

	if inner != nil {
		fields["error"] = inner
	}

	kvStrs := []string{}
	for k, v := range fields {
		kvStrs = append(kvStrs, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Errorf("%s %s", message, strings.Join(kvStrs, ","))
}

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
		return nil, newErrorWithFields(fields, "mgmtIP is required", nil)
	}

	opts := RestClientOptions{
		Insecure:  insecure,
		TraceHTTP: traceHTTP,
	}

	host := fmt.Sprintf("%s://%s", "https", mgmtIP)

	restClient, err := NewRestClient(context.Background(), host, username, password, opts)
	if err != nil {
		logger.Error("failed to create rest client")
		return nil, newErrorWithFields(fields, "failed to create rest client", err)
	}

	unity := &Unity{
		client: restClient,
		connectionInfo: &ConnectionInfo{
			MgmtIP: mgmtIP,
		},
	}
	return unity, nil
}
