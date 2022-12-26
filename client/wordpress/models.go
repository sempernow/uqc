package wordpress

import "github.com/sempernow/uqc/client"

const DateZeroWP = "1970-01-01T00:00:00"

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

	// Endpoint : /wp-json
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Home        string `json:"home,omitempty"`
	GMTOffset   int    `json:"gmt_offset,omitempty"`
}

type Status struct {
	Object string `json:"object,omitempty"`
	Code   int    `json:"code,omitempty"`
}

// Docker configs paths
const (
	PathCfgSitesListJSON = "/sites_list_json"
	PathCfgSitesListCSV  = "/sites_list_csv"
)

// WordPress REST API endpoints
// https://developer.wordpress.org/rest-api/reference/
const (
	SiteURI    = "/wp-json/?_fields=name,description,url,home,gmt_offset"
	PostsURI   = "/wp-json/wp/v2/posts?_fields=id,date,date_gmt,link,modified,modified_gmt,slug,GUID,title,content,excerpt,author,categories,tags,comment_status"
	TagsURI    = "/wp-json/wp/v2/tags?_fields=id,name,slug,count&per_page=100"
	CatsURI    = "/wp-json/wp/v2/categories?_fields=id,name,slug,count&per_page=100"
	AuthorsURI = "/wp-json/wp/v2/users?_fields=id,name,slug,avatar_urls&per_page=100"
)

// Post contains a subset of keys from its WordPress
// REST API namesake of the Posts endpoint.
// https://developer.wordpress.org/rest-api/reference/posts/
type Post struct {
	ID          int      `json:"id,omitempty"`
	Date        string   `json:"date,omitempty"`         // @ New
	DateGMT     string   `json:"date_gmt,omitempty"`     // @ New
	Modified    string   `json:"modified,omitempty"`     // @ Edit
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
