package client

import (
	"fmt"

	"github.com/sempernow/uqc/kit/convert"
)

func (env *Env) PrettyPrint() error {
	save := env.Client.Pass
	env.Client.Pass = fmt.Sprintf("%.3s•••", env.Client.Pass)
	_, err := fmt.Println(convert.PrettyPrint(env))
	env.Client.Pass = save
	//fmt.Println("\nenv.Args.Num(1) :", env.Args.Num(1))
	return err
}
