// --------------------------------------------------------------------
// generic_time.go -- Generic Time format parser.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package util

import (
	"fmt"
	"time"
)

var GenericTimeFormats []string = []string{
	"2006", "06", "Jan 06", "Jan-06", "Jan 2006", "January 2006", "January 06",
	"2006-01-02", "06-01-02", "1/2/06", "1/2/2006",
	"Jan 2, 06", "January 2, 06", "Jan 2, 2006", "January 2, 2006",
	"2006-01-02 15:04:05", "2006-01-02T15:04:05",
	"2006-1-2 15:4:5", "2006-1-2T15:4:5",
	"1/2/06 15:04:05", "1/2/2006 15:04:05",
	"1/2/06 15:4:5", "1/2/2006 15:4:5"}

func ParseGenericTime(s string) (time.Time, error) {
	for _, ft := range GenericTimeFormats {
		t, err := time.Parse(ft, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("Unable to convert %q to a time value.", s)
}
