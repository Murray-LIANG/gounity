package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Murray-LIANG/gounity"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func parseResourceName(url string) (string, error) {
	// url is like `/api/xxx` without host name/IP and port
	parts := strings.Split(url, "/")
	// i.e. `/api/instances/lun/sv_1` and `/api/types/lun/instances`
	if len(parts) < 4 {
		msg := "cannot find resource name in url"
		log.WithField("url", url).Error(msg)
		return "", errors.New(msg)
	}
	return parts[3], nil
}

type mockIndex struct {
	Indices []struct {
		URL      string      `json:"url"`
		Body     interface{} `json:"body"`
		Response string      `json:"response"`
	} `json:"indices"`
}

func mockServerHandler(resp http.ResponseWriter, req *http.Request) {
	reqURL := req.URL.String() // url is like `/api/xxx` without host IP and port
	unescapeURL, err := url.QueryUnescape(reqURL)
	if err != nil {
		log.WithError(err).Error("failed to get unescape url")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	resource, err := parseResourceName(unescapeURL)
	if err != nil {
		log.WithError(err).WithField("unescapeURL", unescapeURL).Error(
			"failed to get resource type from unescape url")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Error("failed to get current working directory")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	resourceDataDir := filepath.Join(cwd, "testdata", resource)
	indexFilePath := filepath.Join(resourceDataDir, "index.json")

	indicesBytes, err := ioutil.ReadFile(indexFilePath)
	if err != nil {
		log.WithError(err).WithField("filepath", indexFilePath).Error(
			"failed to read index.json file")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	dec := json.NewDecoder(bytes.NewReader(indicesBytes))
	dec.UseNumber()
	var indices mockIndex
	if err = dec.Decode(&indices); err != nil {
		log.WithError(err).WithField("filepath", indexFilePath).Error(
			"failed to parse index.json file")
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	dec = json.NewDecoder(req.Body)
	dec.UseNumber()
	var reqBody interface{}
	if err = dec.Decode(&reqBody); err != nil && err != io.EOF {
		log.WithError(err).Error("failed to decode request body")
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	respFileName := ""
	for _, index := range indices.Indices {
		if index.URL == unescapeURL {
			log.WithField("requestBody", reqBody).WithField(
				"mockBody", index.Body).Debug("check if request and mock body matches")
			if reflect.DeepEqual(reqBody, index.Body) {
				respFileName = index.Response
				break
			}
		}
	}
	if respFileName == "" {
		log.WithFields(map[string]interface{}{
			"filepath":    indexFilePath,
			"urlAfterIP":  unescapeURL,
			"requestBody": reqBody,
		}).Error("failed to find response for request in index.json")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	respFilePath := filepath.Join(resourceDataDir, respFileName)
	if respBytes, err := ioutil.ReadFile(respFilePath); err != nil {
		log.WithField("respFilePath", respFilePath).Error(
			"failed to read the response file")
		resp.WriteHeader(http.StatusNotFound)
	} else {
		mockError, err := gounity.ParseUnityError(bytes.NewReader(respBytes))
		if mockError != nil {
			if code := gounity.GetUnityErrorStatusCode(mockError); code != -1 {
				resp.WriteHeader(code)
			}
		} else if err != nil && err != io.EOF &&
			err != gounity.ErrUnableParseRespToError {

			resp.WriteHeader(http.StatusNotFound)
		}
		resp.Write(respBytes)
	}
}

func setupMockServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(mockServerHandler))
}

type testContext struct {
	mockServer *httptest.Server
	context    context.Context
	RestClient gounity.RestClient
	Unity      gounity.UnityConnector
}

func NewTestContext() (*testContext, error) {
	mockServer := setupMockServer()
	ctx := context.Background()
	restClient, err := gounity.NewRestClient(
		ctx, mockServer.URL, "", "", gounity.NewRestClientOptions(true, true),
	)
	if err != nil {
		return nil, err
	}
	unity, err := gounity.NewUnity(extractIp(mockServer.URL), "", "", true)
	if err != nil {
		return nil, err
	}
	return &testContext{
		mockServer: mockServer, context: ctx, RestClient: restClient, Unity: unity,
	}, nil
}

func (c *testContext) TearDown() {
	c.mockServer.Close()
}

func extractIp(url string) string {
	return strings.Split(url, "/")[2]
}
