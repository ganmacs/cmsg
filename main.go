package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const version string = "0.0.1"

var flags = []cli.Flag{
	cli.IntFlag{
		Name:  "j",
		Value: 50,
		Usage: "Specifies the number of parallel.",
	},
	cli.BoolFlag{
		Name:  "r, raw",
		Usage: "Format each commit message as raw format",
	},
	cli.StringSliceFlag{
		Name:  "s",
		Usage: "search messages.",
	},
	cli.StringSliceFlag{
		Name:  "i",
		Usage: "ignore commit messages.",
	},
	cli.StringSliceFlag{
		Name:  "S",
		Usage: "search repositories.",
	},
	cli.StringSliceFlag{
		Name:  "I",
		Usage: "ignore repositories.",
	},
	},
}

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "cmsg"
	app.Usage = "Commit message"
	app.Version = version
	app.Author = "ganmacs"
	app.Email = "ganmacs@gmail.com"
	app.Action = action
	app.Flags = flags

	return app
}

func action(c *cli.Context) {
	rch := NewRepositoryChannel(c).Start()
	cch := NewCommitLogChannel(c, rch).Start()
	NewCommitLogPrinter(c, cch).Print()
}
