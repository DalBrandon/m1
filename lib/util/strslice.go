// --------------------------------------------------------------------
// strslice.go -- functions that deal with string slices
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

import (
	"strings"
)

// CloneStringSlice makes a separate copy of a slice of strings.  If
// the input is nil, a non-nil, but empty slice is returned.
func CloneStringSlice(src []string) (dst []string) {
	if src == nil {
		return make([]string, 0)
	}
	dst = make([]string, len(src))
	if src == nil {
		return dst
	}
	for i := range src {
		dst[i] = src[i]
	}
	return dst
}

// SameStringSlice compares two slices of strings and returns true if they
// are exacly the same -- same len, same order of contents.
func SameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SplitComma returns a slice of strings from input split by commas.
// Each string in the slice has leading and trailing whitespace removed.
// Empty string returns empty slice.  For exampe: " a ,   b, c "
// returns ["a","b","c"].
func SplitComma(line string) []string {
	if Blank(line) {
		return []string{}
	}
	parts := strings.Split(line, ",")
	lst := make([]string, len(parts))
	for i, s := range parts {
		lst[i] = strings.TrimSpace(s)
	}
	return lst
}

// InStringSlice checks to see if item is in the list.
func InStringSlice(lst []string, item string) bool {
	for _, s := range lst {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveStringFromSlice finds the all occurances of the item
// in the slice and returns a new slice with the item removed.
// Note: The input slice is not harmed.
func RemoveStringFromSlice(lst []string, item string) []string {
	lout := CloneStringSlice(lst)
	for {
		// Find the item to remmove
		indx := -1
		for i, s := range lout {
			if s == item {
				indx = i
				break
			}
		}
		if indx < 0 {
			return lout
		}
		lout = append(lout[:indx], lout[indx+1:]...)
	}
}

// SameSliceStringContents returns true if the slices have
// the same elements, even if the elements are not in the
// same order.
func SameSliceStringContents(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int, len(a))
	for _, x := range a {
		m[x] = m[x] + 1
	}
	for _, x := range b {
		m[x] = m[x] - 1
	}
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}

// CommonPrefix returns the longest string that matches all the
// suffixes in the input slice.
func CommonPrefix(lst []string) string {
	suffix := ""
	if len(lst) <= 0 {
		return suffix
	}
	guess := lst[0]
	for _, r := range guess {
		test := suffix + string(r)
		for _, s := range lst {
			if !strings.HasPrefix(s, test) {
				return suffix
			}
		}
		suffix = test
	}
	return suffix
}

// SplitArgs splits the input line into argument and respects
// double quotes.  For example, "p1 p2=\"hi there\" p3" is
// returned as ["p1", "p2=\"hi there\"", "p3"]
func SplitArgs(line string) []string {
	args := make([]string, 0, 30)
	inquote := false
	currentarg := ""
	for _, c := range strings.Split(line, "") {
		if inquote {
			if c == "\"" {
				inquote = false
			}
			currentarg += c
		} else {
			if c == "\"" {
				inquote = true
				currentarg += c
			} else if c == " " || c == "\t" || c == "\n" || c == "\r" {
				if currentarg != "" {
					args = append(args, currentarg)
				}
				currentarg = ""
			} else {
				currentarg += c
			}
		}
	}
	if currentarg != "" {
		args = append(args, currentarg)
	}
	return args
}
