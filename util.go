package gounity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	reqTraceHead = `
--------------------
GOUNITY HTTP REQUEST
--------------------
`
	respTraceHead = `
---------------------
GOUNITY HTTP RESPONSE
---------------------
`
)

func dumpBody(header *http.Header) bool {
	return header.Get(HeaderKeyContentType) != HeaderValueContentTypeBinaryOctetStream
}

func traceRequest(ctx context.Context, req *http.Request) {

	reqBuffer, err := httputil.DumpRequest(req, dumpBody(&req.Header))
	if err != nil {
		return
	}
	log.Debug(strings.Join([]string{reqTraceHead, string(reqBuffer)}, "\n"))
}

func traceResponse(ctx context.Context, resp *http.Response) {

	respBuffer, err := httputil.DumpResponse(resp, dumpBody(&resp.Header))
	if err != nil {
		return
	}
	log.Debug(strings.Join([]string{respTraceHead, string(respBuffer)}, "\n"))
}

const (
	typeStorageResource    = "storageResource"
	actionCreateLun        = "createLun"
	actionModifyLun        = "modifyLun"
	actionCreateFilesystem = "createFilesystem"
	actionModifyFilesystem = "modifyFilesystem"
)

const (
	pathAPITypes     = "api/types"
	pathAPIInstances = "api/instances"
)

func buildUrl(baseURL, fields string, filter *filter) string {
	queryParams := map[string]string{"compact": "true", "fields": fields}
	if filter != nil {
		queryParams["filter"] = filter.String()
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		log.WithError(err).WithField("baseURL", baseURL).Error("failed to parse url")
		return ""
	}
	query := u.Query()
	for k, v := range queryParams {
		query.Add(k, v)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func queryCollectionUrl(res, fields string, filter *filter) string {
	return buildUrl(strings.Join([]string{pathAPITypes, res, "instances"}, "/"), fields,
		filter)
}

func queryInstanceUrl(res, id, fields string) string {
	return buildUrl(strings.Join([]string{pathAPIInstances, res, id}, "/"), fields, nil)
}

func postTypeUrl(typeName, action string) string {
	return strings.Join([]string{pathAPITypes, typeName, "action", action}, "/")
}

func postInstanceUrl(typeName, resId, action string) string {
	return strings.Join([]string{pathAPIInstances, typeName, resId, "action", action}, "/")
}

// UnityErrorMessage defines the error message struct returned by Unity.
type unityErrorMessage struct {
	Message string `json:"en-US"`
}

// UnityError defines the error struct returned by Unity.
type unityError struct {
	ErrorCode      int                 `json:"errorCode"`
	HttpStatusCode int                 `json:"httpStatusCode"`
	Messages       []unityErrorMessage `json:"messages"`
	Message        string
}

type unityErrorResp struct {
	Error *unityError `json:"error,omitempty"`
}

func (e *unityError) Error() string {
	return fmt.Sprintf(
		"error from unity, status code: %v, error code: %v, message: %v",
		e.HttpStatusCode, e.ErrorCode, e.Message,
	)
}

var ErrUnableParseRespToError = errors.New("unable parse response body to unity error")

func ParseUnityError(reader io.Reader) (error, error) {
	resp := &unityErrorResp{}
	if err := json.NewDecoder(reader).Decode(resp); err != nil {
		return nil, err
	}
	if resp.Error == nil {
		// not a `unity error` json in reader
		return nil, ErrUnableParseRespToError
	}

	respError := resp.Error
	respError.Message = respError.Messages[0].Message
	return respError, nil
}

func gbToBytes(gb uint64) uint64 {
	return gb * 1024 * 1024 * 1024
}

type idRepresent struct {
	Id string `json:"id"`
}
