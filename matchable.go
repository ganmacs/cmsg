package main

import "regexp"

type Matchable interface {
	Match(re *regexp.Regexp) bool
}
