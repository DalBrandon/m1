// --------------------------------------------------------------------
// checkdir.go -- Proper check for a directory.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

import (
	"os"
	"strings"
)

func DirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return false
	}
	return true
}

func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

func GetHomeDir() string {
	env := os.Environ()
	for _, e := range env {
		wrds := strings.Split(e, "=")
		if len(wrds) != 2 {
			continue
		}
		if wrds[0] == "HOME" {
			return strings.TrimSpace(wrds[1])
		}
	}
	return "."
}
