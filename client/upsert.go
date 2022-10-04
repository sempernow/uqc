package client

import (
	"time"

	"github.com/imroc/req/v3"
)

const ENDPT_UPSERT_KEY = "/key/m/upsert"
const ENDPT_UPSERT_TKN = "/m/upsert"

// UpsertStatus must fit message.UpsertStatus
type UpsertStatus struct {
	//IDx  int    `db:"idx" json:"idx"` // Required @ SELECT, e.g., @ Retrieve()
	ID    string `db:"msg_id" json:"msg_id"`
	Mode  int    `db:"mode" json:"mode"` // 201, 204, 404
	Error string `db:"-" json:"error,omitempty"`
}

// Message contains the payload to be decoded
// into a message.UpsertMessage by Uqrate API.
type Message struct {
	ID      string   `json:"msg_id,omitempty"` //... Key NOT EXIST @ Uqrate struct
	ChnID   string   `json:"chn_id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Summary string   `json:"summary,omitempty"`
	Body    string   `json:"body,omitempty"`
	Cats    []string `json:"cats,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	URI     string   `json:"uri,omitempty"`

	// DateCreate time.Time `db:"date_create" json:"date_create,omitempty"`
	//... Not exist @ uqrate mirror
	DateUpdate time.Time `json:"date_update,omitempty"`

	//  post_date	           post_date_gmt          post_modified          post_modified_gmt
	//  2022-09-29 19:34:40    2022-09-29 19:34:40    2022-09-29 19:34:40    2022-09-29 19:34:40 @ New
	//  2022-09-26 19:28:11    2022-09-26 19:28:11    2022-09-27 19:23:25    2022-09-27 19:23:25 @ Modified
}

// UpsertMsgByKey; a long-form Message  of an externally-hosted Channel.Slug,
// per key-authenticated ($1) POST request (X-API-KEY) to its uqrate API service endpoint.
// 	Defaults: key: env.Client.Key.
//func (env *Env) UpsertMsgByKey(msg *Message, mid, key string) *Response {
func (env *Env) UpsertMsgByKey(msg *Message, key string) *Response {
	var (
		ups = UpsertStatus{}
		rtn = Response{}
	)
	if key == "" {
		key = env.Client.Key
	}
	if key == "" {
		rtn.Error = "missing key"
	}

	if msg.ID == "" {
		rtn.Error = "missing message id"
	}
	if msg.Title == "" {
		rtn.Error = "missing message title"
	}
	if msg.Body == "" {
		rtn.Error = "missing message body"
	}

	if rtn.Error != "" {
		return &rtn
	}

	endpt := env.BaseAPI + ENDPT_UPSERT_KEY + "/" + msg.ID
	msg.ID = ""

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetHeader("x-api-key", key).
		SetResult(&ups).
		SetError(&ups).
		SetBody(&msg).
		Post(endpt)

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

// UpsertMsgByTkn; a long-form Message of an externally-hosted Channel.Slug ($2),
// per token-authenticated ($1) POST request to its uqrate API service endpoint.
// 	Defaults: slug: env.Channel.Slug, token: env.Client.Token.
func (env *Env) UpsertMsgByTkn(msg *Message, args ...string) *Response {
	var (
		jwt  = env.Client.Token
		slug = env.Channel.Slug

		ups = UpsertStatus{}
		rtn = Response{}
	)
	if len(args) > 0 {
		if args[0] != "" {
			jwt = args[0]
		}
	}
	if len(args) > 1 {
		if args[1] != "" {
			slug = args[1]
		}
	}
	if jwt == "" {
		rtn.Error = "missing token"
	}
	if slug == "" {
		rtn.Error = "missing channel slug"
	}

	if msg.ID == "" {
		rtn.Error = "missing message id"
	}
	if msg.Title == "" {
		rtn.Error = "missing message title"
	}
	if msg.Body == "" {
		rtn.Error = "missing message body"
	}

	if rtn.Error != "" {
		return &rtn
	}

	endpt := env.BaseAPI + ENDPT_UPSERT_TKN + "/" + slug + "/" + msg.ID
	msg.ID = ""

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetBearerAuthToken(jwt).
		SetResult(&ups).
		SetError(&ups).
		SetBody(&msg).
		Post(endpt)

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
