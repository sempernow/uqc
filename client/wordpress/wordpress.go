// Package wordpress handles REST APIs of WordPress and Uqrate
// to process []wordpress.Post into []client.Message for upsert at Uqrate.
package wordpress

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sempernow/kit/convert"
	"github.com/sempernow/kit/id"
	"github.com/sempernow/kit/str"
	"github.com/sempernow/uqc/client"
)

// NewWordPress contains app environment and per-site configuration.
func NewWordPress(env *client.Env, site *Site) *WP {
	return &WP{
		Env:  env,
		Site: site,
	}
}

// MakeSitesList creates []Sites from a CSV file.
// Those values are the export of an SQL query (hosts_channels.sql)
// for relevant records (users and channels) in Uqrate data store.
func MakeSitesList(env *client.Env) []Site {

	sites := []Site{}

	// Open the sites-list CSV file

	env.Logger.Printf("INFO : Try read %s from Docker config\n", PathCfgSitesListCSV)
	bb, err := os.ReadFile(PathCfgSitesListCSV)
	if err != nil {
		env.Logger.Printf("INFO : Try read %s from cache\n", env.SitesListCSV)
		bb, err = os.ReadFile(filepath.Join(env.Assets, env.SitesListCSV))
		if err != nil {
			env.Logger.Printf("ERR @ ReadFile : %s\n", err.Error())
			return sites
		}
	}
	r := csv.NewReader(bytes.NewReader(bb))

	// Append each CSV record to sites list.
	for {
		cc, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			sites[0].Error = err.Error()
			return sites
		}
		if len(cc) < 5 {
			sites[0].Error = errors.New("malformed CSV : too few fields").Error()
			return sites
		}
		if cc[1] == "slug" {
			continue
		}
		// Get additional site params dynamically and merge with those from CSV.
		wp := NewWordPress(env, &Site{HostURL: cc[2]})
		wp.Site = &Site{
			UserHandle: cc[0],
			ChnSlug:    cc[1],
			HostURL:    cc[2],
			OwnerID:    cc[3],
			ChnID:      cc[4],
		}
		wp.SiteGot()

		sites = append(sites, *wp.Site)
	}
	return sites
}

// GetSitesList retrieves []Sites list from its cache (JSON)
// if exist, else makes and caches anew.
func GetSitesList(env *client.Env) []Site {
	sites := []Site{}
	env.Logger.Printf("INFO : Try read %s from Docker config\n", PathCfgSitesListJSON)
	j, err := os.ReadFile(PathCfgSitesListJSON)
	if err != nil {
		env.Logger.Printf("INFO : Try read %s from cache\n", env.SitesListJSON)
		j = env.GetCache(env.SitesListJSON)
		if len(j) == 0 {
			env.Logger.Printf("INFO : Make new sites list\n")
			sites = MakeSitesList(env)
			if err := env.SetCache(env.SitesListJSON, convert.Stringify(sites)); err != nil {
				env.Logger.Printf("ERR : setting cache\n")
				return sites
			}
			j = env.GetCache(env.SitesListJSON)
		}
	}
	if len(j) == 0 {
		env.Logger.Printf("ERR : sites list (%s) EMPTY or NOT FOUND\n", env.SitesListJSON)
	}
	if err := json.Unmarshal(j, &sites); err != nil {
		env.Logger.Printf("ERR : unmarshalling json of '%s'\n", env.SitesListJSON)
		return sites
	}
	return sites
}

// SiteGot retrieves dynamic fields of Site from a site,
// and merges it into existing site record (wp.Site) by reference.
func (wp WP) SiteGot() {

	j, err := wp.getWP(SiteURI)

	if err != nil {
		wp.Site.Error = err.Error()
		return
	}
	if j == "" {
		wp.Site.Error = "GET returned nothing"
		return
	}
	if err := json.Unmarshal([]byte(j), &wp.Site); err != nil {
		wp.Site.Error = err.Error()
		log.Printf("ERR : Unmarshalling : %s\n", err.Error())
	}
}

// SitePosts retrieves wp.Site.Posts; the WordPress-normalized []Post list from a Site.
func (wp WP) SitePosts() {
	j, err := wp.getWP(PostsURI)
	if err != nil {
		wp.Site.Error = err.Error()
		return
	}
	if j == "" {
		wp.Site.Error = "GET returned nothing"
		return
	}

	if err := json.Unmarshal([]byte(j), &wp.Site.Posts); err != nil {
		wp.Site.Error = err.Error()
		log.Printf("ERR : Unmarshalling : %s\n", err.Error())
	}
}

