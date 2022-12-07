package main

import (
	"fmt"

	"github.com/DataDog/go-tuf"
	"github.com/flynn/go-docopt"
)

func init() {
	register("get-threshold", cmdGetThreshold, `
usage: tuf get-threshold <role>

Gets the threshold for a role.  
`)
}

func cmdGetThreshold(args *docopt.Args, repo *tuf.Repo) error {
	role := args.String["<role>"]

	threshold, err := repo.GetThreshold(role)
	if err != nil {
		return err
	}

	fmt.Println("The threshold for", role, "role is", threshold)
	return nil
}
