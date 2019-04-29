package gounity

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

const (
	// HeaderKeyContentType is key name of `Content-Type`.
	HeaderKeyContentType = "Content-Type"
	// HeaderValueContentTypeJSON is `json` value of `Content-Type`.
	HeaderValueContentTypeJSON = "application/json"
	// HeaderValueContentTypeBinaryOctetStream is `binary` value of `Content-Type`.
	HeaderValueContentTypeBinaryOctetStream = "binary/octet-stream"

	emcCsrfTokenName = "EMC-CSRF-TOKEN"
)

type restClient struct {
	http      *http.Client
	host      string
	username  string
	password  string
	authToken string
	csrfToken string
	traceHttp bool
}

type restClientOptions struct {
	insecure  bool
	traceHttp bool
}

// NewRestClientOptions returns a rest client option for creating rest client.
func NewRestClientOptions(insecure, traceHttp bool) *restClientOptions {
	return &restClientOptions{insecure: insecure, traceHttp: traceHttp}
}

// NewRestClient returns a new REST client to Unity.
func NewRestClient(
	ctx context.Context, host, username, password string, opts *restClientOptions,
) (*restClient, error) {

	if host == "" {
		return nil, errors.New("missing host")
	}

	cookieJar, err := cookiejar.New(
		&cookiejar.Options{PublicSuffixList: publicsuffix.List},
	)
	if err != nil {
		return nil, err
	}

	c := &restClient{
		http: &http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: opts.insecure,
				},
			},
		},
		host: host, username: username, password: password,
		traceHttp: opts.traceHttp,
	}

	return c, nil
}

func (c *restClient) pingPong(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	body io.Reader,
) (*http.Response, error) {

	msg := newMessage().withFields(
		map[string]interface{}{
			"method":  method,
			"path":    path,
			"headers": headers,
			"body":    body,
		},
	)
	var err error

	urlParts := []string{c.host}
	if len(path) > 0 {
		if path[0] == '/' {
			path = path[1:]
		}
		urlParts = append(urlParts, path)
	}

	fullURL, err := url.Parse(strings.Join(urlParts, "/"))
	if err != nil {
		return nil, errors.Wrapf(
			err, "parse url failed: %s", msg.withField("urlParts", urlParts),
		)
	}

	req, err := http.NewRequest(method, fullURL.String(), body)
	if err != nil {
		return nil, errors.Wrapf(err, "new request failed: %s", msg)
	}

	for header, value := range headers {
		req.Header.Set(header, value)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("X-EMC-REST-CLIENT", "true")
	req.Header.Set(emcCsrfTokenName, c.csrfToken)

	if c.traceHttp {
		traceRequest(ctx, req)
	}
	req = req.WithContext(ctx)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "http request failed: %s", msg)
	}
	if c.traceHttp {
		traceResponse(ctx, resp)
	}
	return resp, nil
}

func (c *restClient) do(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	body interface{},
) (*http.Response, error) {

	msg := newMessage().withFields(
		map[string]interface{}{
			"method":  method,
			"path":    path,
			"headers": headers,
			"body":    body,
		},
	)
	var bodyCache bytes.Buffer
	var reader io.Reader
	var readerAgain io.Reader
	contentType := ""
	if r, ok := body.(io.ReadCloser); ok {
		reader = io.TeeReader(r, &bodyCache)
		defer r.Close()
		readerAgain = &bodyCache
		contentType = HeaderValueContentTypeBinaryOctetStream

	} else if body != nil {
		bodyBuffer := &bytes.Buffer{}
		enc := json.NewEncoder(bodyBuffer)
		if err := enc.Encode(body); err != nil {
			return nil, errors.Wrapf(err, "encode request body failed: %s", msg)
		}
		reader = io.TeeReader(bodyBuffer, &bodyCache)
		readerAgain = &bodyCache
		contentType = HeaderValueContentTypeJSON
	}

	if contentType != "" {
		if headers == nil {
			headers = map[string]string{}
		}
		headers[HeaderKeyContentType] = contentType
	}

	if resp, err := c.pingPong(
		ctx, method, path, headers, reader,
	); resp != nil && resp.StatusCode != http.StatusUnauthorized {

		// Update csrf token only when getting unauthorized response
		return resp, err
	}

	if err := c.updateCsrf(ctx); err != nil {
		return nil, errors.Wrapf(err, "try to update csrf token failed: %s", msg)
	}

	return c.pingPong(ctx, method, path, headers, readerAgain)
}

func (c *restClient) updateCsrf(ctx context.Context) error {
	logrus.Debug("updating csrf token")

	resp, err := c.pingPong(ctx, http.MethodGet, "/api/types/system/instances", nil, nil)
	if err != nil {
		return errors.Wrapf(err, "update csrf token failed")
	}

	csrfToken := resp.Header.Get(emcCsrfTokenName)
	if c.csrfToken != csrfToken {
		c.csrfToken = csrfToken
		logrus.Debug("csrf token updated")
	}
	return nil
}

// DoWithHeaders sends a REST request with headers.
func (c *restClient) DoWithHeaders(
	ctx context.Context, method, path string,
	headers map[string]string, body, resp interface{},
) error {

	msg := newMessage().withFields(
		map[string]interface{}{
			"host":    c.host,
			"method":  method,
			"path":    path,
			"headers": headers,
			"body":    body,
		},
	)

	rawResp, err := c.do(ctx, method, path, headers, body)
	if err != nil {
		return errors.Wrapf(err, "http request with headers failed: %s", msg)
	}
	defer rawResp.Body.Close()

	switch {
	case rawResp == nil:
		return nil
	case rawResp.StatusCode >= 200 && rawResp.StatusCode <= 299:
		if resp == nil {
			return nil
		}
		dec := json.NewDecoder(rawResp.Body)
		if err = dec.Decode(resp); err != nil && err != io.EOF {
			return errors.Wrapf(err, "unable to decode response into %+v: %s", resp, msg)
		}
	default:
		unityError, err := ParseUnityError(rawResp.Body)
		if err != nil {
			return errors.Wrapf(
				err,
				"unknown error: %s", msg.withField("status code", rawResp.StatusCode),
			)
		}
		return unityError
	}
	return nil
}

// Do sends a REST request.
func (c *restClient) Do(
	ctx context.Context, method, path string, body, resp interface{},
) error {
	return c.DoWithHeaders(ctx, method, path, nil, body, resp)
}

// Get sends a REST request via GET method.
func (c *restClient) Get(
	ctx context.Context, path string, headers map[string]string, resp interface{},
) error {
	return c.DoWithHeaders(ctx, http.MethodGet, path, headers, nil, resp)
}

// Post sends a REST request via POST method.
func (c *restClient) Post(
	ctx context.Context, path string, headers map[string]string, body, resp interface{},
) error {
	return c.DoWithHeaders(ctx, http.MethodPost, path, headers, body, resp)
}

// Delete sends a REST request via DELETE method.
func (c *restClient) Delete(
	ctx context.Context, path string, headers map[string]string, body, resp interface{},
) error {
	return c.DoWithHeaders(ctx, http.MethodDelete, path, headers, nil, resp)
}
