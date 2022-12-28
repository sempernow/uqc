package client

import (
	"fmt"

	"github.com/sempernow/kit/convert"
)

func (env *Env) PrettyPrint() error {
	save := []string{env.Client.Pass, env.Client.Key, env.SitesPass}
	env.Client.Pass = fmt.Sprintf("%.3s•••", save[0])
	env.Client.Key = fmt.Sprintf("%.3s•••", save[1])
	env.SitesPass = fmt.Sprintf("%.3s•••", save[2])

	_, err := fmt.Println(convert.PrettyPrint(env))

	env.Client.Pass = save[0]
	env.Client.Key = save[1]
	env.SitesPass = save[2]

	return err
}
