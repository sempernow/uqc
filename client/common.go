// Package client provides an http client for accessing uqrate services.
package client

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

const UserAgent = "uqrate/client/v1 (https://uqrate.org)"

// Service base
const (
	BASE_AOA = "/aoa/v1"
	BASE_API = "/api/v1"
)

type Env struct {
	NS      string    `json:"ns,omitempty"`
	Args    conf.Args `json:"args,omitempty"`
	Client  `json:"client,omitempty"`
	Service `json:"service,omitempty"`
	Channel `json:"channel,omitempty"`
}

// Request parameters
type Client struct {
	User       string        `json:"user,omitempty"`
	Pass       string        `json:"pass,omitempty"`
	UserAgent  string        `json:"user_agent,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
	TraceLevel int           `json:"trace_level,omitempty"`
	TraceDump  bool          `json:"trace_dump,omitempty"`
	TraceFpath string        `json:"trace_fpath,omitempty"`
}

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

type Channel struct {
	// <hostname> (domain)
	Host string `json:"host,omitempty"`
	// <scheme>://<hostname>:<port>
	HostURL string `json:"host_url,omitempty"` // Channel.Host
	Slug    string `json:"slug,omitempty"`     // Channel.Slug
	OwnerID string `json:"owner_id,omitempty"` // Channel.OwnerID (UUID v4)

	// Thread-root message (long-form)
	ThreadID string `json:"thread_id,omitempty"` // Message.ID (UUID v5)
}

type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

type JWT struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// Print the so-formatted args to os.Stderr.
// So callers may cleanly dump Response.Body (JSON/HTML) to os.Stdout.
func ghostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
