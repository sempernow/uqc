// Package htp/main provides a CLI interface for the app's http client.
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
	upsert  :     Upsert a hosted long-form message.  
	              	upsert <token> <slug> <mid> <title> <summary> <body>
	newest  :     Retrieve list of newest messages.
	example :     Run the client's parent-repo example.

	The client always returns a client.Response struct. (WIP)

	Associated environment variables : app.NewEnv(..) and Makefile.settings .
	Command override declared value of APP_FOO_BAR with --foo-bar=newVALUE .
`

func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != client.ErrHelp {
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
		// if err := ioutil.WriteFile(env.Client.TraceFname, []byte(resp.Body), 0644); err != nil {
		// 	return errors.Wrap(err, "trace")
		// }

	case "token":
		resp := env.Token()
		fmt.Print(types.Stringify(&resp))

	case "upsert":
		token := env.Args.Num(1)
		slug := env.Args.Num(2)
		mid := env.Args.Num(3)
		resp := env.UpsertMessage(token, mid, &client.Message{
			Title:   env.Args.Num(4),
			Summary: env.Args.Num(5),
			Body:    env.Args.Num(6),
		}, slug)
		fmt.Print(types.Stringify(&resp))

	case "newest":
		//lt := env.Args.Num(1)
		if err := client.Newest(env); err != nil {
			return errors.Wrap(err, "newest")
		}
	case "example": // Repo example
		if err := client.Example(env); err != nil {
			return errors.Wrap(err, "example")
		}
	default: // Info / Describe

		fmt.Printf("\nCommands:\n")
		fmt.Printf("%s\n", DESCRIBE)

		return client.ErrHelp
	}

	return nil
}