// GetTkn retrieves JWT for env.Client.User; get from cache; fetch on miss.
func (wp WP) GetTkn() string {
	key := client.CacheKeyTknPrefix + wp.Env.Client.User
	bb := wp.Env.GetCache(key)
	if len(bb) == 0 {
		//log.Printf("INFO : cache miss @ %s\n", key)
		rsp := wp.Env.Token()
		if rsp.Code != 200 {
			log.Printf("ERR : Token(..) %s : %s\n", wp.Env.Client.User, rsp.Error)
			return ""
		}
		if err := wp.Env.SetCache(key, rsp.Body); err != nil {
			log.Printf("ERR : GetTkn: %s : %s\n", wp.Env.Client.User, err.Error())
			return ""
		}
		return rsp.Body
	}
	return convert.BytesToString(bb)
}

// getWP retrieves response (JSON) of a WordPress API endpoint; get from cache; fetch on miss.
func (wp WP) getWP(uri string) (string, error) {
	url := wp.Site.HostURL + uri
	key := urlToFname(url)

	// First try cache.
	bb := wp.Env.GetCache(key)
	if len(bb) == 0 {
		//log.Printf("INFO : cache miss @ %s\n", key)

		// Hit the site softly
		rsp := wp.Env.Get(url, client.JSON)
		time.Sleep(time.Millisecond * 300)
		wp.Site.Status.Object = uri
		wp.Site.Status.Code = rsp.Code

		if rsp.Error != "" {
			// Write regardless to prevent future fetches
			wp.Env.SetCache(key, "")
			return "", errors.New(rsp.Error)
		}
		wp.Env.SetCache(key, rsp.Body)

		return rsp.Body, nil
	}
	return convert.BytesToString(bb), nil
}

// Posts2Msgs denormalizes a WordPress post into a Uqrate message.
func (wp WP) PostsToMsgs() []client.Message {
	list := []client.Message{}
	for _, post := range wp.Site.Posts {
		list = append(list, wp.PostToMsg(&post))
	}
	return list
}

// Posts2Msgs denormalizes a *Post into a client.Message,
// retrieving the various WordPress objects (referenced at Post subkeys)
// as needed to populate Message keys (.Cats, .Tags).
// Message.ID is a static UUID (v5) per Message.ChnID namespace and Message.URI name.
func (wp WP) PostToMsg(post *Post) client.Message {
	msg := client.Message{}

	msg.ChnID = wp.Site.ChnID
	msg.URI = linkToURI(post.Link)
	//msg.ID = uuid.NewV5(uuid.Must(uuid.FromString(msg.ChnID)), strings.ToLower(msg.URI)).String()
	msg.ID, _ = id.UUIDv5(msg.ChnID, msg.URI)
	if msg.ID == "" {
		log.Printf("ERR : UUIDv5 fail : URI: %s .\n", msg.URI)
		return client.Message{URI: msg.URI}
	}

	msg.Title = post.Title.Rendered
	msg.Body = post.Content.Rendered

	if post.Excerpt.Rendered != "" {
		msg.Summary = post.Excerpt.Rendered
	}

	if true {
		if len(post.Categories) > 0 {
			msg.Cats = wp.objNameList(CatsURI, post.Categories)
		}
		if len(post.Tags) > 0 {
			msg.Tags = wp.objNameList(TagsURI, post.Tags)
		}
	}
	// Add the author's name to the list of tags for this message.
	uri := appendToURL(AuthorsURI, convert.IntToString(post.Author))
	if author := wp.objName(uri); author != "" {
		if !strings.Contains(author, "s") {
			msg.Tags = append(msg.Tags, author)
		}
	}

	sanitize(msg.Cats)
	sanitize(msg.Tags)

	// Recover the post timestamp

	off := wp.Site.GMTOffset
	msg.DateUpdate = ToRFC3339(post.ModifiedGMT, 0)
	if IsUnixZero(msg.DateUpdate) {
		msg.DateUpdate = ToRFC3339(post.Modified, off)
	}
	if IsUnixZero(msg.DateUpdate) {
		msg.DateUpdate = ToRFC3339(post.DateGMT, 0)
	}
	if IsUnixZero(msg.DateUpdate) {
		msg.DateUpdate = ToRFC3339(post.Date, off)
	}
	if IsUnixZero(msg.DateUpdate) || msg.DateUpdate.IsZero() {
		msg.DateUpdate = time.Now().Truncate(1 * time.Second).UTC()
		log.Printf("WARN : msg.DateUpdate : NOT FOUND : Set to current time.\n")
	}
	log.Printf("INFO : msg.DateUpdate : %v : GMT Offset was: %d\n", msg.DateUpdate, off)

	return msg
}

// sanitize each name of list.
func sanitize(names []string) {
	for i, name := range names {
		names[i] = str.CleanAlphaNum(name, 35)
	}
}

type object struct {
	ID   int
	Name string
}

// objName retrieves the name referenced (by ID) in a WordPress Post,
// per object type (.Author, .Categories, .Tags), from its (API) URI.
func (wp WP) objName(uri string) string {
	j, err := wp.getWP(uri)
	if err != nil {
		return ""
	}
	o := object{}
	if err := json.Unmarshal([]byte(j), &o); err != nil {
		return ""
	}
	return o.Name
}

