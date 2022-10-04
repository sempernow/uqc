# [`sempernow/uqc`](https://github.com/sempernow/uqc "GitHub")

An http client and CLI (demo) for the [`uqrate`](https://uqrate.org "uqrate.org") project. The uqrate client wraps that of [imroc/req](https://github.com/imroc/req "GitHub") . 

```bash
$ go get -u github.com/sempernow/uqc
```

## `client`

```golang
import "github.com/sempernow/uqc/client"
```

The client package provides an http client as a golang library to access uqrate services. Its functions return a `client.Response`.

```golang
type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
```

## `cli`

This commandline app is a template for building a standalone CLI (binary). Each of its commands is a function of the `client` package, so it also serves as __a reference for utilizing the `client` package__ in other Golang packages.

```bash
$ go build ./app/cli
```

### `token`

```bash
$ cli --service-base-url=https://uqrate.org \
      --client-user=$user \
      --client-pass=$pass \
      token
```
- Override environment at commandline:  
  `--foo-bar=newValue` overrides `APP_FOO_BAR`

### `uptkn` | `upkey`

```bash
mid=$(uuid -v 5 ns:OID /2022/09/15/uqrate-client-test)
$ cli uptkn "$json" "$mid" "${tkn}" "${APP_CHANNEL_SLUG}"
$ cli upkey "$json" "$mid" "${key}
```

`$json`

```json
{
	"body": "Testing uqc upsert.",
	"title": "Uqrate Client Test",
	"keywords": ["foo", "bar"],
	"uri": "/2022/09/15/testing-uqc-upsert"
}
```

## &nbsp;