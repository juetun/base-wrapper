package utils

import (
	"fmt"
	"regexp"
)

func RoutePathMath(regexpString, path string) (matched bool, err error) {
	matched, err = regexp.Match(regexpString, []byte(path))
	return
}

func RoutePathToRegexp(path string) (regexpString string, err error) {
	var mat *regexp.Regexp
	mat, err = regexp.Compile(":[^/]+")
	regexpString = fmt.Sprintf("^%s$", mat.ReplaceAllString(path, "([^/]+)"))
	return
}

