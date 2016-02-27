package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/codegangsta/cli"
)

type RepositoryChannel struct {
	SearchPatterns []string
	IgnorePatterns []string
	Out            RepoChan
}

type RepoChan chan *Repository

func NewRepositoryChannel(c *cli.Context) *RepositoryChannel {
	return &RepositoryChannel{
		SearchPatterns: append(c.StringSlice("S"), GitConfigs("cmsg.searchPaths")...),
		IgnorePatterns: append(c.StringSlice("I"), GitConfigs("cmsg.ignorePaths")...),
		Out:            make(RepoChan),
	}
}

func (ch *RepositoryChannel) Start() *RepositoryChannel {
	go func() {
		repos := ch.repos()
		for r := range repos.search(ch.SearchPatterns).ingore(ch.IgnorePatterns) {
			ch.Out <- r
		}
		close(ch.Out)
	}()
	return ch
}

func (ch *RepositoryChannel) repos() RepoChan {
	ret := make(RepoChan)
	root := GhqRoots()
	cmd := exec.Command("ghq", "list", "-p")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		cmd.Start()
		defer cmd.Wait()
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			ret <- newRepository(scanner.Text(), root)
		}
		close(ret)
	}()
	return ret
}

func (ch RepoChan) search(repos []string) RepoChan {
	for _, repo := range repos {
		if repo != "" {
			ch = ch.collect(repo, true)
		}
	}
	return ch
}

func (ch RepoChan) ingore(repos []string) RepoChan {
	for _, repo := range repos {
		if repo != "" {
			ch = ch.collect(repo, false)
		}
	}
	return ch
}

func (ch RepoChan) collect(repo string, include bool) RepoChan {
	ret := make(RepoChan)
	re, _ := regexp.Compile(repo)

	go func() {
		for repository := range ch {
			if include {
				if !repository.Match(re) {
					ret <- repository
				}
			} else {
				if repository.Match(re) {
					ret <- repository
				}
			}
		}
		close(ret)
	}()
	return ret
}
