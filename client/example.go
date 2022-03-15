package client

import (
	"fmt"

	"github.com/imroc/req/v3"
)

func Example(env *Env) error {
	// For test, you can create and send a request with the global default
	// client, use DevMode to see all details, try and suprise :)
	// req.DevMode()
	// req.Get("https://api.github.com/users/imroc")

	var result interface{}

	// Create and send a request with the custom client and settings
	client := req.C(). // Use C() to create a client
				SetUserAgent(env.UserAgent). // Chainable client settings
				SetTimeout(env.Timeout).
				DevMode()
	resp, err := client.R(). // Use R() to create a request
					SetHeader("Accept", "application/vnd.github.v3+json"). // Chainable request settings
					SetPathParam("username", "imroc").                     // k-v pairs; parameterize the URL (below)
					SetQueryParam("page", "1").
					SetResult(&result).
					Get("https://api.github.com/users/{username}/repos")
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
