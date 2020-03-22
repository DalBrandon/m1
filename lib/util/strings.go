// --------------------------------------------------------------------
// strings.go -- string utilities.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

import (
	"fmt"
	"strings"
)

// Returns true if the input is a blank string, a pointer to a blank string, or
// nil. A blank string is one that contains only whitespace.  Input object that
// are not strings are considered "not blank".
func Blank(ss interface{}) bool {
	if ss == nil {
		return true
	}
	s, ok := ss.(string)
	if !ok {
		sptr, ok := ss.(*string)
		if !ok {
			return false
		}
		s = *sptr
	}
	s = strings.TrimSpace(s)
	if len(s) <= 0 {
		return true
	}
	return false
}

// ContainsOnly checks to see of all the runes in s are also
// in slist.
func ContainsOnly(s, slist string) bool {
	for _, c := range s {
		i := strings.IndexRune(slist, c)
		if i < 0 {
			return false
		}
	}
	return true
}

// ContainsRune checks to see if at least one of the runes
// in slist appears in s.
func ContainsRune(s, slist string) bool {
	for _, c := range s {
		i := strings.IndexRune(slist, c)
		if i >= 0 {
			return true
		}
	}
	return false
}

// ForceRuneSet checks to make sure that all runes in
// s are contained in slist. If they are not, they
// are replaced with r.
func ForceRuneSet(s, slist string, r rune) string {
	result := ""
	for _, c := range s {
		i := strings.IndexRune(slist, c)
		if i >= 0 {
			result += string(c)
		} else {
			result += string(r)
		}
	}
	return result
}

// SelStr selects a string based on the condition flag.
// sTrue is returned if condition is true, otherwise sFalse is returned.
func SelStr(sTrue, sFalse string, condition bool) string {
	if condition {
		return sTrue
	} else {
		return sFalse
	}
}

// FixStrLen either adds spaces or chops the input so that the
// output is exactly n characters long.  If the string is chopped,
// a suffix can be provided to indicate the chopping operation.
// NOTE: currently only works for ascii strings.
func FixStrLen(s string, n int, suffix string) string {
	ns := len(suffix)
	ngot := len(s)
	if ngot == n {
		return s
	}
	if ngot < n {
		for i := 0; i < n-ngot; i++ {
			s += " "
		}
		return s
	}
	if ns > n {
		return suffix[:n]
	}
	n = n - ns
	s = s[:n] + suffix
	return s
}

// StrLeft adds spaces to the left side of the string utill
// the string length is n chars long.  If the string is already
// n chars long, or greater, nothing is done.
func StrLeft(s string, n int) string {
	ns := len(s)
	if ns >= n {
		return s
	}
	nadd := n - ns
	ss := ""
	for i := 0; i < nadd; i++ {
		ss += " "
	}
	return ss + s
}

// CleanStr replaces the unprintable ascii chars (such as newline)
// with the replacement string.
func CleanStr(s, replace string) string {
	n := len(s)
	sout := ""
	for i := 0; i < n; i++ {
		c := s[i]
		if c < 32 || c >= 127 {
			sout += replace
		} else {
			sout += string(c)
		}
	}
	return sout
}

// CleanForWeb replaces dangerous chars with *.  For right now,
// this is only a double quote and unprintable chars.
func CleanForWeb(s string) string {
	n := len(s)
	sout := ""
	replace := "*"
	for i := 0; i < n; i++ {
		c := s[i]
		if c < 32 || c >= 127 {
			sout += replace
		} else if c == 34 {
			sout += replace
		} else {
			sout += string(c)
		}
	}
	return sout
}

// CleanForCafe replaces removes non 7-bit ascii characters
// with the exception that a 65533 unicode character is
// replaced with an 'e'.  This fixes the common problem with
// data from excell files where excell converted "cafe" to
// "cafe'".  Note that a 65533 is the unicode "replacement"
// character that is used to indicate an unknown or invalid
// character.
func CleanForCafe(s string) string {
	sout := ""
	for _, y := range s {
		if y < 0 {
			continue
		}
		if y < 127 {
			sout = sout + string(y)
		}
		if y == 65533 {
			sout = sout + "e" // Fixes the 'cafe' problem introduced by excell
		}
	}
	return sout
}

// StrToBool converts a string into a bool, trying to interpert
// it as loosely as possible.  If a reasonable value for a bool
// is not found, the default and an error is returned.
func StrToBool(s string, defaultval bool) (val bool, err error) {
	s = strings.ToLower(s)
	if s == "true" || s == "t" || s == "yes" || s == "y" {
		return true, nil
	}
	if s == "false" || s == "f" || s == "no" || s == "n" {
		return false, nil
	}
	return defaultval, fmt.Errorf("Unrecognizable bool string value.")
}