// objNameList retrieves the list of names referenced (by ID) in a WordPress Post,
// per object type (.Author, .Categories, .Tags), by its reference list, from its (API) URI.
func (wp WP) objNameList(uri string, want []int) []string {

	j, err := wp.getWP(uri)

	if err != nil {
		return []string{}
	}
	oo := []object{}
	if err := json.Unmarshal([]byte(j), &oo); err != nil {
		return []string{}
	}
	var (
		names []string
		miss  []int
		got   bool
	)
	// Match id:name
	for _, id := range want {
		for _, o := range oo {
			if got {
				continue
			}
			if id == o.ID {
				names = append(names, o.Name)
				got = true
			}
		}

		if !got {
			miss = append(miss, id)
		}
		got = false
	}
	// If any want are missing, then get per id
	if len(miss) > 0 {
		//log.Printf("INFO : name(s) miss (%d) @ uri: %s\n", len(miss), uri)
		static := uri
		for _, id := range miss {
			uri = appendToURL(uri, convert.IntToString(id))

			if name := wp.objName(uri); name != "" {
				names = append(names, name)
			}
			uri = static
		}
	}
	return names
}

// appendToURL(..) de/re/constructs URL as necessary to append a slug (a).
//
//	"../foo/bar?x=b,c&d=7" + a => "../foo/bar/a?x=b,c&d=7"
//	"../foo/bar"           + a => "../foo/bar/a"
//	"../foo/bar/"          + a => "../foo/bar/a/"
func appendToURL(url, a string) string {
	ss := strings.Split(url, "?")

	if len(ss) == 2 {
		url = ss[0] + "/" + a
		url = url + "?" + ss[1]
	}

	if len(ss) == 1 {
		if url[len(url)-1:] == "/" {
			url = url + a + "/"
		} else {
			url = url + "/" + a
		}
	}
	return url
}

// urlToFname converts url to rtn (cache fname), e.g.,
//
//	url : "https://TheWpSite.com/wp-json/wp/v2/posts?author=7"
//	rtn : "TheWpSite.com_posts.json"
//	url : "https://TheWpSite.com/wp-json/wp/v2/users/7"
//	rtn : "TheWpSite.com_users.7.json"
func urlToFname(url string) string {
	site := fqdn(url)
	var obj, fname string
	if strings.Contains(url, "wp-json/wp/v2") {
		obj = strings.Split(url, "wp-json/wp/v2")[1]
	} else {
		if strings.Contains(url, "wp-json/") {
			obj = strings.Split(url, "wp-json/")[1]
		}
	}
	obj = strings.Split(obj, "?")[0]
	ss := strings.Split(obj, "/")
	if len(ss) > 1 {
		fname = site + "_"
	} else {
		fname = site + "."
	}
	for _, s := range ss {
		if s != "" {
			fname += s + "."
		}
	}
	return fname + "json"
}

// url : "https://TheWpSite.com/wp-json/wp/v2/posts?author=7"
// rtn : "TheWpSite.com"
func fqdn(url string) string {
	ss := strings.Split(url, "//")
	fname := ss[0]
	if len(ss) > 1 {
		fname = ss[1]
	}
	ss = strings.Split(fname, "/")
	return ss[0]
}

// "https://foo.bar.baz/a/b" to "/a/b"
func linkToURI(link string) string {
	x := strings.SplitAfter(link, "://")
	x = strings.SplitAfter(x[1], "/")
	return "/" + strings.Join(x[1:], "")
}

// ToRFC3339 parses WordPress-posts ($post) date into Golang time (GMT).
func ToRFC3339(date string, offset int) time.Time {
	// ToRFC3339("2021-12-21T13:15:31", -6)
	//==>         2021-12-21 19:15:31 +0000 UTC
	// WP REST API @ /wp-json :
	// "gmt_offset": -6
	// WP REST API @ /wp-json/wp/v2/posts :
	// "date":         "2021-12-21T13:15:31",
	// "date_gmt":     "2021-12-21T19:15:31",
	// "modified":     "2022-01-11T14:22:07",
	// "modified_gmt": "2022-01-11T20:22:07",
	off := "Z"
	if offset < 0 {
		if offset > -26 {
			off = fmt.Sprintf("%.2d:00", offset)
		}
	} else {
		if offset < 26 {
			off = fmt.Sprintf("+%.2d:00", offset)
		}
	}
	t, _ := time.Parse(time.RFC3339, (date + off))
	return t.Truncate(1 * time.Second).UTC()
}

// IsUnixZero tests for "1970-01-01 00:00:00 +0000 UTC".
// Unlike time pkg t.IsZero(), which tests for "0001-01-01 00:00:00 +0000 UTC".
func IsUnixZero(t time.Time) bool {
	return t == time.Unix(0, 0).UTC()
}
