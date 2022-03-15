package client

import (
	"github.com/imroc/req/v3"
)

const UPSERT_ENDPT = "/m/host"

// UpsertStatus must fit message.UpsertStatus
type UpsertStatus struct {
	//IDx  int    `db:"idx" json:"idx"` // Required @ SELECT, e.g., @ Retrieve()
	ID    string `db:"msg_id" json:"msg_id"`
	Mode  int    `db:"mode" json:"mode"` // 201, 204, 404
	Error string `db:"-" json:"error,omitempty"`
}

// Upsert a long-form, externally-hosted Message to Channel.Slug, defaulting to env.Slug.
func (env *Env) UpsertMessage(token, mid string, msg *Message, slug ...string) *Response {

	if msg.Body == "" {
		return &Response{Error: "message body missing"}
	}
	chn := slug[0]
	if chn == "" {
		chn = env.Slug
	}
	if chn == "" {
		return &Response{Error: "channel slug missing"}
	}

	endpt := env.BaseAPI + UPSERT_ENDPT + "/" + chn + "/" + mid

	result := &UpsertStatus{}

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	resp, err := client.R().
		SetBearerAuthToken(token).
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
