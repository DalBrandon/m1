// --------------------------------------------------------------------
// pwhash_test.go -- Test it.
//
// Created 2018-09-29 DLB
// --------------------------------------------------------------------

package pwhash

import (
	"fmt"
	"testing"
)

func Test_dohash(t *testing.T) {
	tests := []string{"I love lucy", "love", ""}

	for _, pw := range tests {
		hash, err := HashPassword(pw)
		if err != nil {
			t.Fatalf("%q failed to hash. Err=%v\n", pw, err)
		}
		ok := CheckPasswordHash(pw, hash)
		if !ok {
			t.Fatalf("%q failed loop back test.\n", pw)
		}
		fmt.Printf("%q => %q\n", pw, hash)
	}
}
