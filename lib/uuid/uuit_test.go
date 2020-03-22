// --------------------------------------------------------------------
// uuid_test.go -- test
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package uuid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	//            12345678901234567890123456789012
	su1 string = "1234567890abcdeEDCBA9876543210AB"
	su2 string = "1234567890abcdeEDCBA987654321012"
	suz string = "00000000000000000000000000000000"
)

func Test_UUID(t *testing.T) {

	t.Log("Testing the UUID methods.")
	// Make sure two new uuids are not the same.
	u1 := New()
	u2 := New()

	if u1 == u2 {
		t.Fatalf("New uuids equal! %v = %v", u1, u2)
	}
	t.Log("UUID creates different objects -- okay.")

	scheck := fmt.Sprintf("%s", u1)
	if scheck != u1.String() {
		t.Fatalf("Stringer interface not working.  String() != %%s")
	}
	t.Log("UUID stringer interfacw works.")

	if u1.String() == u2.String() {
		t.Fatalf("Strings of new uuids equal! %s = %s", u1.String(), u2.String())
	}
	if len(u1.String()) != 32 {
		t.Fatalf("String output is not 32 chars.")
	}
	t.Log("String output of UUIDs -- okay.")

	u3, err := FromString(su1)
	if err != nil {
		t.Fatalf("Unable to use FromString(). %v", err)
	}
	u4, err := FromString(u3.String())
	if u3 != u4 {
		t.Fatalf("Same uuids found unequal.")
	}
	t.Log("UUIDs can be created from strings -- okay.")

	var u5, u6 UUID
	u7 := new(UUID)
	if u5.IsZero() != true || u6.IsZero() != true || u7.IsZero() != true {
		t.Fatalf("Zero tests failed for uninitized uuids.")
	}
	if u5 != u6 || u5 != *u7 {
		t.Fatalf("Unable to compare uninitialized zero uuids.")
	}
	t.Log("Uninitalized uuids are zero -- okay.")

	u8, err := FromString(suz)
	if err != nil {
		t.Fatalf("Unable to create zero uuid using FromString.")
	}
	if u8.IsZero() == false {
		t.Fatalf("Zero created uuid found not to be zero.")
	}
	if u5 != u8 {
		t.Fatalf("Unitialized uuid not equal to zero uuid.")
	}
	t.Log("Uninitialized uuids and zero created uuids are same -- okay.")

	_, err = FromString("")
	if err == nil {
		t.Fatalf("No error produced from a blank string in FromString")
	}

	u9, err := FromString0("")
	if err != nil {
		t.Fatalf("Error produced by FromString0 on a blank string. %v", err)
	}
	if !u9.IsZero() {
		t.Fatalf("Output from FromString(\"\") not zero.")
	}

	// One more check of equivalance!
	var u10 UUID
	if u9 != u10 || u9.String() != u10.String() {
		t.Fatalf("Equal check fails for two different zero uuids.")
	}
	u11, err := FromString(suz)
	if err != nil {
		t.Fatalf("FromString fails on zero input.")
	}
	if u9 != u10 || u9 != u11 || u11.String() != u9.String() || u10.String() != u11.String() {
		t.Fatalf("Some combination of equal zero is wrong!")
	}

	u12 := New()
	u13 := ForceStr(u12.String())
	if u12 != u13 || u12.String() != u13.String() || u12 == u9 || u13 == u10 || u9.String() == u13.String() {
		t.Fatalf("Some combination of zero and not zero does not compare correctly.")
	}

	u14 := ForceStr(strings.ToLower(su1))
	u15 := ForceStr(strings.ToUpper(su1))

	if u14 != u15 || u14.String() != u15.String() {
		t.Fatalf("The same hex does not equal due to case differeneces.\n%q\n%q", u14.ustr, u15.ustr)
	}

	var u16 UUID
	if u16.String() != suz {
		t.Fatalf("String() on unitialized uuid not correct. %q", u16.String())
	}

}

type dummy struct {
	U1     UUID
	U2     UUID
	UArray []UUID
}

