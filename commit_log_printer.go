package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

type CommitLogPrinter struct {
	in      chan *CommitLog
	oneline bool
}

func NewCommitLogPrinter(c *cli.Context, ch *CommitLogChannel) *CommitLogPrinter {
	return &CommitLogPrinter{
		in:      ch.Out,
		oneline: !c.Bool("r"),
	}
}

func (p *CommitLogPrinter) Print() {
	for c := range p.in {
		fmt.Println(c.Format(p.oneline))
	}
}
