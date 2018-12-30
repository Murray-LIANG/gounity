package gounity

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

const (
	// HeaderKeyContentType is key name of `Content-Type`.
	HeaderKeyContentType = "Content-Type"
	// HeaderValueContentTypeJSON is `json` value of `Content-Type`.
	HeaderValueContentTypeJSON = "application/json"
	// HeaderValueContentTypeBinaryOctetStream is `binary` value of `Content-Type`.
	HeaderValueContentTypeBinaryOctetStream = "binary/octet-stream"
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
		http: &http.Client{Jar: cookieJar},
		host: host, username: username, password: password,
		traceHttp: opts.traceHttp,
	}

	if opts.insecure {
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c, nil
}

func setDefaultContentType(
	header *http.Header, headers map[string]string, defaultValue string,
) {
	if v, ok := headers[HeaderKeyContentType]; ok {
		defaultValue = v
	}
	header.Set(HeaderKeyContentType, defaultValue)
}

func (c *restClient) pingPong(
	ctx context.Context, msg *message,
	method, path string, headers map[string]string, body interface{},
) (*http.Response, error) {

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

	var req *http.Request
	if reader, ok := body.(io.ReadCloser); ok {
		req, err = http.NewRequest(method, fullURL.String(), reader)
		defer reader.Close()

		setDefaultContentType(
			&req.Header, headers, HeaderValueContentTypeBinaryOctetStream)
	} else if body != nil {
		bodyBuffer := &bytes.Buffer{}
		enc := json.NewEncoder(bodyBuffer)
		if err = enc.Encode(body); err != nil {
			return nil, errors.Wrapf(err, "encode request body failed: %s", msg)
		}
		req, err = http.NewRequest(method, fullURL.String(), bodyBuffer)
		setDefaultContentType(
			&req.Header, headers, HeaderValueContentTypeJSON)
	} else {
		req, err = http.NewRequest(method, fullURL.String(), nil)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "new request failed: %s", msg)
	}

	isContentTypeSet := req.Header.Get(HeaderKeyContentType) != ""

	for header, value := range headers {
		if header == HeaderKeyContentType && isContentTypeSet {
			continue
		}
		req.Header.Add(header, value)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("X-EMC-REST-CLIENT", "true")
	return c.doWithRetryOnce(ctx, req, msg)
}

func (c *restClient) doWithRetryOnce(
	ctx context.Context, req *http.Request, msg *message,
) (*http.Response, error) {

	var err error
	for count := 0; count < 2; count++ {
		req.Header.Set("EMC-CSRF-TOKEN", c.csrfToken)

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
		c.csrfToken = resp.Header.Get("EMC-CSRF-TOKEN")
		return resp, nil
	}
	return nil, errors.Wrapf(err, "http request failed: %s", msg)
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

	rawResp, err := c.pingPong(ctx, msg, method, path, headers, body)
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
