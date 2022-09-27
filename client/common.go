// Package client provides an http client for accessing uqrate services.
package client

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

const UserAgent = "uqrate/client/v1 (https://uqrate.org)"

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

// Response is the return of all (exported) client function calls.
type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

// JWT must fit response body of Token(..) request on success.
type JWT struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// Env is the receiver of all (exported) client functions,
// and contains all parameters defining the client environment.
type Env struct {
	NS      string    `json:"ns,omitempty"`
	Args    conf.Args `json:"args,omitempty"`
	Client  `json:"client,omitempty"`
	Service `json:"service,omitempty"`
	Channel `json:"channel,omitempty"`
}

// Client contains all request parameters.
type Client struct {
	User       string        `json:"user,omitempty"`
	Pass       string        `json:"pass,omitempty"`
	Token      string        `json:"token,omitempty"`
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

// Channel regards that to which a message is upserted; at store of Service host.
type Channel struct {
	// host     :            <hostname>           : domain
	// host url : <scheme>://<hostname>:<port>    : url
	HostURL string `json:"host_url,omitempty"` // Channel.HostURL
	Slug    string `json:"slug,omitempty"`     // Channel.Slug
	OwnerID string `json:"owner_id,omitempty"` // Channel.OwnerID (UUID v4)

	// Thread-root message (long-form)
	ThreadID string `json:"thread_id,omitempty"` // Message.ID (UUID v5)
}

// Message must fit message.Message of Service.
type Message struct {
	Title      string   `json:"title,omitempty"`
	Summary    string   `json:"summary,omitempty"`
	Body       string   `json:"body,omitempty"`
	Keywords   []string `json:"keywords,omitempty"`
	Categories []string `json:"categories,omitempty"`
	URI        string   `json:"uri,omitempty"`

	// DateCreate time.Time `db:"date_create" json:"date_create,omitempty"`
	//... Not exist @ uqrate mirror
	DateUpdate time.Time `db:"date_update" json:"date_update,omitempty"`
}

// ghostPrint(..) prints the so-formatted args to os.Stderr.
// This is useful in any client function that both prints and returns Response.Body;
// it preserves a clean os.Stdout to which caller may then dump the return (JSON/HTML).
func ghostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
