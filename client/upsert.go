package client

import (
	"fmt"

	"github.com/imroc/req/v3"
)

const upsertINSTRUCT = `
	upsert <token> <mid> <&Message> [<slug>]

	- Default slug: env.Slug
	- Message.Body must not be empty.
`

// Message must fit message.Message
type Message struct {
	Title   string `json:"title,omitempty"`
	Summary string `json:"summary,omitempty"`
	//Banner   string   `json:"banner,omitempty"`
	Body string `json:"body,omitempty"`
	//Keywords []string `json:"keywords,omitempty"`
	//Ctime    string   `json:"ctime,omitempty"`
}

// UpsertStatus must fit message.UpsertStatus
type UpsertStatus struct {
	//IDx  int    `db:"idx" json:"idx"` // Required @ SELECT, e.g., @ Retrieve()
	ID   string `db:"msg_id" json:"msg_id"`
	Mode int    `db:"mode" json:"mode"` // 201, 204, 404

	// Server-sent (@ ResponseError)
	Error string `db:"-" json:"error,omitempty"`
}

// Upsert a long-form, externally-hosted message.
func (env *Env) UpsertMessage(token, mid string, msg *Message, slug ...string) *Response {

	if msg.Body == "" {
		fmt.Printf("%s\n", upsertINSTRUCT)
		return &Response{Error: ErrHelp.Error()}
	}
	chn := slug[0]
	if chn == "" {
		chn = env.Slug
	}
	endpt := env.BaseAPI + "/m/host/" + chn + "/" + mid

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
