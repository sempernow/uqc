// Package client provides an http client of uqrate services.
package client

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

// Levels
// 	REQUEST is sans client-added latencies.
// 	CLIENT includes client-added latencies.
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

// ghostPrint(..) prints the so-formatted args to os.Stderr.
// This is useful in any client function that both prints and returns Response.Body;
// it preserves a clean os.Stdout to which caller may then dump the return (JSON/HTML).
func ghostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
