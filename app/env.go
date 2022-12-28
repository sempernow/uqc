// Package app provides runtime configuration of the CLI from OS environment variables.
package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sempernow/kit/convert"
	"github.com/sempernow/uqc/client"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

// NS sets an app-wide namespace prefix for exploitation by `conf` pkg.
// Any environment var of form `<NS>_*` overwrites the configured default
// value of its field-name match at the `conf`-associated struct.
// (Field defaults are set per struct tag: `conf:"default:<VAL>"`.)
// The resets occur per `conf.Parse(..)`,
// typically at each service's initialization stage. See `main.run()`.
const NS = "APP"

// TplExt is app-wide
const TplExt = ".gohtml"

// Errors @ app layer
const (
	NotImplemented = "Not implemented"
)

// Build-time parameters : values set thereof by Golang compiler/linker per ldflags.
//
//	See Makefile.settings for sources of values
//	See Dockerfile : RUN go build -ldflags="-X '${MODULE}/app.Maker=${VENDOR}' -X ...
//	See docker image inspect : .[].Config.Labels
var (
	// Maker is the copyright holder : VENDOR
	Maker string = "@src"
	// SVN is the git HEAD hash : SVN
	SVN string = "@src"
	// Version (semantic-versioning) : VER_<SVC>
	Version string = "0.0.0"
	// Built is the build timestamp : BUILT
	Built string = "2001-01-01T01:01:01Z" //... RFC3339
)

// Service base
const (
	BASE_AOA = "/aoa/v1"
	BASE_API = "/api/v1"
)

// NewEnv returns the environment-configured receiver for client functions.
func NewEnv(osArgs []string) (*client.Env, error) {

	var cfg struct {
		conf.Version
		Args          conf.Args
		Assets        string `conf:"default:assets"`
		Cache         string `conf:"default:cache"`
		SitesPass     string `conf:"default:aPass,noprint"`
		SitesListCSV  string `conf:"default:host_channels.csv"`
		SitesListJSON string `conf:"default:_sites.json"`

		Client struct { // APP_CLIENT_*
			User  string `conf:"default:aUser"`
			Pass  string `conf:"default:aPass,noprint"`
			Token string `conf:"default:-"`
			Key   string `conf:"default:-"`

			UserAgent  string        `conf:"default:uqc/dev"`
			Timeout    time.Duration `conf:"default:5s"`
			TraceLevel int           `conf:"default:1"`
			TraceDump  bool          `conf:"default:false"`
			TraceFpath string        `conf:"default:./client.trace-resp.dump"`
		}
		Service struct {
			BaseURL string `conf:"default:http://localhost:3000"`
			// BaseAOA string `conf:"default:http://localhost:3333/aoa/v1"`
			// BaseAPI string `conf:"default:http://localhost:3000/api/v1"`
			// BasePWA string `conf:"default:http://localhost:3030"`
		}
		Channel struct {
			HostURL  string `conf:"default:https://blog.site"`
			OwnerID  string `conf:"default:00000000-0000-0000-0000-000000000000"`
			Slug     string `conf:"default:aSlug"`
			ID       string `conf:"default:00000000-0000-0000-0000-000000000000"`
			ThreadID string `conf:"default:00000000-0000-0000-0000-000000000000"`
		}
	}
	year, _, _ := time.Now().Date()
	cfg.Version.Desc = "© " + convert.IntToString(year) + " " + Maker
	cfg.Version.SVN = SVN

	// Parse the `cfg` literal. This resets any field therein from its default value
	// to that of its matching namespaced (NS) environment variable.
	// Pattern-match example: cfg.Assets.PathRoot <= APP_ASSETS_PATH_ROOT .
	// CLI Override : --assets-path-root=/some/other/path
	if err := conf.Parse(osArgs[1:], NS, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(NS, &cfg)
			if err != nil {
				return &client.Env{}, errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return &client.Env{}, nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(NS, &cfg)
			if err != nil {
				return &client.Env{}, errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return &client.Env{}, nil
		}
		return &client.Env{}, errors.Wrap(err, "parsing config")
	}

	return &client.Env{
		Logger:        log.New(os.Stdout, NS+" ", log.LstdFlags),
		Args:          cfg.Args,
		NS:            NS,
		Assets:        cfg.Assets,
		Cache:         cfg.Cache,
		SitesPass:     cfg.SitesPass,
		SitesListCSV:  cfg.SitesListCSV,
		SitesListJSON: cfg.SitesListJSON,

		Build: client.Build{
			Desc:    cfg.Desc,
			Maker:   Maker,
			SVN:     fmt.Sprintf("%.12s", SVN),
			Version: Version,
			Built:   Built,
			Year:    Built[:4],
		},

		Client: client.Client{
			User:  cfg.Client.User,
			Pass:  cfg.Client.Pass,
			Token: cfg.Client.Token,
			Key:   cfg.Client.Key,

			UserAgent:  cfg.Client.UserAgent,
			Timeout:    cfg.Client.Timeout,
			TraceLevel: cfg.Client.TraceLevel,
			TraceDump:  cfg.Client.TraceDump,
			TraceFpath: cfg.Client.TraceFpath,
		},

		Service: client.Service{
			BaseURL: cfg.Service.BaseURL,
			BaseAOA: cfg.Service.BaseURL + BASE_AOA,
			BaseAPI: cfg.Service.BaseURL + BASE_API,
			BasePWA: cfg.Service.BaseURL,
			// BaseAOA: cfg.Service.BaseAOA,
			// BaseAPI: cfg.Service.BaseAPI,
			// BasePWA: cfg.Service.BasePWA,
		},

		Channel: client.Channel{
			HostURL: cfg.Channel.HostURL,
			//OwnerID:  cfg.Channel.OwnerID,
			Slug: cfg.Channel.Slug,
			ID:   cfg.Channel.ID,
			//ThreadID: cfg.Channel.ThreadID,
		},
	}, nil
}
