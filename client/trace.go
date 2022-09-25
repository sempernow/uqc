package client

import (
	"io/ioutil"
	"os"

	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
)

// Levels
// 	REQUEST is sans client-added latencies.
// 	CLIENT includes client-added latencies.
const (
	REQUEST = iota + 1
	CLIENT
)

// Trace hits the declared endpoint (any) with a GET request,
// prints response-timing info to os.Stderr, and returns response body else info.
// Dumps body to file instead if both TraceDump flag and TraceFpath set (see Env.Client).
// https://github.com/imroc/req#Debugging
func (env *Env) Trace(endpt, cType string) *Response {

	switch cType {
	case "json":
		cType = JSON
	default:
		cType = HTML
	}

	opt := &req.DumpOptions{
		Output:         os.Stderr,
		RequestHeader:  true,
		ResponseBody:   false,
		RequestBody:    false,
		ResponseHeader: true,
		Async:          false,
	}

	var (
		result interface{} // unncessary?
		resp   *req.Response
		err    error
	)
	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout).
		SetCommonDumpOptions(opt).
		EnableDumpAll()

	switch env.TraceLevel {

	case CLIENT:
		client.EnableTraceAll()
		resp, err = client.R().
			SetHeader("Accept", cType).
			Get(endpt)
	case REQUEST:
		fallthrough
	default:
		resp, err = client.R().EnableTrace().
			SetHeader("Accept", cType).
			SetResult(&result).
			SetError(&result).
			Get(endpt)
	}

	if err != nil {
		return &Response{Error: errors.Wrap(err, "trace").Error()}
	}
	trace := resp.Request.TraceInfo()
	ghostPrint("%v\n%s\n%v\n\n", trace.Blame(), "----------", trace)

	if err != nil {
		return &Response{
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return &Response{
			Code:  resp.StatusCode,
			Error: resp.Status,
		}
	}

	// Dump successful response to file (env.TraceFpath), conditionally.
	if env.TraceDump && (env.TraceFpath != "") {
		if err := ioutil.WriteFile(env.TraceFpath, resp.Bytes(), 0644); err != nil {
			return &Response{
				Code:  resp.StatusCode,
				Body:  resp.String(),
				Error: errors.Wrap(err, "@ ioutil.WriteFile(..) : '"+env.TraceFpath+"'").Error(),
			}
		}
		return &Response{
			Code: resp.StatusCode,
			Body: "Response body dumped to: " + env.TraceFpath,
		}
	}

	return &Response{
		Code: resp.StatusCode,
		Body: resp.String(),
	}
}
