package commands

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

// WordPressPost abides its subset of JSON response body
// keys from WordPress API endpoint /wp-json/wp/v2/posts .
// https://developer.wordpress.org/rest-api/reference
// https://developer.wordpress.org/rest-api/reference/posts/
type WordPressPost struct {
	ID          int      `json:"id,omitempty"`
	DateGMT     string   `json:"date_gmt,omitempty"`
	ModifiedGMT string   `json:"modified_gmt,omitempty"`
	GUID        Rendered `json:"guid,omitempty"`    // https://base.com?p=29343
	Link        string   `json:"link,omitempty"`    // https://base.com/2022/09/title-string
	Slug        string   `json:"slug,omitempty"`    // /title-string
	Title       Rendered `json:"title,omitempty"`   // Subkey: Rendered
	Content     Rendered `json:"content,omitempty"` // Subkey: Rendered
	Excerpt     Rendered `json:"excerpt,omitempty"` // Subkey: Rendered
	Author      int      `json:"author,omitempty"`
	Categories  []int    `json:"categories,omitempty"` //
	Tags        []int    `json:"tags,omitempty"`       // @

	// Code  int    `json:"code,omitempty"`
	// Error string `json:"error,omitempty"`
}

type Rendered struct {
	Rendered string `json:"rendered,omitempty"`
}

// ConvertPosts maps a WordPress-API /posts response body (JSON)
// into a Uqrate []messsage.Message list.
// midSrc is either user_id + chn_slug or chn_id
func Posts2Msgs(j, midSrc string) ([]client.Message, error) {

	posts := []WordPressPost{}
	if err := json.Unmarshal([]byte(j), &posts); err != nil {
		return []client.Message{}, errors.Wrap(err, "decoding JSON posts")
	}
	msg := []client.Message{}
	for _, post := range posts {
		msg = append(msg, convertPost(&post, midSrc))
	}

	return msg, nil
}

// ConvertPost maps a WordPressPost to a message.Message of Uqrate
func convertPost(post *WordPressPost, midSrc string) client.Message {
	// Filter by post.Author ???
	msg := client.Message{}

	uri := link2uri(post.Link)
	msg.ID = uuid.NewV5(uuid.NamespaceOID, (midSrc + uri)).String()
	msg.URI = uri

	msg.Title = post.Title.Rendered
	msg.Body = post.Content.Rendered

	if post.Excerpt.Rendered != "" {
		msg.Summary = post.Excerpt.Rendered
	}

	// TODO : See uqrate project : flatten and segregate into tags and jt_tags_xid table
	if false {
		if len(post.Tags) > 0 {
			msg.Tags = getTags(post.Tags)
		}
		if len(post.Categories) > 0 {
			msg.Tags = getCats(post.Categories)
		}
	}

	msg.DateUpdate, _ = time.Parse(time.RFC3339, post.ModifiedGMT+"Z")

	return msg
}

func getTags(tags []int) []string {
	// Retrieve from our WP data store
	// else fetch from WordPress site : GET /wp-json/wp/v2/tags/<id>
	k := []string{}
	for _, tag := range tags {
		k = append(k, types.IntToString(tag))
	}
	return k
}
func getCats(cats []int) []string {
	// Retrieve from our WP data store
	// else fetch from WordPress site : GET /wp-json/wp/v2/categories/<id>
	c := []string{}
	for _, cat := range cats {
		c = append(c, types.IntToString(cat))
	}
	return c
}

// https://foo.bar.baz/a/b/c => /a/b/c
func link2uri(link string) string {
	x := strings.SplitAfter(link, "://")
	x = strings.SplitAfter(x[1], "/")
	return "/" + strings.Join(x[1:], "")
}
func toRFC3339(date string) string {
	return date + "Z"
}
