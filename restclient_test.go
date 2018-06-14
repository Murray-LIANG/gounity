package gounity

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func setup() (RestClient, context.Context) {
	f, err := os.OpenFile("/tmp/gounity.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})

	ctx := context.Background()
	cli, _ := NewRestClient(ctx, "https://10.245.101.39",
		"admin", "Password123!", RestClientOptions{Insecure: true, TraceHTTP: true})
	return cli, ctx
}

func TestPingPong(t *testing.T) {
	cli, ctx := setup()
	resp, err := cli.PingPong(ctx, http.MethodGet, "api/instances/pool/pool_1", nil,
		nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDoWithHeaders(t *testing.T) {
	cli, ctx := setup()
	lunResp := &lunResp{}
	err := cli.DoWithHeaders(ctx, http.MethodGet,
		fmt.Sprintf("api/instances/lun/sv_1?fields=%s", fieldsLUN), nil, nil, lunResp)

	assert.Nil(t, err)
	lun := lunResp.Content
	assert.Equal(t, "sv_1", lun.ID)
	assert.Equal(t, "pool_1", lun.Pool.ID)
}

func TestGet(t *testing.T) {
	cli, ctx := setup()
	lunResp := &lunResp{}
	err := cli.Get(ctx,
		fmt.Sprintf("api/instances/lun/sv_1?fields=%s", fieldsLUN), nil, lunResp)

	assert.Nil(t, err)
	lun := lunResp.Content
	assert.Equal(t, "sv_1", lun.ID)
	assert.Equal(t, "pool_1", lun.Pool.ID)
}
