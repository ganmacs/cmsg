package main

import (
	"os/exec"
	"strings"
)

func GitConfigs(key string) []string {
	config := GitConfig(key)
	return strings.Split(config, ",")
}

func GitConfig(key string) string {
	return ExecCommand("git", "config", key)
}

func GhqRoots() string {
	return ExecCommand("ghq", "root", "--all")
}

func ExecCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	buff, _ := cmd.Output()
	if len(buff) > 0 {
		return string(buff[:len(buff)-1])
	}
	return string(buff)
}
