// --------------------------------------------------------------------
// strings_test.go -- Test the string functions
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package util

import (
	"fmt"
	"testing"
)

func Test_CleanStr(t *testing.T) {
	x := "abc\ndef\tghi"
	y := "."
	z := CleanStr(x, y)
	if z != "abc.def.ghi" {
		t.Fatalf("CleanStr(%q,%q) fails with output = %q", x, y, z)
	}
}

func Test_Chop(t *testing.T) {
	x := "123456789"
	y := FixStrLen(x, 9, "")
	if y != x {
		t.Fatalf("FixStrLen fails (t,5,\"\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 5, "")
	if y != "12345" {
		t.Fatalf("FixStrLen fails (t,5,\"\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 15, "")
	if y != "123456789      " {
		t.Fatalf("FixStrLen fails (t,15,\"\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 15, "...")
	if y != "123456789      " {
		t.Fatalf("FixStrLen fails (t,15,\"...\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 5, "...")
	if y != "12..." {
		t.Fatalf("FixStrLen fails (t,5,\"...\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 3, "...")
	if y != "..." {
		t.Fatalf("FixStrLen fails (t,3,\"...\"): Input=%q, Output=%q", x, y)
	}
	y = FixStrLen(x, 2, "...")
	if y != ".." {
		t.Fatalf("FixStrLen fails (t,2,\"...\"): Input=%q, Output=%q", x, y)
	}
	x = ""
	y = FixStrLen(x, 2, "...")
	if y != "  " {
		t.Fatalf("FixStrLen fails (t,2,\"...\"): Input=%q, Output=%q", x, y)
	}
}

func Test_Blank(t *testing.T) {
	if !Blank(" \n\n ") {
		t.Fatalf("Blank() returned false for spaces and newlines.")
	}
	if !Blank(nil) {
		t.Fatalf("Blank() returned false for nil.")
	}
	if Blank("  Hi There  ") {
		t.Fatalf("Blank() returned true for non blank string.")
	}
	sptr := new(string)
	*sptr = "  "
	if !Blank(sptr) {
		t.Fatalf("Blank() returned false for pointer to blank string.")
	}
	*sptr = "  Stuff "
	if Blank(sptr) {
		t.Fatalf("Blank() returned true for pointer to non-blank string.")
	}
	num := 123
	if Blank(num) {
		t.Fatalf("Blank() returned true for non-string argument.")
	}
}

func Test_ContainsOnly(t *testing.T) {

	slist := "abcdefghijklmnopqrstuvwxyz"
	t1 := "MoonLight"
	t2 := "starlight"
	t3 := "!bang"
	if ContainsOnly(t1, slist) != false {
		t.Fatalf("Contains Only Failed. Did not report false for input (%q, %q)", t1, slist)
	}
	if ContainsOnly(t2, slist) != true {
		t.Fatalf("Contains Only Failed. Did not report true for input (%q, %q)", t2, slist)
	}
	if ContainsOnly(t3, slist) != false {
		t.Fatalf("Contains Only Failed. Did not report true for input (%q, %q)", t3, slist)
	}
}

func Test_ContainsRune(t *testing.T) {

	slist := "ABCML"
	t1 := "MoonLight"
	t2 := "starlight"
	t3 := "!bang"
	if ContainsRune(t1, slist) != true {
		t.Fatalf("ContainsRune Failed. Did not report true for input (%q, %q)", t1, slist)
	}
	if ContainsRune(t2, slist) != false {
		t.Fatalf("ContainsRune Failed. Did not report false for input (%q, %q)", t2, slist)
	}
	if ContainsRune(t3, slist) != false {
		t.Fatalf("ContainsRune Failed. Did not report false for input (%q, %q)", t3, slist)
	}
}

func Test_ForceRuneSet(t *testing.T) {

	slist1 := "abcdefghijklmnopqrstuvwxyz"
	tin1 := "This is cool"
	tck1 := "_his_is_cool"
	tout1 := ForceRuneSet(tin1, slist1, '_')
	if tout1 != tck1 {
		t.Fatalf("FourceRuneSet failed. Reported %q for input %q.", tout1, tin1)
	}

	slist2 := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz "
	tin2 := "This is cool! But, not right? ##Sorry##"
	tck2 := "This is cool* But* not right* **Sorry**"
	tout2 := ForceRuneSet(tin2, slist2, '*')
	if tout2 != tck2 {
		t.Fatalf("FourceRuneSet failed. Reported %q for input %q.", tout2, tin2)
	}
}

func Test_SplitComma(t *testing.T) {
	sline := " a ,   b, c "
	lst := SplitComma(sline)
	if len(lst) != 3 {
		t.Fatalf("SplitComma failed to return 3 strings. Found = %d", len(lst))
	}
	testlst := []string{"a", "b", "c"}
	if !SameStringSlice(lst, testlst) {
		t.Fatalf("SplitComma failed.\nExpected=%v\nFound=%v\n", testlst, lst)
	}
	sline = ""
	lst = SplitComma(sline)
	if len(lst) != 0 {
		t.Fatalf("SplitComma failed. Instead of empty slice, found = %v", lst)
	}
}

func Test_InStringSlice(t *testing.T) {
	lst := []string{"a", "b", "c"}
	if InStringSlice(lst, "e") || !InStringSlice(lst, "a") {
		t.Fatalf("InStringSlice failed.")
	}
}

func Test_RemoveStringFromSlice(t *testing.T) {
	lst := []string{"a", "b", "c"}
	lout := RemoveStringFromSlice(lst, "a")
	if !SameStringSlice([]string{"b", "c"}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 1: %v", lout)
	}
	lout = RemoveStringFromSlice(lst, "b")
	if !SameStringSlice([]string{"a", "c"}, lout) {
		fmt.Printf("len(lout)=%d\n", len(lout))
		t.Fatalf("RemoveStringFromSlice failed test 2: %v", lout)
	}
	lout = RemoveStringFromSlice(lst, "c")
	if !SameStringSlice([]string{"a", "b"}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 3: %v", lout)
	}
	lout = RemoveStringFromSlice(lst, "d")
	if !SameStringSlice([]string{"a", "b", "c"}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 3: %v", lout)
	}
	lst = []string{"a", "b"}
	lout = RemoveStringFromSlice(lst, "a")
	if !SameStringSlice([]string{"b"}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 4: %v", lout)
	}
	lout = RemoveStringFromSlice(lst, "b")
	if !SameStringSlice([]string{"a"}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 5: %v", lout)
	}
	lst = []string{"a"}
	lout = RemoveStringFromSlice(lst, "a")
	if !SameStringSlice([]string{}, lout) {
		t.Fatalf("RemoveStringFromSlice failed test 6: %v", lout)
	}
	lout = RemoveStringFromSlice([]string{}, "b")
	if len(lout) != 0 {
		t.Fatalf("RemoveStringFromSlice failed with empty input.")
	}
	lout = RemoveStringFromSlice(nil, "b")
	if len(lout) != 0 {
		t.Fatalf("RemoveStringFromSlice failed with nil input.")
	}
	lst = []string{"a", "b", "a", "c", "a"}
	lout = RemoveStringFromSlice(lst, "a")
	if !SameStringSlice([]string{"b", "c"}, lout) {
		t.Fatalf("RemoveStringFromSlice fail test 7: %v", lout)
	}
	lst = []string{"a", "b", "", "c"}
	lout = RemoveStringFromSlice(lst, "")
	if !SameStringSlice([]string{"a", "b", "c"}, lout) {
		t.Fatalf("RemoveStringFromSlice fail test 8: %v", lout)
	}

}

func Test_SameSliceStringContents(t *testing.T) {
	lst0 := []string{}
	lst1 := []string{"a", "b", "c"}
	lst2 := []string{"b", "a", "c"}
	lst3 := []string{"a", "c", "e"}
	lst4 := []string{"a", "b", "c", "d"}
	var lst5 []string
	if SameSliceStringContents(lst5, lst1) {
		t.Fatalf("SameSliceStringContents failed with one nil input.")
	}
	if !SameSliceStringContents(lst5, lst0) {
		t.Fatalf("SameSliceStringContents failed with nil and empty slice.")
	}
	if SameSliceStringContents(lst1, lst0) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst0)
	}
	if !SameSliceStringContents(lst1, lst2) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst2)
	}
	if SameSliceStringContents(lst1, lst3) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst3)
	}
	if SameSliceStringContents(lst1, lst4) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst3)
	}
	if SameSliceStringContents(lst1, lst0) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst0)
	}
	if !SameSliceStringContents(lst1, lst1) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst1, lst1)
	}
	lst6 := []string{"a", "a", "b", "b"}
	lst7 := []string{"a", "b", "a", "b"}
	lst8 := []string{"c", "b", "a", "a"}
	if !SameSliceStringContents(lst6, lst7) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst6, lst7)
	}
	if SameSliceStringContents(lst6, lst8) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst6, lst7)
	}
	if SameSliceStringContents(lst7, lst8) {
		t.Fatalf("SameSliceStringContents failed with %v and %v", lst6, lst7)
	}
}

