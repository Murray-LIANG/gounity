package gounity

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
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

// RestClient is a client to send REST request.
type RestClient interface {
	// Do sends a REST request.
	Do(
		ctx context.Context,
		method, path string,
		body, resp interface{},
	) error

	// DoWithHeaders sends a Rest request with headers.
	DoWithHeaders(
		ctx context.Context,
		method, path string,
		headers map[string]string,
		body, resp interface{},
	) error

	//PingPong sends a Rest request and returns the raw response body.
	PingPong(
		ctx context.Context,
		method, path string,
		headers map[string]string,
		body interface{},
	) (*http.Response, error)

	// Get sends a request using GET method.
	Get(
		ctx context.Context,
		path string,
		headers map[string]string,
		resp interface{},
	) error

	// Post sends a request using POST method.
	Post(
		ctx context.Context,
		path string,
		headers map[string]string,
		body, resp interface{},
	) error

	// Delete sends a request using DELETE method.
	Delete(
		ctx context.Context,
		path string,
		headers map[string]string,
		body, resp interface{},
	) error
}

type client struct {
	http      *http.Client
	host      string
	username  string
	password  string
	authToken string
	csrfToken string
	traceHTTP bool
}

// RestClientOptions are options for the REST client.
type RestClientOptions struct {
	// Insecure indicates whether or not to suppress SSL errors.
	Insecure  bool
	TraceHTTP bool
}

// NewRestClient returns a new REST client to Unity.
func NewRestClient(ctx context.Context, host, username, password string,
	opts RestClientOptions) (RestClient, error) {

	if host == "" {
		return nil, newGounityError("missing host")
	}

	cookieJar, err := cookiejar.New(
		&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	c := &client{http: &http.Client{Jar: cookieJar}, host: host,
		username: username, password: password, traceHTTP: opts.TraceHTTP}

	if opts.Insecure {
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return c, nil
}

func setDefaultContentType(header *http.Header, headers map[string]string,
	defaultValue string) {
	if v, ok := headers[HeaderKeyContentType]; ok {
		defaultValue = v
	}
	header.Set(HeaderKeyContentType, defaultValue)
}

func (c *client) PingPong(ctx context.Context, method, path string,
	headers map[string]string, body interface{}) (*http.Response, error) {

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
		return nil, err
	}

	var req *http.Request
	if reader, ok := body.(io.ReadCloser); ok {
		req, err = http.NewRequest(method, fullURL.String(), reader)
		defer reader.Close()

		setDefaultContentType(&req.Header, headers,
			HeaderValueContentTypeBinaryOctetStream)
	} else if body != nil {
		bodyBuffer := &bytes.Buffer{}
		enc := json.NewEncoder(bodyBuffer)
		if err = enc.Encode(body); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, fullURL.String(), bodyBuffer)
		setDefaultContentType(&req.Header, headers,
			HeaderValueContentTypeJSON)
	} else {
		req, err = http.NewRequest(method, fullURL.String(), nil)
	}

	if err != nil {
		return nil, err
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
	return c.doWithRetryOnce(ctx, req)
}

func (c *client) doWithRetryOnce(ctx context.Context, req *http.Request) (*http.Response, error) {
	var err error
	for count := 0; count < 2; count++ {
		req.Header.Set("EMC-CSRF-TOKEN", c.csrfToken)

		if c.traceHTTP {
			traceRequest(ctx, req)
		}
		req = req.WithContext(ctx)
		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}
		if c.traceHTTP {
			traceResponse(ctx, resp)
		}
		c.csrfToken = resp.Header.Get("EMC-CSRF-TOKEN")
		return resp, nil
	}
	return nil, err
}

func (c *client) DoWithHeaders(ctx context.Context, method, path string,
	headers map[string]string, body, resp interface{}) error {
	rawResp, err := c.PingPong(ctx, method, path, headers, body)
	if err != nil {
		return err
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
			log.WithError(err).Error(
				fmt.Sprintf("unable to decode response into %+v", resp))
			return err
		}
	default:
		unityError, err := parseUnityError(rawResp.Body)
		if err != nil {
			log.WithError(err).Error(
				fmt.Sprintf("unable to decode response into unity error"))
			return err
		}
		return unityError
	}
	return nil
}

func (c *client) Do(ctx context.Context, method, path string, body,
	resp interface{}) error {
	return c.DoWithHeaders(ctx, method, path, nil, body, resp)
}

func (c *client) Get(ctx context.Context, path string, headers map[string]string,
	resp interface{}) error {

	return c.DoWithHeaders(ctx, http.MethodGet, path, headers, nil, resp)
}

func (c *client) Post(ctx context.Context, path string, headers map[string]string,
	body, resp interface{}) error {

	return c.DoWithHeaders(ctx, http.MethodPost, path, headers, body, resp)
}

func (c *client) Delete(ctx context.Context, path string, headers map[string]string,
	body, resp interface{}) error {

	return c.DoWithHeaders(ctx, http.MethodDelete, path, headers, nil, resp)
}