func dummyEqual(d1, d2 dummy) bool {
	if d1.U1 != d2.U1 || d1.U2 != d2.U2 {
		return false
	}
	if len(d1.UArray) != len(d2.UArray) {
		return false
	}
	for i, u := range d1.UArray {
		if u != d2.UArray[i] {
			return false
		}
	}
	return true
}

func Test_json(t *testing.T) {

	t.Log("Begining JSON tests...")
	u1, err1 := FromString(su1)
	u2, err2 := FromString(su2)
	if err1 != nil || err2 != nil {
		t.Fatalf("Unable to create uuids. %v %v", err1, err2)
	}
	d1 := dummy{U1: u1, U2: u2, UArray: []UUID{u1, u2}}
	buf1 := new(bytes.Buffer)
	enc := json.NewEncoder(buf1)
	enc.Encode(d1)

	var d2 dummy
	buf1out := bytes.NewReader(buf1.Bytes())
	dec := json.NewDecoder(buf1out)
	err := dec.Decode(&d2)
	if err != nil {
		t.Fatalf("Unable to decode encoded json. %v", err)
	}
	if !dummyEqual(d1, d2) {
		t.Fatalf("Round trip error with json. %v != %v", d1, d2)
	}

	var d3 dummy
	j1 := bytes.NewReader([]byte(jtest1))
	dec2 := json.NewDecoder(j1)
	err = dec2.Decode(&d3)
	if err != nil {
		t.Fatalf("Unable to decode json. %v", err)
	}

	if d3.U1.String() != strings.ToUpper(su1) {
		t.Fatalf("Bad json decode. U1 incorrect. U1=%s != %s", d3.U1.String(), su1)
	}
	if d3.U2.String() != strings.ToUpper(su2) {
		t.Fatalf("Bad json decode. U2 incorrect. U2=%s != %s", d3.U2.String(), su2)
	}
	if len(d3.UArray) != 2 {
		t.Fatalf("Bad json decode. D3 Array len != 2")
	}
	if d3.UArray[0].String() != strings.ToUpper(su1) {
		t.Fatalf("Bad json decode. UArray[0] incorrect. UArray[0]=%s != %s", d3.UArray[0].String(), su1)
	}
	if d3.UArray[1].String() != strings.ToUpper(su2) {
		t.Fatalf("Bad json decode. UArray[1] incorrect. UArray[1]=%s != %s", d3.UArray[0].String(), su2)
	}

	var d4 dummy
	j2 := bytes.NewReader([]byte(jtest2))
	dec3 := json.NewDecoder(j2)
	err = dec3.Decode(&d4)
	if err != nil {
		t.Fatalf("Unable to decode json. %v", err)
	}

	if d4.U1.IsZero() != true {
		t.Fatalf("Bad zero on json decode. U1 should be zero. Found=%s from %q", d4.U1.String(), d4.U1.ustr)
	}
	if d4.U2.IsZero() != true {
		t.Fatalf("Bad zero on json decode. U2 should be zero. Found=%s from %q", d4.U2.String(), d4.U2.ustr)
	}
	if len(d4.UArray) != 3 {
		t.Fatalf("Bad json decode. D4 Array len != 3")
	}
	if d4.UArray[0].IsZero() != true {
		t.Fatalf("Bad zero on json decode. UArray[0] should be zero. Found=%s from %q", d4.UArray[0].String(), d4.UArray[0].ustr)
	}
	if d4.UArray[1].IsZero() != true {
		t.Fatalf("Bad zero on json decode. UArray[1] should be zero. Found=%s from %q", d4.UArray[1].String(), d4.UArray[1].ustr)
	}
}

const jtest1 string = `
{
    "U1": "1234567890abcdeedCBA9876543210AB",
    "U2": "1234567890ABCdeedCBA987654321012",
    "UArray": [
        "1234567890abcDEEDCBA9876543210AB",
        "1234567890ABCdeedcbA987654321012"
    ]
}
`
const jtest2 string = `
{
    "U1": "",
    "UArray": [ "",
        "00000000000000000000000000000000",
        "1234567890ABCDEEDCBA987654321012"
    ]
}
`
