package gounity_test

// func TestPingPong(t *testing.T) {
// 	ctx, err := testutil.NewTestContext()
// 	assert.Nil(t, err, "failed to setup rest client to mock server")
// 	defer ctx.TearDown()

// 	resp, err := ctx.restClient.pingPong(
// 		ctx.context, http.MethodGet,
// 		fmt.Sprintf("api/instances/lun/sv_1?compact=true&fields=%s", typeFieldsLun),
// 		nil, nil)

// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)
// }

// func TestDoWithHeaders(t *testing.T) {
// 	ctx, err := testutil.NewTestContext()
// 	assert.Nil(t, err, "failed to setup rest client to mock server")
// 	defer ctx.TearDown()

// 	instResp := &instanceResp{}
// 	err = ctx.restClient.DoWithHeaders(
// 		ctx.context, http.MethodGet,
// 		fmt.Sprintf("api/instances/lun/sv_1?compact=true&fields=%s", typeFieldsLun),
// 		nil, nil, instResp)
// 	assert.Nil(t, err)

// 	lun := &Lun{}
// 	err = json.Unmarshal(instResp.Content, lun)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "sv_1", lun.Id)
// 	assert.Equal(t, "pool_1", lun.Pool.Id)
// }

// func TestGet(t *testing.T) {
// 	ctx, err := testutil.NewTestContext()
// 	assert.Nil(t, err, "failed to setup rest client to mock server")
// 	defer ctx.TearDown()

// 	instResp := &instanceResp{}
// 	err = ctx.restClient.Get(
// 		ctx.context,
// 		fmt.Sprintf("api/instances/lun/sv_1?compact=true&fields=%s", typeFieldsLun),
// 		nil, instResp)
// 	assert.Nil(t, err)

// 	lun := &Lun{}
// 	err = json.Unmarshal(instResp.Content, lun)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "sv_1", lun.Id)
// 	assert.Equal(t, "pool_1", lun.Pool.Id)
// }
