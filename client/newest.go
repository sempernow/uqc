package client

import (
	"fmt"

	"github.com/imroc/req/v3"
)

// Newest retreives list of newest messages
// OBSOLETE : TODO : update to fit newer spec
func Newest(env *Env) error { // ☩ go run ./app/htp --services-root=https://swarm.foo newest |jq .

	endpt := env.BaseAPI + "/ml/top/5/10"

	timeout := env.Timeout

	var result interface{}

	// Create and send a request with the custom client and settings
	client := req.C(). // Use C() to create a client
				SetUserAgent(UserAgent). // Chainable client settings
				SetTimeout(timeout)      //.DevMode()
	resp, err := client.R(). // Use R() to create a request
					SetHeader("Accept", "application/json").
		//SetPathParam("mid", "e7309529-e651-4108-b1fb-2ed56a30f229").
		//SetQueryParam("page", "1").
		SetResult(&result).
		Get(endpt)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
