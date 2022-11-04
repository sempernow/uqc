package client

import (
	"github.com/imroc/req/v3"
)

/******************************************************************************
USAGE:
		rsp, err := env.Token()
		if err != nil {
			return err
		}
		fmt.Println(kit.Stringify(rsp)) // JSON
*****************************************************************************/

const TKN_ENDPT = "/a/token"

// JWT must fit response body of Token(..) request on success.
type JWT struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// Token retrieves an access token (JWT) per Basic Auth request.
//
//	Defaults: user (args[0]): Env.Client.User, pass (args[1]): Env.Client.Pass
func (env *Env) Token(args ...string) *Response {
	var (
		user  = env.Client.User
		pass  = env.Client.Pass
		endpt = env.BaseAOA + TKN_ENDPT
		got   = JWT{}
		rtn   = Response{}
	)
	if len(args) > 0 {
		if args[0] != "" {
			user = args[0]
		}
	}
	if len(args) > 1 {
		if args[1] != "" {
			pass = args[1]
		}
	}
	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetBasicAuth(user, pass).
		SetResult(&got).
		SetError(&got).
		Get(endpt)

	if err != nil {
		return &Response{
			Error: err.Error(),
		}
	}

	if rsp.IsError() {
		return &Response{
			Code:  rsp.StatusCode,
			Error: got.Error,
		}
	}

	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = rsp.Status
		return &rtn
	}
	if rsp.IsSuccess() {
		rtn.Body = rsp.String()
	}

	//GhostPrint("rtn: %s", rtn.Body)
	return &Response{
		Code: rsp.StatusCode,
		Body: got.Token,
		//Body: rtn.Body,
	}
}
