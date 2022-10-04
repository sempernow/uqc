package wordpress

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
