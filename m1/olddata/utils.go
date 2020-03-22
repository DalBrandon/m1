// --------------------------------------------------------------------
// utils.go -- Utilities for olddata package.
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

package olddata

import (
	"dbe/lib/util"
)

// findmax returns the maxium number in the argument list.
func findmax(ifields ...int) int {
	n := -1
	for i := 0; i < len(ifields); i++ {
		if ifields[i] > n {
			n = ifields[i]
		}
	}
	return n
}

// remove_commas delete any commas found in the input string.
func remove_commas(s string) string {
	n := len(s)
	sout := ""
	for i := 0; i < n; i++ {
		c := s[i]
		if c != ',' {
			sout += string(c)
		}
	}
	return sout
}

func all_fields_blank(ss []string) bool {
	for _, s := range ss {
		if !util.Blank(s) {
			return false
		}
	}
	return true
}
