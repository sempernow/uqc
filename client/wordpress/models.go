package wordpress

import "github.com/sempernow/uqc/client"

// WP contains a WordPress site configuration
type WP struct {
	Env     *client.Env
	Site    *Site
	Cleaner func(string) string
}

// Site contains that required to map a WordPress site to a Uqrate Channel.
type Site struct {
	UserHandle string `json:"user_handle,omitempty"`
	ChnSlug    string `json:"chn_slug,omitempty"`
	HostURL    string `json:"host_url,omitempty"`
	OwnerID    string `json:"owner_id,omitempty"`
	ChnID      string `json:"chn_id,omitempty"`
	Posts      []Post `json:"posts,omitempty"`
	Error      string `json:"error,omitempty"`
	Status     `json:"status,omitempty"`
}

type Status struct {
	Object string `json:"object,omitempty"`
	Code   int    `json:"code,omitempty"`
}

const (
	PostsURI   = "/wp-json/wp/v2/posts"
	TagsURI    = "/wp-json/wp/v2/tags?per_page=50&_fields=id,name,slug"
	CatsURI    = "/wp-json/wp/v2/categories?per_page=50&_fields=id,name,slug"
	AuthorsURI = "/wp-json/wp/v2/users?_fields=id,name,slug,avatar_urls.96"
)

// Post contains a subset of keys from its WordPress
// REST API namesake of the Posts endpoint.
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
	Error string `json:"error,omitempty"`
}

// Rendered contains that of certain Post keys.
type Rendered struct {
	Rendered string `json:"rendered,omitempty"`
}
