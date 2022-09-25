package client

import (
	"github.com/imroc/req/v3"
)

const UPSERT_ENDPT = "/m/upt"

// UpsertStatus must fit message.UpsertStatus
type UpsertStatus struct {
	//IDx  int    `db:"idx" json:"idx"` // Required @ SELECT, e.g., @ Retrieve()
	ID    string `db:"msg_id" json:"msg_id"`
	Mode  int    `db:"mode" json:"mode"` // 201, 204, 404
	Error string `db:"-" json:"error,omitempty"`
}

// Upsert a long-form, externally-hosted Message to Channel.Slug, defaulting to env.Slug.
func (env *Env) UpsertMessage(msg *Message, mid string, args ...string) *Response {

	if msg.Body == "" {
		return &Response{Error: "message body missing"}
	}
	jwt := env.Client.Token
	if len(args) > 0 {
		jwt = args[0]
	}
	slug := env.Channel.Slug
	if len(args) > 1 {
		slug = args[1]
	}
	if jwt == "" {
		return &Response{Error: "missing token"}
	}
	if slug == "" {
		return &Response{Error: "missing channel slug"}
	}

	endpt := env.BaseAPI + UPSERT_ENDPT + "/" + slug + "/" + mid
	// endpt = "http://swarm.foo:3000/api/v1/m/upt/TestHostSlug/1b6a7bdb-50c1-5fff-9cba-9279ca073fa5"

	// fmt.Printf("endpt: %s\nmid: %s\nslug: %s\njwt: %s\nmsg: %v\n",
	// 	endpt, mid, slug, jwt, types.PrettyPrint(&msg),
	// )

	result := &UpsertStatus{}

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	resp, err := client.R().
		SetBearerAuthToken(jwt).
		SetResult(&result).
		SetError(&result).
		SetBody(&msg).
		Post(endpt)

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
		Body: result.ID,
	}
}
