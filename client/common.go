// Package client provides an http client of uqrate services.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	Args      conf.Args `json:"args,omitempty"`
	NS        string    `json:"ns,omitempty"`
	Assets    string    `json:"assets,omitempty"`
	Cache     string    `json:"cache,omitempty"`
	SitesPass string    `json:"sites_pass,omitempty"`
	Client    `json:"client,omitempty"`
	Service   `json:"service,omitempty"`
	Channel   `json:"channel,omitempty"`
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
	ID string `json:"chn_id,omitempty"` // Channel.ID
	// host     :            <hostname>           : domain
	// host url : <scheme>://<hostname>:<port>    : url
	HostURL string `json:"host_url,omitempty"` // Channel.HostURL
	Slug    string `json:"slug,omitempty"`     // Channel.Slug
	OwnerID string `json:"owner_id,omitempty"` // Channel.OwnerID (UUID v4)

	// Thread-root message (long-form)
	ThreadID string `json:"thread_id,omitempty"` // Message.ID (UUID v5)
}

// GhostPrint(..) prints args per format, e.g., "want: %s\nhave: %s\n", to os.Stderr.
func GhostPrint(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// SetCache writes data to key file in Env.Cache folder.
func (env *Env) SetCache(key, data string) error {
	if key == "" || data == "" {
		key = "ERR @ env.SetCache(..) : missing parameter(s)"
		GhostPrint("\nkey: '%s'\ndata: '%s'\n", key, data)
		return errors.New(key)
	}
	err := ioutil.WriteFile(
		filepath.Join(env.Cache, key),
		[]byte(data),
		0664,
	)
	if err != nil {
		GhostPrint("\nERR @ env.SetCache(..) : writing file: %v\n", err)
		return err
	}
	fi, err := os.Stat(filepath.Join(env.Cache, key))
	if err == nil {
		GhostPrint("\nsize: %vKB @ file: '%s'\n", (fi.Size() / 1024), key)
	}
	return nil
}

// GetCache reads key file of Env.Cache folder.
func (env *Env) GetCache(key string) []byte {
	if key == "" {
		GhostPrint("\nERR @ env.GetCache(..) : missing key parameter\n")
		GhostPrint("\nkey: '%s'\n", key)
		return []byte{}
	}
	bb, err := ioutil.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nERR @ env.GetCache(..) : reading file: %s\n", err)
		return []byte{}
	}
	//GhostPrint("\n@ %s_CACHE : %s\n", env.NS, filepath.Join(env.Cache, key))
	return bb
}

// GetCacheJSON reads key file of Env.Cache folder into struct of pointer.
func (env *Env) GetCacheJSON(key string, ptr interface{}) {
	if key == "" {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : missing key parameter\n")
	}
	bb, err := ioutil.ReadFile(filepath.Join(env.Cache, key))
	if err != nil {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : reading file: %s\n", err)
	}
	if err := json.Unmarshal(bb, &ptr); err != nil {
		GhostPrint("\nERR @ env.GetCacheJSON(..) : decoding JSON: %s\n", err)

	}
}
