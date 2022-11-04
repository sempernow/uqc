package commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/client/wordpress"
	"github.com/sempernow/uqc/kit/convert"
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
			client.GhostPrint("\nERR @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		} else {
			client.GhostPrint("\n@ %s : HTTP %d\n", env.Client.User, rsp.Code)
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
			filepath.Join(env.Assets, "wp", "assets", "media", "avatars", (site.UserHandle + avatar)),
		); err != nil {
			avatar = "wordpress" + avatar
		} else {
			avatar = site.UserHandle + avatar
		}
		if _, err := os.ReadFile(
			filepath.Join(env.Assets, "wp", "assets", "media", "banners", (site.UserHandle + banner)),
		); err != nil {
			banner = "uqrate-banner.webp"
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
			client.GhostPrint("\nERR @ %s : HTTP %d\n", env.Client.User, rsp.Code)
		} else {
			client.GhostPrint("\n@ %s : HTTP %d\n", env.Client.User, rsp.Code)
		}
	}
}

// PurgeCacheTkns removes tkn.* cache.
func PurgeCacheTkns(env *client.Env) {
	client.GhostPrint("\nDelete Tokens (tkn.*) @ %s\n", env.Cache)
	sites := wordpress.GetSitesList(env)
	for _, site := range sites {
		fname := "tkn." + site.UserHandle
		if err := os.Remove(filepath.Join(env.Cache, "keys", fname)); err != nil {
			//client.GhostPrint("\nERR @ os.Remove : %v\n", err)
		} else {
			client.GhostPrint("\nDEL @ %s\n", fname)
		}
	}
}

// PurgeCachePosts removes *_posts.json and *_msgs.json cache.
func PurgeCachePosts(env *client.Env) {
	client.GhostPrint("\nDelete Posts and Messages Cache : *{_posts,_msgs}.json : @ %s\n", env.Cache)
	sites := wordpress.GetSitesList(env)
	for _, site := range sites {
		domain := strings.Split(site.HostURL, "//")[1]
		fname := domain + "_posts.json"
		if err := os.Remove(filepath.Join(env.Cache, fname)); err != nil {
			//client.GhostPrint("\nERR @ os.Remove : %v\n", err)
		} else {
			client.GhostPrint("\nDEL @ %s\n", fname)
		}
		fname = domain + "_msgs.json"
		if err := os.Remove(filepath.Join(env.Cache, fname)); err != nil {
			//client.GhostPrint("\nERR @ os.Remove : %v\n", err)
		} else {
			client.GhostPrint("\nDEL @ %s\n", fname)
		}
	}
}

// UpsertPosts converts []Post into []client.Message of all sites in []Site list,
// upserting the Uqrate messages to their associated channel (mirror) per site.
func UpsertPosts(env *client.Env) {
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
		client.GhostPrint("\n=== @ %d : %s\n", i, site.UserHandle)

		wp = wordpress.NewWordPress(env, &site)
		wp.SitePosts()
		if len(wp.Site.Posts) == 0 {
			client.GhostPrint("\nWARN : NO Posts @ site : %s : %s\n", site.UserHandle, wp.Site.Error)
			continue
		}

		msgs = wp.PostsToMsgs()
		if len(msgs) == 0 {
			client.GhostPrint("\nWARN : NO Messages converted for site : %s\n", site.UserHandle)
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
			client.GhostPrint("\nINFO : Messages Upserted @ %s : rsp: %s\n", site.UserHandle, convert.Stringify(rsp))
		}

		domain := strings.Split(site.HostURL, "//")[1]
		if err := env.SetCache(domain+"_msgs.json", convert.Stringify(msgs)); err != nil {
			client.GhostPrint("\nERR : setting *_msgs.json cache : %s : %s\n", site.UserHandle, err.Error())
			//continue
		}
	}
}
