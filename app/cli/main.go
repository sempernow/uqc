// Package app/cli/main provides a CLI interface for uqrate's http client.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sempernow/uqc/app"
	"github.com/sempernow/uqc/app/cli/commands"
	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/kit/timestamp"
	"github.com/sempernow/uqc/kit/types"

	"github.com/pkg/errors"
)

const DESCRIBE = `
	env     :     PrettyPrint the environment (Env) struct.
	dump    :     Dump response body of any endpoint to STDOUT and HTTP status to STDERR.
	              	dump $url ['html'|'json'(default)]
	trace   :     Trace/Debug any endpoint to STDERR and response body to STDOUT.
	              	trace $url ['html'|'json'(default)]  (to file per Makefile.settings)
	token   :     Get access token (JWT) per Basic Auth.
	              	token [$user $pass] |jq -Mr .body
	uptkn   :     Upsert a long-form message of hosted channel using JWT authentication.
	              	uptkn $json [$jwt [$slug]]
	uptkey  :     Upsert a long-form message of hosted channel using API key authentication.
	              	uptkn $json [$key]
	wpfetch :     Fetch WordPress Posts from the declared URL 
	              and dump JSON response body to file @ ./wp_posts.<DOMAIN>.json
	              	wpfetch $url
	wpupkey :     Convert and Upsert fetched posts (JSON file) of a WordPress site.

	Associated environment variables : app.NewEnv(..) and Makefile.settings .
	Command override any APP_* value : APP_FOO_BAR with --foo-bar=newVALUE .

	Run any per ` + "`make gorun`" + ` using $makeargs :

	    $ export makeargs='cli trace https://jsonplaceholder.typicode.com/todos/1'
	    $ make gorun
`

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

func main() {
	if err := run(); err != nil {
		if errors.Cause(err) != ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s", err)
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
		rsp := env.Trace(endpt, format)
		fmt.Printf("%s", rsp.Body)
		fmt.Fprintf(os.Stderr, "%s", rsp.Error)

	case "dump":
		endpt := env.Args.Num(1)
		format := env.Args.Num(2)
		rsp := env.Dump(endpt, format)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)

	case "token":
		user := env.Args.Num(1)
		pass := env.Args.Num(2)
		rsp := env.Token(user, pass)
		fmt.Printf("%s", types.Stringify(rsp))

	case "uptkn":
		// Upsert 1 JSON Message
		j := env.Args.Num(1)
		jwt := env.Args.Num(2)
		slug := env.Args.Num(3)
		msg := client.Message{}
		if err := json.Unmarshal([]byte(j), &msg); err != nil {
			return errors.Wrap(err, "decoding JSON message")
		}
		msg.Body = timestamp.TimeStringZulu(time.Now().UTC()) + " per JWT"
		rsp := env.UpsertMsgByTkn(&msg, jwt, slug)
		fmt.Printf("\n%s\n", types.Stringify(rsp))

	case "upkey":
		// Upsert 1 JSON Message
		j := env.Args.Num(1)
		key := env.Args.Num(2)
		msg := client.Message{}
		if err := json.Unmarshal([]byte(j), &msg); err != nil {
			return errors.Wrap(err, "decoding JSON message")
		}
		msg.Body = timestamp.TimeStringZulu(time.Now().UTC()) + " per ApiKey"
		rsp := env.UpsertMsgByKey(&msg, key)
		fmt.Printf("\n%s\n", types.Stringify(rsp))

	case "wpfetch":
		// Fetch per WordPress site : Any endpoint : /posts, /tags, /categories, /users
		rsp := env.Dump(env.Args.Num(1), client.JSON)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)

	// posts := []client.WordPressPost{}
	// if err := json.Unmarshal([]byte(rsp.Body), &posts); err != nil {
	// 	return errors.Wrap(err, "decoding JSON posts")
	// }
	// fmt.Printf("%s\n", types.Stringify(posts))
	case "test":
		commands.Test(env)

	case "wpuptkn":
		// Upsert /posts of WordPress site (JSON) per JWT auth.
		// Convert from the WordPress JSON to Messages struct
		path := env.Args.Num(1)
		jwt := env.Args.Num(2)
		slug := env.Args.Num(3)
		ownerslug := env.Args.Num(4)
		json, _ := ioutil.ReadFile(path)
		msgs, err := commands.Posts2Msgs(string(json), ownerslug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		//fmt.Printf("\n%s\n", types.Stringify(msgs))
		// fmt.Printf("\n%s\n", types.PrettyPrint(msgs))
		for _, msg := range msgs {
			rsp := env.UpsertMsgByTkn(&msg, jwt, slug)
			fmt.Printf("\n%s\n", types.Stringify(rsp))
		}
	case "wpupkey":
		// Upsert /posts of WordPress site (JSON) per API key auth.
		// Convert from the WordPress JSON to Messages struct
		path := env.Args.Num(1)
		key := env.Args.Num(2)
		cid := env.Args.Num(3)
		json, _ := ioutil.ReadFile(path)
		msgs, err := commands.Posts2Msgs(string(json), cid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		//fmt.Printf("\n%s\n", types.Stringify(msgs))
		// fmt.Printf("\n%s\n", types.PrettyPrint(msgs))
		for _, msg := range msgs {
			rsp := env.UpsertMsgByKey(&msg, key)
			fmt.Printf("\n%s\n", types.Stringify(rsp))
		}

	default: // Ghost print so pipe okay: ... |jq .
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "%s\n", DESCRIBE)

		return nil
	}

	return nil
}
