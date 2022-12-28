package client

import (
	"github.com/imroc/req/v3"
	"github.com/sempernow/kit/convert"
)

// UpsertMsgByTkn performs a POST request to Uqrate's API service endpoint
// for upserting a long-form Message (msg) of an externally-hosted
// Channel.Slug (slug) using bearer-token (token) authorization.
//
//	Defaults:
//		token (args[0]): env.GetCache(client.CacheKeyTknPrefix + env.Client.User)
//		                 @ ${APP_CACHE}/tkn.${APP_CLIENT_USER}
//		slug  (args[1]): env.Channel.Slug
//		                 @ ${APP_CHANNEL_SLUG}
func (env *Env) UpsertMsgByTkn(msg *Message, args ...string) *Response {
	var (
		jwt  string
		slug = env.Channel.Slug

		got = UpsertStatus{}
		rtn = Response{}
	)

	if len(args) > 0 {
		if args[0] != "" {
			jwt = args[0]
		}
	}
	if jwt == "" {
		jwt = convert.BytesToString(env.GetCache(CacheKeyTknPrefix + env.Client.User))
	}
	if jwt == "" {
		rtn.Error = "missing token"
	}

	if len(args) > 1 {
		if args[1] != "" {
			slug = args[1]
		}
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
	msg.ChnID = ""

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetBearerAuthToken(jwt).
		SetResult(&got).
		SetError(&got).
		SetBody(&msg).
		Post(endpt)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = got.Error
		return &rtn
	}
	rtn.Body = got.ID
	return &rtn
}

// UpsertMsgByKey performs a POST request to Uqrate's API service endpoint
// for upserting a long-form Message (msg) of an externally-hosted Channel.
// Authorization to that protected endpoint is by ApiKey (key),
// scoped to its target channel, sent as value of X-API-KEY header.
//
//	Defaults: key: env.GetCacheJSON(..)
//	               @ "${APP_CACHE}/keys/key." + msg.ChnID + ".json"
func (env *Env) UpsertMsgByKey(msg *Message, key string) *Response {
	var (
		got = UpsertStatus{}
		rtn = Response{}
	)
	if key == "" {
		k := ApiKey{}
		env.GetCacheJSON("/key."+msg.ChnID+".json", &k)
		key = k.Key
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
	msg.ChnID = ""

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetHeader("x-api-key", key).
		SetResult(&got).
		SetError(&got).
		SetBody(&msg).
		Post(endpt)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = got.Error
		return &rtn
	}
	rtn.Body = got.ID
	return &rtn
}
