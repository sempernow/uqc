// Package commands provides high-level processing of a WordPress SitesList.
package commands

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/sempernow/kit/convert"
	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/client/wordpress"
)

const (
	SUFFIX_POSTS = "_posts.json"
	SUFFIX_MSGS  = "_msgs.json"
)

// UpsertChannels of sites list with values therein.
func UpsertChannels(env *client.Env) {
	sites := wordpress.GetSitesList(env)
	env.Client.Pass = env.SitesPass
	for _, site := range sites {
		wp := wordpress.NewWordPress(env, &site)
		env.Client.User = site.UserHandle
		tkn := wp.GetTkn()
		if tkn == "" {
			continue
		}
		chn := client.Channel{
			ID:      site.ChnID,
			OwnerID: site.OwnerID,
			Slug:    site.ChnSlug,
			//Tags:    []string{},        //... mutates per message upsert.
			//Title:   site.Name,         //... set @ User (Site) record.
			//About:   site.Description,  //... set @ User (Site) record.
		}
		rsp := wp.Env.PostByTkn(tkn, env.Service.BaseAPI+"/c/upsert", &chn)
		if rsp.Code > 299 {
			env.Logger.Printf("ERR : PostByTkn @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		} else {
			env.Logger.Printf("INFO : PostByTkn @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		}
	}
}

// UpdateUsers of sites list with values therein.
func UpdateUsers(env *client.Env) {
	sites := wordpress.GetSitesList(env)

	// All sites mirrored hereby share common password
	env.Client.Pass = env.SitesPass

	for _, site := range sites {
		wp := wordpress.NewWordPress(env, &site)
		env.Client.User = site.UserHandle

		tkn := wp.GetTkn()

		if tkn == "" {
			continue
		}

		// Get/Set avatar and banner

		var (
			avatar = "-avatar.webp"
			banner = "-banner.webp"
		)
		if _, err := os.ReadFile(
			filepath.Join(env.Assets, "media", "avatars", (site.UserHandle + avatar)),
		); err != nil {
			avatar = "wordpress" + avatar
		} else {
			avatar = site.UserHandle + avatar
		}
		if _, err := os.ReadFile(
			filepath.Join(env.Assets, "media", "banners", (site.UserHandle + banner)),
		); err != nil {
			banner = "uqrate" + banner
		} else {
			banner = site.UserHandle + banner
		}

		// Set payload

		user := client.User{
			Display: site.Name,
			About:   site.Description,
			Avatar:  avatar,
			Banner:  banner,
		}
		if len(user.Display) > client.MaxUserDisplay {
			if len(user.About) == 0 {
				user.About = user.Display
			}
			user.Display = user.Display[:client.MaxUserDisplay]
		}

		// Update site (user) record
		rsp := wp.Env.PutByTkn(tkn, env.Service.BaseAPI+"/u/"+site.OwnerID, &user)
		if rsp.Code > 299 {
			env.Logger.Printf("ERR : PutByTkn @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		} else {
			env.Logger.Printf("INFO : PutByTkn @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		}
	}
}

// PurgeCacheTkns removes token cache.
func PurgeCacheTkns(env *client.Env) {
	env.Logger.Printf("INFO : PurgeCacheTkns ("+client.CacheKeyTknPrefix+"*) @ %s\n", env.Cache)
	sites := wordpress.GetSitesList(env)
	for _, site := range sites {
		fname := client.CacheKeyTknPrefix + site.UserHandle
		if err := os.Remove(filepath.Join(env.Cache, fname)); err != nil {
		} else {
			env.Logger.Printf("INFO : DEL @ %s\n", fname)
		}
	}
}

// PurgeCachePosts removes posts and messages cache.
func PurgeCachePosts(env *client.Env) {
	env.Logger.Printf("INFO: PurgeCachePosts @ %s\n", env.Cache)
	sites := wordpress.GetSitesList(env)
	for _, site := range sites {
		domain := strings.Split(site.HostURL, "//")[1]
		fname := domain + SUFFIX_POSTS
		if err := os.Remove(filepath.Join(env.Cache, fname)); err != nil {
		} else {
			env.Logger.Printf("INFO : DEL @ %s\n", fname)
		}
		fname = domain + SUFFIX_MSGS
		if err := os.Remove(filepath.Join(env.Cache, fname)); err != nil {
		} else {
			env.Logger.Printf("INFO : DEL @ %s\n", fname)
		}
	}
}

// UpsertPosts converts []Post into []client.Message of all sites in []Site list,
// upserting the Uqrate messages to their associated channel (mirror) per site.
func UpsertPosts(env *client.Env) {
	PurgeCacheTkns(env)
	PurgeCachePosts(env)
	sites := wordpress.GetSitesList(env)
	env.Channel.Slug = "Mirror"
	env.Client.Pass = env.SitesPass
	var (
		wp   *wordpress.WP
		msgs []client.Message
	)

	// Process each site in sites list

	for i, site := range sites {
		if site.ChnSlug == "slug" {
			// @ CSV header (first row)
			continue
		}
		env.Logger.Printf("INFO : Site #%d : %s\n", i, site.UserHandle)

		wp = wordpress.NewWordPress(env, &site)
		wp.SitePosts()
		if len(wp.Site.Posts) == 0 {
			env.Logger.Printf("WARN : NO SitePosts @ %s : %s\n", site.UserHandle, wp.Site.Error)
			continue
		}

		msgs = wp.PostsToMsgs()
		if len(msgs) == 0 {
			env.Logger.Printf("WARN : NO PostsToMsgs @ %s\n", site.UserHandle)
			continue
		}

		// Get access token for upsert of this user's channel
		wp.Env.Client.User = site.UserHandle
		tkn := wp.GetTkn()
		if tkn == "" {
			continue
		}

		for _, msg := range msgs {
			rsp := env.UpsertMsgByTkn(&msg)
			env.Logger.Printf("INFO : UpsertMsgByTkn @ %s : HTTP %d\n", site.UserHandle, rsp.Code)
		}

		domain := strings.Split(site.HostURL, "//")[1]
		if err := env.SetCache(domain+SUFFIX_MSGS, convert.Stringify(msgs)); err != nil {
			env.Logger.Printf("ERR : SetCache @ %s : *"+SUFFIX_MSGS+" : %s\n", site.UserHandle, err.Error())
		}
	}
}

// UpsertPostsChron repeatedly runs the UpsertPosts task once per hours, forever.
func UpsertPostsChron(env *client.Env, hours int) {
	out, err := conf.String(env)
	if err != nil {
		env.Logger.Printf("ERR : generating config for output\n")
	}
	env.Logger.Printf("INFO : Config ...\n%v\n", out)

	i := 1
	for {
		env.Logger.Printf("INFO : UpsertPosts : BEGIN #%d\n", i)

		UpsertPosts(env)

		env.Logger.Printf("INFO : UpsertPosts : END #%d\n", i)

		time.Sleep(time.Duration(hours) * time.Hour)
		i += 1
	}
}
