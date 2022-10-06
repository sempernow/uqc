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

## `get`

```bash
$ cli get https://jsonplaceholder.typicode.com/todos/1
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

Upsert message by token or key

```bash
$ cli up{tkn,key} "$json" "${tkn,key}"
```

`$json`

```json
    {
        "msg_id":"1b6a7bdb-50c1-5fff-9cba-9279ca073fa5",
        "chn_id":"5cb6d760-37a2-47e0-8d7a-c86af9ed222f",
        "body": "Testing uqc upsert.",
        "title": "uqc Test",
        "summary": "Success!",
        "categories": ["foo", "foo-bar"],
        "keywords": ["foo bar"],
        "uri": "/foo/bar"
    }
```

## &nbsp;