func Test_CommonPrefix(t *testing.T) {
	lst0 := []string{}
	lst1 := []string{"aabceasdf", "aafrz", "aaqfy"}
	lst2 := []string{"abc", "def", "ghi"}
	lst3 := []string{"/home/jim/bob/file1", "/home/jim/sue/file2", "/home/jim/kim/file3"}
	lst4 := []string{"", "aa", "bb", "cc", "dd"}
	lst5 := []string{"aa", "bb", "", "dd"}
	if s := CommonPrefix(lst0); s != "" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst0)
	}
	if s := CommonPrefix(lst1); s != "aa" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst1)
	}
	if s := CommonPrefix(lst2); s != "" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst2)
	}
	if s := CommonPrefix(lst3); s != "/home/jim/" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst3)
	}
	if s := CommonPrefix(lst4); s != "" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst4)
	}
	if s := CommonPrefix(lst5); s != "" {
		t.Fatalf("CommonPrefix failed. Got %q for %v.", s, lst5)
	}
}

func Test_ParseArgs(t *testing.T) {
	type pair struct {
		Num     int
		TestStr string
		Ans     []string
	}
	tstlst := make([]pair, 0, 10)
	tstlst = append(tstlst, pair{1, "", []string{}})
	tstlst = append(tstlst, pair{2, "one", []string{"one"}})
	tstlst = append(tstlst, pair{3, " one     two  ", []string{"one", "two"}})
	tstlst = append(tstlst, pair{4, "one\ntwo\tthree\n", []string{"one", "two", "three"}})
	tstlst = append(tstlst, pair{5, "a one b", []string{"a", "one", "b"}})
	tstlst = append(tstlst, pair{6, `"one two" three`, []string{"\"one two\"", "three"}})
	tstlst = append(tstlst, pair{7, `abc" two one "def g`, []string{"abc\" two one \"def", "g"}})
	tstlst = append(tstlst, pair{8, `"`, []string{"\""}})
	tstlst = append(tstlst, pair{9, `""`, []string{"\"\""}})
	tstlst = append(tstlst, pair{10, `" "`, []string{"\" \""}})
	tstlst = append(tstlst, pair{11, `ab cd "ef `, []string{"ab", "cd", "ef "}})
	tstlst = append(tstlst, pair{12, `ab cd "e f" gh`, []string{"ab", "cd", "e f", "gh"}})
	tstlst = append(tstlst, pair{13, `cmd p1="one two" p2=three`, []string{"cmd", "p1=\"one two\"", "p2=three"}})

	for _, tt := range tstlst {
		ans := SplitArgs(tt.TestStr)
		if len(ans) != len(tt.Ans) {
			t.Fatalf("ParseArgs failed on test %d. Input=%q, Out=%v.", tt.Num, tt.TestStr, ans)
		}
	}
}

