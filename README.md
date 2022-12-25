# [`sempernow/uqc`](https://github.com/sempernow/uqc "GitHub")

An HTTP client library and CLI fitting the API of [`Uqrate`](https://uqrate.org "uqrate.org") services. 

```bash
$ go get -u github.com/sempernow/uqc
```

## `client` package

Each of its functions return a `client.Response`.

```golang
type Response struct {
	Body  string `json:"body,omitempty"`
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}
```

## `cli` package

The buildable CLI. Its commands execute functions of `client` or `wordpress` packages. 
See `Makefile` recipes for the commands configured for `go run ...` execution (sans build).

```bash
$ go build ./app/cli
```

### Menu of project recipes

```bash
$ make
```