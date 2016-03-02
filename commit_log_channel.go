package main

import (
	"regexp"
	"sync"

	"github.com/codegangsta/cli"
)

type CommitLogChannel struct {
	SearchPatterns []string
	IgnorePatterns []string
	noMerges       bool
	in             chan *Repository
	Out            chan *CommitLog
	limit          chan bool
}

func NewCommitLogChannel(c *cli.Context, input *RepositoryChannel) *CommitLogChannel {
	return &CommitLogChannel{
		SearchPatterns: append(c.StringSlice("s"), GitConfigs("cmsg.searchCommits")...),
		IgnorePatterns: append(c.StringSlice("i"), GitConfigs("cmsg.ignoreCommits")...),
		noMerges:       c.Bool("no-merges"),
		in:             input.Out,
		Out:            make(CChan, c.Int("j")),
		limit:          make(chan bool, c.Int("j")),
	}
}

func (ch *CommitLogChannel) Start() *CommitLogChannel {
	go func() {
		commits := ch.commits()
		for c := range commits.search(ch.SearchPatterns).ignore(ch.IgnorePatterns) {
			ch.Out <- c
		}
		close(ch.Out)
	}()
	return ch
}

func (ch *CommitLogChannel) commits() CChan {
	ret := make(CChan)
	wg := new(sync.WaitGroup)

	go func() {
		for repo := range ch.in {
			wg.Add(1)
			ch.limit <- true

			go func(r *Repository) {
				for clog := range CommitLogs(r, ch.noMerges) {
					ret <- clog
				}
				wg.Done()
				<-ch.limit
			}(repo)
		}
		wg.Wait()
		close(ret)
	}()
	return ret
}

type CChan chan *CommitLog

func (rc CChan) search(repos []string) CChan {
	for _, r := range repos {
		if r != "" {
			rc = rc.collect(r, true)
		}
	}
	return rc
}

func (ch CChan) ignore(repos []string) CChan {
	for _, r := range repos {
		if r != "" {
			ch = ch.collect(r, false)
		}
	}
	return ch
}

func (ch CChan) collect(r string, search bool) CChan {
	ret := make(CChan)
	re, _ := regexp.Compile(r)
	go func() {
		for c := range ch {
			if search {
				if !c.Match(re) {
					ret <- c
				}
			} else {
				if c.Match(re) {
					ret <- c
				}
			}
		}
		close(ret)
	}()
	return ret
}