// func Test_SameStringSliceContents(t *testing.T) {
//  if !SameSliceStringContents(nil, nil) {
//      t.Fatalf("SameSliceStringContents failed on test 1.")
//  }
//  if !SameSliceStringContents([]string{}, nil) {
//      t.Fatalf("SameSliceStringContents failed on test 2.")
//  }
//  if !SameSliceStringContents([]string{}, []string{}) {
//      t.Fatalf("SameSliceStringContents failed on test 3.")
//  }
//  if SameSliceStringContents([]string{"a"}, []string{""}) {
//      t.Fatalf("SameSliceStringContents failed on test 4.")
//  }
//  if SameSliceStringContents([]string{"a"}, []string{"a", "a"}) {
//      t.Fatalf("SameSliceStringContents failed on test 5.")
//  }
//  if SameSliceStringContents([]string{"a", "b"}, []string{"b", "c"}) {
//      t.Fatalf("SameSliceStringContents failed on test 6.")
//  }
//  if !SameSliceStringContents([]string{"a", "b", "c"}, []string{"a", "b", "c"}) {
//      t.Fatalf("SameSliceStringContents failed on test 7.")
//  }
//  if !SameSliceStringContents([]string{"a", "b"}, []string{"b", "a"}) {
//      t.Fatalf("SameSliceStringContents failed on test 8.")
//  }
//  if SameSliceStringContents([]string{"", "b"}, []string{"c", ""}) {
//      t.Fatalf("SameSliceStringContents failed on test 9.")
//  }
//  if SameSliceStringContents([]string{"a", "b", "a"}, []string{"a", "b", "b"}) {
//      t.Fatalf("SameSliceStringContents failed on test 10.")
//  }
//  if !SameSliceStringContents([]string{"a", "b", "a"}, []string{"b", "a", "a"}) {
//      t.Fatalf("SameSliceStringContents failed on test 11.")
//  }
// }

