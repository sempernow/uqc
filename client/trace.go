package client

import (
	"os"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

// Trace hits the declared endpoint (any) with a GET request,
// prints response-timing info to os.Stderr, and returns response body else info.
// Dumps body to file instead if both TraceDump flag and TraceFpath set (see Env.Client).
// https://github.com/imroc/req#Debugging
func (env *Env) Trace(endpt, cType string) *Response {

	if strings.ToLower(cType) == "html" {
		cType = HTML
	} else {
		cType = JSON
	}

	opts := &req.DumpOptions{
		Output:         os.Stderr,
		RequestHeader:  true,
		ResponseBody:   false,
		RequestBody:    false,
		ResponseHeader: true,
		Async:          false,
	}

	var (
		err error
		rsp *req.Response
		rtn Response
	)
	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout).
		SetCommonDumpOptions(opts).
		EnableDumpAll() //.EnableDebugLog()

	switch env.TraceLevel {

	case CLIENT:

		rsp, err = client.EnableTraceAll().R().
			SetHeader("Accept", cType).
			Get(endpt)

	case REQUEST:
		fallthrough
	default:

		rsp, err = client.R().EnableTrace().
			SetHeader("Accept", cType).
			Get(endpt)
	}

	if err != nil {
		rtn.Error = errors.Wrap(err, "trace").Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	trace := rsp.Request.TraceInfo()
	GhostPrint("%v\n%s\n%v\n\n", trace.Blame(), "----------", trace)

	if rsp.IsError() {
		rtn.Error = rsp.Status
		return &rtn
	}
	if rsp.IsSuccess() {
		rtn.Body = rsp.String()
	}

	// Dump successful response to file (env.TraceFpath) conditionally per env setting.
	if env.TraceDump && (env.TraceFpath != "") {
		if err := os.WriteFile(env.TraceFpath, rsp.Bytes(), 0644); err != nil {
			rtn.Error = errors.Wrap(
				err, "@ WriteFile : '"+env.TraceFpath+"'",
			).Error()

			return &rtn
		}
		rtn.Body = "Response body dumped to: " + env.TraceFpath
		return &rtn
	}
	return &rtn
}
