package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	tuf "github.com/DataDog/go-tuf/client"
	"github.com/dustin/go-humanize"
	"github.com/flynn/go-docopt"
)

func init() {
	register("list", cmdList, `
usage: tuf-client list [-s|--store=<path>] <url>

Options:
  -s <path>    The path to the local file store [default: tuf.db]

List available target files.
  `)
}

func cmdList(args *docopt.Args, client *tuf.Client) error {
	if _, err := client.Update(); err != nil {
		return err
	}
	targets, err := client.Targets()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "PATH\tSIZE")
	for path, meta := range targets {
		fmt.Fprintf(w, "%s\t%s\n", path, humanize.Bytes(uint64(meta.Length)))
	}
	return nil
}
