// --------------------------------------------------------------------
// fmtstrslice.go -- Formats a string slice
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

func FormatStrSlice(slst []string) string {
	sout := ""
	bFirst := true
	for _, s := range slst {
		if !bFirst {
			sout += ", "
		}
		bFirst = false
		sout += s
	}
	return sout
}
