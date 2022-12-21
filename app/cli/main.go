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

	//"github.com/sempernow/uqc/app/cli/commands"
	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/client/wordpress"
	"github.com/sempernow/uqc/kit/convert"
	"github.com/sempernow/uqc/kit/timestamp"

	"github.com/pkg/errors"
)

const DESCRIBE = `
	env         :     PrettyPrint the environment (Env) struct.
	get         :     Dump response body of GET to STDOUT and HTTP status to STDERR.
	                  	get $url ['html'|'json'(default)]
	posttkn     :     Dump response body of token-authenticated POST 
	                  	to STDOUT and HTTP status to STDERR.
	postkey     :     Dump response body of key-authenticated POST 
	                  	to STDOUT and HTTP status to STDERR.
	                  	postkey $url ['html'|'json'(default)]
	trace       :     Trace/Debug any endpoint to STDERR and response body to STDOUT.
	                  	trace $url ['html'|'json'(default)]  (to file per Makefile.settings)
	token       :     Get access token (JWT) per Basic Auth and store in cache.
	                  	token [$user $pass] |jq -Mr .body
	key         :     Get key from token and store in cache.
	                  	key [$cid] |jq -Mr .body

	siteslist   :     Make a new sites list from CSV sources list (env.SitesListCSV).

	updateusers :     Update all users of sites list.

	upsertchns  :     Upsert all channels of sites list.
	
	upsertposts :     Upsert all posts of all sites on sites list.
	
	upsertpostschron : Repeatedly run upsertposts command every x hours 
	                   	upsertpostschron $hours
	
	purgecachetkns
	purgecacheposts (*_posts.json, *_msgs.json)

	uptkn       :     Upsert a long-form message of hosted channel using JWT authentication.
	                  	uptkn $json [$jwt [$slug]]
	uptkey      :     Upsert a long-form message of hosted channel using API key authentication.
	                  	uptkn $json [$key]
	wpfetch     :      Fetch WordPress Posts from the declared URL 
	                	and dump JSON response body to file @ ./wp_posts.<DOMAIN>.json
	                	wpfetch $url

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

	case "dev":
		j := `{
			"id": 346241,
			"date": "2022-12-21T14:48:25",
			"date_gmt": "1970-01-01T00:00:00",
			"modified": "2022-12-21T14:48:25",
			"modified_gmt": "1970-01-01T00:00:00",
			"slug": "some-slug",
			"link": "https://foo.bar/baz",
			"title": {
				"rendered": "Title"
			}
		}`
		post := wordpress.Post{}
		if err := json.Unmarshal([]byte(j), &post); err != nil {
			return errors.Wrap(err, "decoding JSON message")
		}
		toRFC3339 := func(date string) time.Time {
			t, _ := time.Parse(time.RFC3339, date+"Z")
			return t.UTC()
		}
		isUnixZero := func(t time.Time) bool {
			return t == time.Unix(0, 0).UTC()
		}
		fmt.Println(convert.PrettyPrint(post))
		fmt.Println(
			toRFC3339(post.DateGMT).IsZero(),
			isUnixZero(toRFC3339(post.DateGMT)),
			toRFC3339(post.DateGMT),
		)
	case "env":
		if err := env.PrettyPrint(); err != nil {
			return err
		}
	case "upsertpostschron":
		commands.UpsertPostsChron(env, convert.ToInt(env.Args.Num(1)))

	case "siteslist":
		fname := env.SitesListJSON
		client.GhostPrint("\n=== Make new sites list\n")
		sites := wordpress.MakeSitesList(env)
		if err := env.SetCache(fname, convert.Stringify(sites)); err != nil {
			return err
		}

	case "updateusers":
		commands.UpdateUsers(env)
	case "upsertchns":
		commands.UpsertChannels(env)
	case "upsertposts":
		commands.UpsertPosts(env)
	case "purgecachetkns":
		commands.PurgeCacheTkns(env)
	case "purgecacheposts":
		commands.PurgeCachePosts(env)

	case "site":
		site := wordpress.Site{
			//URL: "https://ComicsGate.org",
			//URL: "https://TheDuran.com",
			HostURL: "https://TheCritic.co.uk",
			ChnID:   "d5750f33-a12d-4719-9600-94fcee80f487",
		}
		wp := wordpress.NewWordPress(env, &site)
		wp.SitePosts()
		// fmt.Println(convert.PrettyPrint(site))
		msgs := wp.PostsToMsgs()
		if len(msgs) == 0 {
			fmt.Fprintf(os.Stderr, "WARN : NO MESSAGES @ %s", err)
		}
		//fmt.Printf("%s", convert.Stringify(mm))
		path := env.Cache + "/" + "TheCritic.co.uk_msgs.json"
		ioutil.WriteFile(path, []byte(convert.Stringify(msgs)), 0664)

	case "trace":
		endpt := env.Args.Num(1)
		format := env.Args.Num(2)
		rsp := env.Trace(endpt, format)
		fmt.Printf("%s", rsp.Body)
		fmt.Fprintf(os.Stderr, "%s", rsp.Error)

	case "get":
		endpt := env.Args.Num(1)
		format := env.Args.Num(2)
		rsp := env.Get(endpt, format)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)

	case "posttkn":
		jwt := env.Args.Num(1)
		url := env.Args.Num(2)
		json := env.Args.Num(3)
		rsp := env.PostByTkn(jwt, url, json)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)
	case "postkey":
		key := env.Args.Num(1)
		url := env.Args.Num(2)
		json := env.Args.Num(3)
		rsp := env.PostByKey(key, url, json)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)

	case "tkn":
		fallthrough
	case "token":
		user := env.Args.Num(1)
		pass := env.Args.Num(2)
		rsp := env.Token(user, pass)
		if user == "" {
			user = env.Client.User
		}
		if user == "" {
			fmt.Fprintf(os.Stderr, "\nMissing user parameter\n")
			return nil
		}
		fname := client.CacheKeyTknPrefix + user
		if err := env.SetCache(fname, rsp.Body); err == nil {
			fmt.Printf("%s", env.GetCache(fname))
		}
		fmt.Printf("%#v", rsp)
	case "key":
		cid := env.Args.Num(1)
		rsp := env.PatchKey(cid)
		fname := client.CacheKeyTknPrefix + cid + ".json"
		if rsp.Error != "" {
			fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		} else {
			if err := env.SetCache(fname, rsp.Body); err == nil {
				fmt.Printf("%s", env.GetCache(fname))
				// k := client.ApiKey{}
				// env.GetCacheJSON(fname, &k)
				// fmt.Printf("\nk: %s\nkey:%s\n", convert.Stringify(k), k.Key)
			}
		}

	case "uptkn":
		// Upsert 1 JSON Message
		j := env.Args.Num(1)
		// jwt := env.Args.Num(2)
		// slug := env.Args.Num(3)
		msg := client.Message{}
		if err := json.Unmarshal([]byte(j), &msg); err != nil {
			return errors.Wrap(err, "decoding JSON message")
		}
		//fmt.Printf("'%s'", jwt)
		msg.Body = timestamp.TimeStringZulu(time.Now().UTC()) + " per JWT"
		//rsp := env.UpsertMsgByTkn(&msg, jwt, slug)
		rsp := env.UpsertMsgByTkn(&msg)
		fmt.Printf("\n%s\n", convert.Stringify(rsp))

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
		fmt.Printf("\n%s\n", convert.Stringify(rsp))
	case "wpfetch":
		// Fetch per WordPress site : Any endpoint : /posts, /tags, /categories, /users
		rsp := env.Get(env.Args.Num(1), client.JSON)
		fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
		fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
		fmt.Printf("%s", rsp.Body)

	// posts := []client.WordPressPost{}
	// if err := json.Unmarshal([]byte(rsp.Body), &posts); err != nil {
	// 	return errors.Wrap(err, "decoding JSON posts")
	// }
	// fmt.Printf("%s\n", convert.Stringify(posts))

	case "wpuptkn":
		site := wordpress.Site{
			//URL: "https://ComicsGate.org",
			//URL: "https://TheDuran.com",
			HostURL: "https://TheCritic.co.uk",
			ChnID:   "d5750f33-a12d-4719-9600-94fcee80f487",
		}
		wp := wordpress.NewWordPress(env, &site)

		jwt := env.Args.Num(1)
		slug := env.Args.Num(2)

		msgs := wp.PostsToMsgs()
		if len(msgs) == 0 {
			fmt.Fprintf(os.Stderr, "WARN : NO MESSAGES @ %s", err)
		}

		for _, msg := range msgs {
			rsp := env.UpsertMsgByTkn(&msg, jwt, slug)
			fmt.Printf("\n%s\n", convert.Stringify(rsp))
		}
	case "wpupkey":
		key := env.Args.Num(1)
		site := wordpress.Site{
			//URL: "https://ComicsGate.org",
			//URL: "https://TheDuran.com",
			HostURL: "https://TheCritic.co.uk",
			ChnID:   "d5750f33-a12d-4719-9600-94fcee80f487",
		}
		wp := wordpress.NewWordPress(env, &site)

		msgs := wp.PostsToMsgs()
		if len(msgs) == 0 {
			fmt.Fprintf(os.Stderr, "WARN : NO MESSAGES @ %s", err)
		}

		for _, msg := range msgs {
			rsp := env.UpsertMsgByKey(&msg, key)
			fmt.Printf("\n%s\n", convert.Stringify(rsp))
		}

	default: // Ghost print so pipe okay: ... |jq .
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "%s\n", DESCRIBE)

		return nil
	}

	return nil
}
