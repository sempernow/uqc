package client

import (
	"github.com/imroc/req/v3"
	"github.com/sempernow/uqc/kit/convert"
)

// PutByKey makes PUT request with header: `X-API-KEY: <KEY>`
func (env *Env) PutByKey(key, url string, data interface{}) *Response {
	var (
		ups = UpsertStatus{}
		rtn = Response{}
	)
	if key == "" {
		key = convert.BytesToString(env.GetCache("/keys/key." + env.Channel.ID))
		if key == "" {
			key = env.Client.Key
		}
	}

	if key == "" {
		rtn.Error = "missing key"
	}
	if rtn.Error != "" {
		return &rtn
	}

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetHeader("x-api-key", key).
		SetResult(&ups).
		SetError(&ups).
		SetBody(&data).
		Put(url)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = ups.Error
		return &rtn
	}
	rtn.Body = ups.ID
	return &rtn
}

// PutByTkn makes a PUT request with header: `Authorization: Bearer <TKN>`
func (env *Env) PutByTkn(tkn, url string, data interface{}) *Response {
	var (
		ups = UpsertStatus{}
		rtn = Response{}
	)
	if tkn == "" {
		//tkn = env.Client.Token
		tkn = convert.BytesToString(env.GetCache("/keys/tkn." + env.Client.User))
	}
	if tkn == "" {
		rtn.Error = "missing token"
	}

	if rtn.Error != "" {
		return &rtn
	}

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetBearerAuthToken(tkn).
		SetResult(&ups).
		SetError(&ups).
		SetBody(&data).
		Put(url)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = ups.Error
		return &rtn
	}
	rtn.Body = ups.ID
	return &rtn
}
