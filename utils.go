package gounity

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

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

func buildCollectionQueryURL(res string, fields string) string {
	return fmt.Sprintf("api/types/%s/instances?fields=%s", res, fields)
}

func buildInstanceQueryURL(res, id string, fields string) string {
	return fmt.Sprintf("api/instances/%s/%s?fields=%s", res, id, fields)
}
