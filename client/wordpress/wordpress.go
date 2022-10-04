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
	"github.com/sempernow/uqc/client/wordpress/wpfetch"
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
			msg.Cats = wpfetch.CatsList(post.Categories)
		}
		if len(post.Tags) > 0 {
			msg.Tags = wpfetch.TagsList(post.Tags)
		}
	}
	if author := wpfetch.Author(post.Author); author != "" {
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
