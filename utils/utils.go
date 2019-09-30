package utils

import (
	"strings"

	"github.com/methlock/code-stats/config"
)

// IsCommentLine checks if line is comment line.
func IsCommentLine(row string) (result bool) {
	for _, comment := range config.LineIdentificators {
		if strings.HasPrefix(row, comment) {
			return true
		}
	}
	return false
}

// JoinPaths correctly joins root and some object in root path.
// No missing slashes and stuff like that.
func JoinPaths(root string, path string) (fullPath string) {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}
	if strings.HasPrefix(path, "/") {
		path = path[0 : len(path)-1]
	}
	return root + path
}

// IsThisIn returns true if something is in something else...
// TODO: with interface?
func IsThisIn(what string, in []string) bool {
	for _, that := range in {
		if what == that {
			return true
		}
	}
	return false
}
