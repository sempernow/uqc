// Package client provides an http client of uqrate services.
package client

import (
	"time"

	"github.com/ardanlabs/conf"
)

const (
	CacheKeyTknPrefix = "tkn."
)

// Levels
//
//	REQUEST is sans client-added latencies.
//	CLIENT includes client-added latencies.
const (
	REQUEST = iota + 1
	CLIENT
)

// Mode
const (
	MODE_JSON   = 1
	MODE_STRUCT = 2
)

// Service base
const (
	BASE_AOA = "/aoa/v1"
	BASE_API = "/api/v1"
)

// Content type for HTTP request header : "Accept: ..."
const (
	JSON = "application/json"
	HTML = "text/html"
)

type CSRF struct {
	CSRF string `json:"csrf"`
}

// Response is the return of all (exported) client function calls.
type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
	// Error `json:"error,omitempty"`
}

// TODO : Mod to this; See uqrate error response;
// has this detailed morphadite error response which gets lost in decode.
type Error struct {
	Error  string   `json:"error,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

// Env is the receiver of all (exported) client functions,
// and contains all parameters defining the client environment.
type Env struct {
	Args          conf.Args `json:"args,omitempty"`
	NS            string    `json:"ns,omitempty"`
	Build         `json:"build"`
	Assets        string `json:"assets,omitempty"`
	Cache         string `json:"cache,omitempty"`
	SitesPass     string `json:"sites_pass,omitempty"`
	SitesListCSV  string `json:"sites_list_csv,omitempty"`
	SitesListJSON string `json:"sites_list_json,omitempty"`
	Client        `json:"client,omitempty"`
	Service       `json:"service,omitempty"`
	Channel       `json:"channel,omitempty"`
}

// Build contains application build info.
type Build struct {
	Desc    string `json:"desc,omitempty"`
	Maker   string `json:"maker,omitempty"`
	SVN     string `json:"svn,omitempty"`
	Version string `json:"version,omitempty"`
	Built   string `json:"built,omitempty"`
	Year    string `json:"year,omitempty"`

	// BuiltAOA string
	// BuiltAPI string
	// BuiltPWA string
}

// Client contains all request parameters.
type Client struct {
	User       string        `json:"user,omitempty"`
	Pass       string        `json:"pass,omitempty"`
	Token      string        `json:"token,omitempty"`
	Key        string        `json:"key,omitempty"`
	UserAgent  string        `json:"user_agent,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
	TraceLevel int           `json:"trace_level,omitempty"`
	TraceDump  bool          `json:"trace_dump,omitempty"`
	TraceFpath string        `json:"trace_fpath,omitempty"`
}

// Service regards that requested by Client; that servicing Message(s) of Channel(s).
type Service struct {
	// <hostname> (domain)
	Host string `json:"host,omitempty"`
	// <scheme>://<hostname>:<port>
	BaseURL string `json:"base_url,omitempty"`
	// <scheme>://<hostname>:<port>/<service>/<version>
	BaseAOA string `json:"base_aoa,omitempty"`
	BaseAPI string `json:"base_api,omitempty"`
	BasePWA string `json:"base_pwa,omitempty"`
}

const MaxUserDisplay = 30

// User regards that of the service account, channel owner, hosting site, and credited author.
type User struct {
	ID      string `json:"user_id,omitempty"` // Not contained in service user.UpdateUser
	Display string `json:"display,omitempty"`
	About   string `json:"about,omitempty"`
	Avatar  string `json:"avatar,omitempty"`
	Banner  string `json:"banner,omitempty"`
}

// Channel regards that to which a message is upserted; at store of Service host.
type Channel struct {
	ID string `json:"chn_id,omitempty"` // Channel.ID
	// host     :            <hostname>           : domain
	// host url : <scheme>://<hostname>:<port>    : url
	HostURL string `json:"host_url,omitempty"` // Channel.HostURL
	Slug    string `json:"slug,omitempty"`     // Channel.Slug
	OwnerID string `json:"owner_id,omitempty"` // Channel.OwnerID (UUID v4)

	// Thread-root message (long-form)
	ThreadID string `json:"thread_id,omitempty"` // Message.ID (UUID v5)

	Title string   `json:"title,omitempty"`
	About string   `json:"about,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

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
// into a message.UpdateMessage by Uqrate API.
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
}
