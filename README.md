# [`sempernow/uqc`](https://github.com/sempernow/uqc "GitHub")

A *developmental* http client and CLI (demo) for the [`uqrate`](https://uqrate.org "uqrate.org") project. The uqrate client wraps that of [imroc/req](https://github.com/imroc/req "GitHub") . 

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

This is the prototype standalone CLI. Each of its commands is a function of the `client` or `wordpress` packages.

```bash
$ go build ./app/cli
```

## `get`

```bash
$ cli get https://jsonplaceholder.typicode.com/todos/1
```