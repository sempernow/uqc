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
	"github.com/sempernow/uqc/kit/types"
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
			msg.Cats = CatsList(post.Categories)
		}
		if len(post.Tags) > 0 {
			msg.Tags = TagsList(post.Tags)
		}
	}
	if author := Author(post.Author); author != "" {
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

const (
	Posts      = "/wp-json/wp/v2/posts"
	Tags       = "/wp-json/wp/v2/tags"
	Categories = "/wp-json/wp/v2/categories"
	Authors    = "/wp-json/wp/v2/users"
)

// Post contains a subset of keys from its WordPress
// REST API namesake at endpoint /wp-json/wp/v2/posts .
// https://developer.wordpress.org/rest-api/reference/posts/
type Post struct {
	ID          int      `json:"id,omitempty"`
	DateGMT     string   `json:"date_gmt,omitempty"`     // @ New
	ModifiedGMT string   `json:"modified_gmt,omitempty"` // @ Edit
	Link        string   `json:"link,omitempty"`         // https://base.com/2022/09/title-string
	Slug        string   `json:"slug,omitempty"`         // /title-string
	GUID        Rendered `json:"guid,omitempty"`         // https://base.com?p=29343
	Title       Rendered `json:"title,omitempty"`
	Content     Rendered `json:"content,omitempty"`
	Excerpt     Rendered `json:"excerpt,omitempty"`
	Author      int      `json:"author,omitempty"`
	Categories  []int    `json:"categories,omitempty"`
	Tags        []int    `json:"tags,omitempty"`

	// Code  int    `json:"code,omitempty"`
	// Error string `json:"error,omitempty"`
}

// Rendered contains that of certain Post keys.
type Rendered struct {
	Rendered string `json:"rendered,omitempty"`
}

func Author(id int) string {
	// Retrieve from /wp-json/wp/v2/authors/:id
	return "author name"
}

func PostsList(env *client.Env, site string) []Post {

	fetchWP(env, site+Posts)

	return []Post{}
}
func TagsList(tags []int) []string {
	// Retrieve from our WP data store
	// else fetch from WordPress site : GET /wp-json/wp/v2/tags/:id
	k := []string{}
	for _, tag := range tags {
		k = append(k, types.IntToString(tag))
	}
	return k
}

func CatsList(cats []int) []string {
	// Retrieve from our WP data store
	// else fetch from WordPress site : GET /wp-json/wp/v2/categories/<id>
	c := []string{}
	for _, cat := range cats {
		c = append(c, types.IntToString(cat))
	}
	return c
}

func fetchWP(env *client.Env, url string) (string, error) {
	rsp := env.Dump(url, client.JSON)
	if rsp.Error != "" {
		return "", errors.New(rsp.Error)
	}
	return rsp.Body, nil
}
