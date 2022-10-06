package client

import (
	"net/http"
	"time"

	"github.com/imroc/req/v3"
	"github.com/sempernow/uqc/kit/convert"
)

const KEY_ENDPT = "/c/key/"

type ApiKey struct {
	XID        string    `db:"xid" json:"xid,omitempty"`
	Scope      int       `db:"scope" json:"scope,omitempty"`
	Name       string    `db:"key_name" json:"key_name,omitempty"`
	Hash       []byte    `db:"key_hash" json:"-"`
	DateCreate time.Time `db:"date_create" json:"date_create"`
	DateUpdate time.Time `db:"date_update" json:"date_update,omitempty"`
	Rotations  int       `db:"rotations" json:"rotations,omitempty"`

	Key   string `db:"-" json:"key,omitempty"`
	Value string `db:"-" json:"key_value,omitempty"`
	Error string `db:"-" json:"error,omitempty"`

	// All possible fields of all scopes (models), else handle per model
	OwnerID string `db:"owner_id" json:"owner_id,omitempty"` // channels.owner_id

	ChnID   string `db:"chn_id" json:"chn_id,omitempty"`     // channels.chn_id
	ChnSlug string `db:"chn_slug" json:"chn_slug,omitempty"` // channels.slug
	HostURL string `db:"host_url" json:"host_url,omitempty"` // channels.host_url

}

// PatchKey makes token-authenticated PATCH request for ApiKey
func (env *Env) PatchKey(cid string, arg ...string) *Response {
	var (
		endpt = env.BaseAPI + KEY_ENDPT
		rtn   = Response{}
		got   = ApiKey{}
		jwt   = convert.BytesToString(env.GetCache("/keys/tkn." + env.Client.User))
	)
	if len(arg) > 0 {
		jwt = arg[0]
	}
	csrf := CSRF{CSRF: "abc123"}

	url := endpt + cid
	//GhostPrint("url: %s\n", url)

	client := req.C().
		SetUserAgent(env.UserAgent).
		SetTimeout(env.Timeout)

	rsp, err := client.R().
		SetBearerAuthToken(jwt).
		SetCookies(&http.Cookie{
			Name:  "_c",
			Value: csrf.CSRF,
		}).
		SetBody(&csrf).
		SetResult(&got).
		SetError(&got).
		Patch(url)

	if err != nil {
		rtn.Error = err.Error()
		return &rtn
	}
	rtn.Code = rsp.StatusCode

	if rsp.IsError() {
		rtn.Error = got.Error
		return &rtn
	}
	if rsp.IsSuccess() {
		rtn.Body = rsp.String()
	}
	return &rtn
}
