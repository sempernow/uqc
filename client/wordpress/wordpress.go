// Package wordpress provides access to the WordPress REST API,
// transforming JSON response bodies of relevant endpoints into []client.Message .
package wordpress

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sempernow/uqc/client"
)

func Test(env *client.Env) {
	fmt.Println(env)
}

// Posts2Msgs transforms a []Post response body (JSON)
// into a []client.Message of a Uqrate Channel (cid).
func Posts2Msgs(j, cid string) ([]client.Message, error) {

	posts := []Post{}
	if err := json.Unmarshal([]byte(j), &posts); err != nil {
		return []client.Message{}, errors.Wrap(err, "decoding JSON posts")
	}
	msg := []client.Message{}
	for _, post := range posts {
		msg = append(msg, convertPost(&post, cid))
	}

	return msg, nil
}

func FetchWP(env *client.Env, endpt string) (string, error) {
	rsp := env.Dump(endpt, client.JSON)
	// fmt.Fprintf(os.Stderr, "HTTP %d\n", rsp.Code)
	// fmt.Fprintf(os.Stderr, "%s\n", rsp.Error)
	// fmt.Printf("%s", rsp.Body)
	if rsp.Error != "" {
		return "", errors.New(rsp.Error)
	}
	return rsp.Body, nil
}

// ConvertPost maps a Post to a client.Message
// having a UUIDv5 ID sourced from client.Message.ChnID (cid).
func convertPost(post *Post, cid string) client.Message {
	msg := client.Message{}

	msg.ChnID = cid
	msg.URI = link2uri(post.Link)
	msg.ID = uuid.NewV5(uuid.NamespaceOID, (msg.ChnID + msg.URI)).String()

	msg.Title = post.Title.Rendered
	msg.Body = post.Content.Rendered

	if post.Excerpt.Rendered != "" {
		msg.Summary = post.Excerpt.Rendered
	}

	if false {
		if len(post.Categories) > 0 {
			msg.Cats = fetch.CatsList(post.Categories)
		}
		if len(post.Tags) > 0 {
			msg.Tags = getTags(post.Tags)
		}
	}
	if author := getAuthor(post.Author); author != "" {
		msg.Tags = append(msg.Tags, author)
	}

	msg.DateUpdate = toRFC3339(post.ModifiedGMT)

	return msg
}

// https://foo.bar.baz/a/b/c => /a/b/c
func link2uri(link string) string {
	x := strings.SplitAfter(link, "://")
	x = strings.SplitAfter(x[1], "/")
	return "/" + strings.Join(x[1:], "")
}
func toRFC3339(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date+"Z")
	return t
}
