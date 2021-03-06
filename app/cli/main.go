// Package app/cli/main provides a CLI interface for uqrate's http client.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sempernow/uqc/app"
	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/kit/types"

	"github.com/pkg/errors"
)

const DESCRIBE = `
	env     :     PrettyPrint the environment (Env) struct.
	trace   :     Trace/Debug an endpoint to os.Stderr; return response body.
	              	trace <url> [json|html(default)]
	token   :     Get access token (JWT) per Basic Auth.
	              	token |jq -Mr .body
	upsert  :     Upsert a hosted long-form message.
	              	upsert <token> <slug> <mid> <title> <summary> <body>

	Associated environment variables : app.NewEnv(..) and Makefile.settings .
	Command override any APP_* value : APP_FOO_BAR with --foo-bar=newVALUE .
`

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != ErrHelp {
			log.Printf("error: %s", err)
		}
		os.Exit(1)
	}
}

func run() error {

	env, err := app.NewEnv(os.Args)

	if err != nil {
		return errors.Wrap(err, "env")
	}

	switch env.Args.Num(0) {

	case "env":
		if err := env.PrettyPrint(); err != nil {
			return err
		}

	case "trace":
		endpt := env.Args.Num(1)
		format := env.Args.Num(2)
		resp := env.Trace(endpt, format)
		fmt.Printf("%s", resp.Body)

	case "token":
		resp := env.Token()
		fmt.Printf("%s", types.Stringify(resp))

	case "upsert":
		token := env.Args.Num(1)
		slug := env.Args.Num(2)
		mid := env.Args.Num(3)
		resp := env.UpsertMessage(token, mid, &client.Message{
			Title:   env.Args.Num(4),
			Summary: env.Args.Num(5),
			Body:    env.Args.Num(6),
		}, slug)
		fmt.Printf("%s", types.Stringify(resp))

	default: // Info / Describe

		fmt.Printf("\nCommands:\n")
		fmt.Printf("%s\n", DESCRIBE)

		return ErrHelp
	}

	return nil
}
