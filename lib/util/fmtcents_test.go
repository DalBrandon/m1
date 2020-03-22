// --------------------------------------------------------------------
// fmtcents_test.go -- Test the cent formatting function
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package util

import (
	"testing"
)

type Centtest struct {
	Value  int
	Result string
}

var cent_tests []Centtest = []Centtest{
	{1, "0.01"},
	{-1, "-0.01"},
	{100, "1.00"},
	{-100, "-1.00"},
	{-212, "-2.12"},
	{-124523, "-1,245.23"},
	{0, "0.00"},
	{-1234567890123, "-12,345,678,901.23"},
	{1234567890123, "12,345,678,901.23"},
	{-4567890123, "-45,678,901.23"},
	{4567890123, "45,678,901.23"},
	{-67890199, "-678,901.99"},
	{67890199, "678,901.99"},
	{-890146, "-8,901.46"},
	{890146, "8,901.46"},
	{-90145, "-901.45"},
	{90145, "901.45"},
	{-4123, "-41.23"},
	{4123, "41.23"},
	{-101, "-1.01"},
	{101, "1.01"},
}

func Test_FmtCent(t *testing.T) {
	for _, x := range cent_tests {
		sout := CentsToStr(x.Value)
		if sout != x.Result {
			t.Fatalf("CentToStr fail. Input = %d, Output = %q, Expected = %q",
				x.Value, sout, x.Result)
		}
	}
}
