package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type CommitLog struct {
	body string
}

func (c *CommitLog) Match(re *regexp.Regexp) bool {
	return nil == re.FindStringIndex(c.body)
}

func (c *CommitLog) Format(oneline bool) string {
	if oneline {
		return toOneline(c.body)
	}
	return c.body
}

func toOneline(body string) string {
	s := strings.Replace(body, "\n", "|", -1)
	return strings.TrimSuffix(s, "|")
}

func CommitLogs(repo *Repository, noMerge bool) chan *CommitLog {
	ch := make(chan *CommitLog)
	var cmd *exec.Cmd
	if noMerge {
		cmd = exec.Command("git", "log", "--pretty=format:["+repo.Name+"] <%an> %B%x07", "--no-merges")
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:["+repo.Name+"] <%an> %B%x07")
	}

	cmd.Dir = repo.Path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		cmd.Start()
		defer cmd.Wait()
		scanner := bufio.NewScanner(stdout)
		scanner.Split(scanCommitLog)
		for scanner.Scan() {
			body := scanner.Text()
			ch <- &CommitLog{body: body}
		}
		close(ch)
	}()
	return ch
}

func scanCommitLog(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var i, j int
	if i = bytes.IndexByte(data, '\x07'); i >= 0 {
		if j = bytes.IndexByte(data[i:], '\n'); j >= 0 {
			return i + j + 1, dropSep(dropCR(data[0 : i+j])), nil
		}
	}

	if atEOF {
		return len(data), dropSep(dropCR(data)), nil
	}

	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func dropSep(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\x07' {
		return data[0 : len(data)-1]
	}
	return data
}
