package client

import (
	"github.com/imroc/req/v3"
)

/******************************************************************************
USAGE:
		tkn, err := client.Token(env)
		if err != nil {
			return err
		}
		fmt.Println(kit.Stringify(tkn)) // JSON
*****************************************************************************/

const tokenINSTRUCT = `
	--client-user=foo
	--client-pass=bar
	--service-root=https://swarm.foo
	token
`

// Token retrieves an access token (JWT) per Basic Auth request.
func (env *Env) Token() *Response {
	var (
		user   = env.Client.User
		pass   = env.Client.Pass
		endpt  = env.BaseAOA + "/a/token"
		result = &JWT{}
	)
	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)
	resp, err := client.R().
		SetBasicAuth(user, pass).
		SetResult(&result).
		SetError(&result).
		Get(endpt)

	if err != nil {
		return &Response{
			Error: err.Error(),
		}
	}
	if resp.IsError() {
		return &Response{
			Code:  resp.StatusCode,
			Error: result.Error,
		}
	}
	return &Response{
		Code: resp.StatusCode,
		Body: result.Token,
	}
}
