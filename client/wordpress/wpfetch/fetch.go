package wpfetch

import (
	"errors"

	"github.com/sempernow/uqc/client"
	"github.com/sempernow/uqc/kit/types"
)

const (
	Posts      = "/wp-json/wp/v2/posts"
	Tags       = "/wp-json/wp/v2/tags"
	Categories = "/wp-json/wp/v2/categories"
	Authors    = "/wp-json/wp/v2/users"
)

func Author(id int) string {
	// Retrieve from /wp-json/wp/v2/authors/:id
	return "author name"
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

func fetchWP(env *client.Env, endpt string) (string, error) {
	rsp := env.Dump(endpt, client.JSON)
	if rsp.Error != "" {
		return "", errors.New(rsp.Error)
	}
	return rsp.Body, nil
}
