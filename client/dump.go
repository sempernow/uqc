package client

import (
	"strings"

	"github.com/imroc/req/v3"
)

// Dump returns the raw HTML|JSON(default) @ Response.Body .
func (env *Env) Dump(url, cType string) *Response {

	var rtn Response

	if url == "" {
		rtn.Error = "missing url"
		return &rtn
	}

	if strings.ToLower(cType) == "html" {
		cType = HTML
	} else {
		cType = JSON
	}

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetHeader("Accept", cType).
		SetError(&rtn).
		Get(url)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = rsp.Status
		return &rtn
	}
	rtn.Body = rsp.String()

	return &rtn
}