func Test_SameStringSliceContents2(t *testing.T) {
	if !SameSliceStringContents(nil, nil) {
		t.Fatalf("SameSliceStringContents failed on test 1.")
	}
	if !SameSliceStringContents([]string{}, nil) {
		t.Fatalf("SameSliceStringContents failed on test 2.")
	}
	if !SameSliceStringContents([]string{}, []string{}) {
		t.Fatalf("SameSliceStringContents failed on test 3.")
	}
	if SameSliceStringContents([]string{"a"}, []string{""}) {
		t.Fatalf("SameSliceStringContents failed on test 4.")
	}
	if SameSliceStringContents([]string{"a"}, []string{"a", "a"}) {
		t.Fatalf("SameSliceStringContents failed on test 5.")
	}
	if SameSliceStringContents([]string{"a", "b"}, []string{"b", "c"}) {
		t.Fatalf("SameSliceStringContents failed on test 6.")
	}
	if !SameSliceStringContents([]string{"a", "b", "c"}, []string{"a", "b", "c"}) {
		t.Fatalf("SameSliceStringContents failed on test 7.")
	}
	if !SameSliceStringContents([]string{"a", "b"}, []string{"b", "a"}) {
		t.Fatalf("SameSliceStringContents failed on test 8.")
	}
	if SameSliceStringContents([]string{"", "b"}, []string{"c", ""}) {
		t.Fatalf("SameSliceStringContents failed on test 9.")
	}
	if SameSliceStringContents([]string{"a", "b", "a"}, []string{"a", "b", "b"}) {
		t.Fatalf("SameSliceStringContents failed on test 10.")
	}
	if !SameSliceStringContents([]string{"a", "b", "a"}, []string{"b", "a", "a"}) {
		t.Fatalf("SameSliceStringContents failed on test 11.")
	}
}

func Test_StrLeft(t *testing.T) {
	if StrLeft("abc", 2) != "abc" {
		t.Fatalf("StrLeft failed on test 1.")
	}
	if StrLeft("abc", 3) != "abc" {
		t.Fatalf("StrLeft failed on test 2.")
	}
	if StrLeft("", 0) != "" {
		t.Fatalf("StrLeft failed on test 3.")
	}
	if StrLeft("", 2) != "  " {
		t.Fatalf("StrLeft failed on test 4.")
	}
	if StrLeft("abc", 4) != " abc" {
		t.Fatalf("StrLeft failed on test 5.")
	}
	if StrLeft("abc", 6) != "   abc" {
		t.Fatalf("StrLeft failed on test 6.")
	}
}
