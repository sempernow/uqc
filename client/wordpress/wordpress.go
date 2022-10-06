// Package wordpress provides access to the WordPress REST API,
// transforming JSON response bodies of relevant endpoints into []client.Message .
package wordpress

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/kit/convert"
	"github.com/sempernow/uqc/kit/str"
)

func NewWordPress(env *client.Env, site *Site) *WP {
	return &WP{
		Env:  env,
		Site: site,
	}
}

// SiteList creates []Sites from a CSV file (host_channels.csv).
// Those values are are the export of an SQL query (hosts_channels.sql)
// for site-host records; uqrate users, each having an associated channel.
func SiteList(env *client.Env) []Site {
	sites := []Site{}
	bb, err := ioutil.ReadFile(filepath.Join(env.Assets, "wp", "host_channels.csv"))
	if err != nil {
		sites[0].Error = err.Error()
		return sites
	}
	r := csv.NewReader(bytes.NewReader(bb))
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
		// [
		// 	NewDiscourses
		// 	Mirror
		// 	https://NewDiscourses.com
		// 	fd8bf171-5c81-4435-b3cd-3fb2d1d09f1b
		// 	8d996155-8ec7-43c0-a506-d6a883267798
		// ]

		sites = append(sites, Site{
			UserHandle: cc[0],
			ChnSlug:    cc[1],
			HostURL:    cc[2],
			OwnerID:    cc[3],
			ChnID:      cc[4],
		})
	}
	return sites
}

// UpsertSites converts and upserts all Posts of each site in []Site list.
func UpsertSites(env *client.Env, sites []Site) {
	env.Channel.Slug = "Mirror"
	env.Client.Pass = env.SitesPass

	var (
		wp    *WP
		msgs  []client.Message
		rsp   *client.Response
		fname string
		err   error
	)
	for i, site := range sites {
		if site.ChnSlug == "slug" {
			continue
		}

		fmt.Printf("%d : '%s'\n", i, site.UserHandle)
		wp = NewWordPress(env, &site)
		wp.SitePosts()
		msgs, err = wp.PostsToMsgs()
		if err != nil {
			client.GhostPrint("\nERR @ PostsToMsgs: %s : %s\n", site.UserHandle, err.Error())
			continue
		}
		env.SetCache("msgs_"+site.UserHandle+".json", convert.Stringify(msgs))

		// Get authorization for upsert of this user's channel
		env.Client.User = site.UserHandle
		rsp = env.Token()
		fname = "/keys/tkn." + site.UserHandle
		if err := env.SetCache(fname, rsp.Body); err != nil {
			client.GhostPrint("\nERR @ Token: %s : %s\n", site.UserHandle, err.Error())
			continue
		}

		for _, msg := range msgs {
			rsp := env.UpsertMsgByTkn(&msg)
			fmt.Printf("\nupsert messages @ %s : %s\n", site.UserHandle, convert.Stringify(rsp))
		}
	}
}

// SitePosts retrieves the WordPress-normalized []Post from a Site.
func (wp WP) SitePosts() {
	j, err := wp.getWP(PostsURI)
	if err != nil {
		wp.Site.Error = err.Error()
	}
	if err := json.Unmarshal([]byte(j), &wp.Site.Posts); err != nil {
		wp.Site.Error = err.Error()
	}
}

// getWP retrieves cached response (JSON); fetches only on miss.
func (wp WP) getWP(uri string) (string, error) {
	// First try cache.
	url := wp.Site.HostURL + uri
	path := filepath.Join(wp.Env.Cache, urlToFname(url))
	jj, miss := ioutil.ReadFile(path)
	if miss != nil {
		fmt.Fprintf(os.Stderr, "Fetch @ cache miss: %s\npath: %s\n", url, path)

		if strings.Contains(uri, wp.Site.Status.Object) && wp.Site.Status.Code > 299 {
			return "", errors.New("site status: " + convert.IntToString(wp.Site.Status.Code))
		}

		rsp := wp.Env.Get(url, client.JSON)
		time.Sleep(time.Millisecond * 300)

		wp.Site.Status.Object = uri
		wp.Site.Status.Code = rsp.Code
		if rsp.Error != "" {
			// Write regardless to prevent recurring fetch on fail
			ioutil.WriteFile(path, []byte(""), 0664)
			return "", errors.New(rsp.Error)
		}
		ioutil.WriteFile(path, []byte(rsp.Body), 0664)
		return rsp.Body, nil
	}
	return convert.BytesToString(jj), nil
}

// Posts2Msgs denormalizes a []Post into a []client.Message of a Uqrate Channel.
func (wp WP) PostsToMsgs() ([]client.Message, error) {
	list := []client.Message{}
	for _, post := range wp.Site.Posts {
		list = append(list, wp.PostToMsg(&post))
	}
	return list, nil
}

// Posts2Msgs denormalizes a *Post into a client.Message of a Uqrate Channel.
// The message ID is a UUID v5 of OID namespace; name is concat of channel ID and message URI.
func (wp WP) PostToMsg(post *Post) client.Message {
	msg := client.Message{}

	msg.ChnID = wp.Site.ChnID
	msg.URI = linkToURI(post.Link)
	msg.ID = uuid.NewV5(uuid.NamespaceOID, (msg.ChnID + msg.URI)).String()

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
	uri := filepath.Join(AuthorsURI, convert.IntToString(post.Author))
	if author := wp.objName(uri); author != "" {
		if !strings.Contains(author, "s") {
			msg.Tags = append(msg.Tags, author)
		}
	}

	makeValid(msg.Cats)
	makeValid(msg.Tags)

	msg.DateUpdate = asRFC3339(post.ModifiedGMT)

	return msg
}

func makeValid(names []string) {
	for i, name := range names {
		names[i] = str.CleanAlphaNum(name, 35)
	}
}

type object struct {
	ID   int
	Name string
}

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
		got   bool
		miss  []int
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
		for _, id := range miss {
			if name := wp.objName(uri + "/" + convert.IntToString(id)); name != "" {
				names = append(names, name)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
	return names
}

// urlToFname converts url to rtn (cache fname), e.g.,
// 	url : "https://TheWpSite.com/wp-json/wp/v2/posts?author=7"
// 	rtn : "TheWpSite.com_posts.json"
// 	url : "https://TheWpSite.com/wp-json/wp/v2/users/7"
// 	rtn : "TheWpSite.com_users.7.json"
func urlToFname(url string) string {
	site := fqdn(url)
	obj := strings.Split(url, "wp-json/wp/v2")[1]
	obj = strings.Split(obj, "?")[0]
	ss := strings.Split(obj, "/")

	fname := site + "_"
	for _, s := range ss {
		if s != "" {
			fname += s + "."
		}
	}
	return fname + "json"
}

// 	url : "https://TheWpSite.com/wp-json/wp/v2/posts?author=7"
// 	rtn : "TheWpSite.com"
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

// asRFC339 parses WordPress ($post) timestamp string into Golang time.
func asRFC3339(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date+"Z")
	return t
}
