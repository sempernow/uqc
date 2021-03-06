# [`sempernow/uqc`](https://github.com/sempernow/uqc "GitHub")

An http client and CLI (demo) for the [`uqrate`](https://uqrate.org "uqrate.org") project. The uqrate client wraps that of [imroc/req](https://github.com/imroc/req "GitHub") . 

## Package `/client`

Package client provides an http client as a golang library to access uqrate services. Its functions return a `client.Response`.

```golang
type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
```

### Use the package

```bash
go get -u github.com/sempernow/uqc
```

```golang
import "github.com/sempernow/uqc/client"
```

See the CLI package, [`main.go`](app/cli/main.go "app/cli/main.go"), for examples of `client` calls.

## Package `/app/cli`

Each command is a function of the `client` package. So it serves as a template for building a standalone CLI (binary), and as a reference for utilizing the `client` package in other Golang packages.

## Demonstate `client` package using the CLI

### `go run ...` 

#### `trace`

```bash
go run ./app/cli trace https://gorest.co.in/public/v2/users json |jq .
```
```bash
GET /public/v2/users HTTP/1.1
...
HTTP/1.1 200 OK
...
the request total time is 1.1144131s, and costs 575.8156ms on tls handshake
----------
TotalTime         : 1.1144131s
DNSLookupTime     : 15.0618ms
TCPConnectTime    : 250.8347ms
TLSHandshakeTime  : 575.8156ms
FirstResponseTime : 266.5266ms
ResponseTime      : 6.094ms
IsConnReused:     : false
RemoteAddr        : 139.59.66.125:443

[
  {
    "id": 3780,
    "name": "Ekdant Kocchar",
    "email": "kocchar_ekdant@hoppe-lang.name",
    "gender": "male",
    "status": "active"
  },
  ...
]
...
```
- Optionally dump to `APP_CLIENT_TRACE_FPATH`

#### `env`

Print JSON of the injected `client.Env`

```bash
go run ./app/cli \
    --service-base-url=https://uqrate.org \
    --client-user=$user \
    --client-pass=$pass \
    --client-timeout=3s \
    --channel-host-url=http://127.0.0.1:5500 \
    --channel-slug=$slug \
    env
```

#### `token`

Get a token (cryptographically-signed JWT) against uqrate-member credentials.

```bash
go run ./app/cli \
    --service-base-url=https://uqrate.org \
    --client-user=$user \
    --client-pass=$pass \
    token |jq -Mr .body
@ Token:
        user: usertest
        pass: 111?????????
```
```bash
eyJ...OfA
```
- The raw token prints to stdout; All else to stderr, so can pipe.

#### `upsert`

Insert/update a long-form message (article) of a uqrate channel as a channel-hosting member and owner of the channel. Such channel hosts have access to the per-article reply messages configured as a uqrate-hosted comments section running in an iframe with each such article at their site. See [`uqrate.js`](https://uqrate.org/sa/scripts/uqrate.js) .

```bash
export token=$(go run ./app/cli --service-base-url=https://uqrate.org --client-user=$user --client-pass=$pass token |jq -Mr .body)

go run ./app/cli \
    --service-base-url=$svc_base_url \
    upsert $token $chn_slug $(uuid -v 5 ns:OID abc123) 'A NEW title' 'A summary.' 'NEW body here.'
```
```bash
# @ 201
{"body":"e8fc2054-11c8-5d32-9bb5-4b857504122e","code":201}

# @ 204
{"code":204}

# @ 401
{"code":401,"error":"token and refresh-reference cookie invalid: http: named cookie not present"}
```

### @ `make gorun` | [`Makefile.settings`](Makefile.settings)

Same as above, yet per [`Makefile` recipe](Makefile) configured per `makeargs` param.

#### `trace`

```bash
export makeargs='trace https://gorest.co.in/public/v2/users json |jq .'
make gorun
```
```bash
bash make.go.run.app.sh cli trace https://gorest.co.in/public/v2/users json |jq .
...
```

#### `token`

```bash
export makeargs='token |jq -Mr .body'
make gorun
```
```bash
bash make.go.run.app.sh cli token
eyJ...gNQ
```

## Notes on [TLS : performance impact](https://blog.yugabyte.com/measuring-the-performance-impact-of-tls-encryption-using-tpcc/ "2021 'Measuring the Performance Impact of TLS Encryption Using TPC-C'")

@ `strace` (Linux utility) 

```bash
strace -o 'out.log' -f -tt curl -H 'Accept: application/json' https://uqrate.org/app
```

## &nbsp;