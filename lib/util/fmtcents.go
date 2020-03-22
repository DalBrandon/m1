// --------------------------------------------------------------------
// fmtcents.go -- Format Cents.
//
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package util

import (
	"fmt"
)

// CentsToStr will convert a value of cents (as an integer) into an
// formatted string for dollars, including commas and negative sign.
// The output string will not contain any blank space.  The input is
// good up to  plus or minus 999,999,999,999.99 dollars.
func CentsToStr(cents int) string {
	neg := ""
	if cents < 0 {
		neg = "-"
		cents = -cents
	}
	dollars := cents / 100
	cents = cents - (dollars * 100)
	bils := dollars / 1000000000
	dollars = dollars - (bils * 1000000000)
	mils := dollars / 1000000
	dollars = dollars - (mils * 1000000)
	thous := dollars / 1000
	dollars = dollars - (thous * 1000)
	if bils > 0 {
		return fmt.Sprintf("%s%d,%03d,%03d,%03d.%02d", neg, bils, mils, thous, dollars, cents)
	}
	if mils > 0 {
		return fmt.Sprintf("%s%d,%03d,%03d.%02d", neg, mils, thous, dollars, cents)
	}
	if thous > 0 {
		return fmt.Sprintf("%s%d,%03d.%02d", neg, thous, dollars, cents)
	}
	if dollars > 0 {
		return fmt.Sprintf("%s%d.%02d", neg, dollars, cents)
	}
	return fmt.Sprintf("%s0.%02d", neg, cents)
}
