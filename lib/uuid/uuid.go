// --------------------------------------------------------------------
// uuid.go -- creates and manuplates uuids
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package uuid

import (
	"encoding/json"
	"epic/lib/util"
	"fmt"
	realuuid "github.com/pborman/uuid"
	"strconv"
	"strings"
	"unicode/utf8"
)

type UUID struct {
	ustr string
}

const (
	// ----------------------12345678901234567890123456789012
	UuidZeroString string = "00000000000000000000000000000000"
)

// Returns a zero valued uuid.
func Zero() UUID {
	return UUID{ustr: ""}
}

// New returns a new UUID object.  Each call to New will return
// a different UUID. Useful for database keys in mysql.  Not for security.
func New() UUID {
	var u UUID
	id := realuuid.NewRandom()
	s := id.String()
	u.ustr = strings.ToUpper(strings.Replace(s, "-", "", -1))
	return u
}

// FromString attemps to build a UUID from an input string.  The input
// must be 32 Hex characters without dashes, or an error is returned.
func FromString(s string) (u UUID, err error) {
	if IsUuidString(s) {
		u.ustr = strings.ToUpper(s)
		if u.ustr == UuidZeroString {
			// Disallow storage of zero string.
			u.ustr = ""
		}
		return u, nil
	}
	return u, fmt.Errorf("Input \"%s\" is not a valid 32 char hex number for UUID.", s)
}

// FromString0 attemps to build a UUID from an input string.  It is
// similar to FromString execpt that if the input is compeletely blank,
// the uuid.Zero() is returned.
func FromString0(s string) (u UUID, err error) {
	if util.Blank(s) {
		return Zero(), nil
	}
	return FromString(s)
}

// ForceStr builds a UUID from an input string, ingoring any errors.  If the
// string is not a valid UUID, the zero uuid is returned.
func ForceStr(s string) UUID {
	u, _ := FromString(s)
	return u
}

// IsZero returns true if the UUID has a zero value.
func (u UUID) IsZero() bool {
	if len(u.ustr) <= 0 {
		return true
	}
	if u.ustr == UuidZeroString {
		return true
	}
	return false
}

// String returns the string representation of the UUID.  The string
// representation is 32 Hex chars, suitable for use in the sql's UNHEX()
// function.
func (u UUID) String() string {
	if len(u.ustr) <= 0 {
		return UuidZeroString
	} else {
		return u.ustr
	}
}

// MarshalJSON is used by the JSON encoder to write this type.
func (u UUID) MarshalJSON() ([]byte, error) {
	s := u.String()
	return json.Marshal(s)
}

// UnmarshalJSON is used by the JSON encoder to read this type.
func (u *UUID) UnmarshalJSON(b []byte) (err error) {
	// Blank types are okay, and treated as zero.
	if strings.TrimSpace(string(b)) == "{}" {
		u.ustr = ""
		return nil
	}
	// Something was transmitted. Try to decode.
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		u.ustr = ""
		return err
	}
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		u.ustr = ""
		return err
	}
	if !IsUuidString(s) {
		u.ustr = ""
		return fmt.Errorf("JSON unmarshal error. String not a 32 char Hex Key: \"%s\"", s)
	}
	u.ustr = strings.ToUpper(s)
	if u.ustr == UuidZeroString {
		// Disallow storage of zero string.
		u.ustr = ""
	}
	return nil
}

// GobEncode is used by the Gob interface to write this type.
func (u UUID) GobEncode() ([]byte, error) {
	// Add a 0x01 byte to the beginning of the byte slice to act as a version
	// number of this object.
	s := u.String()
	b := make([]byte, 0, len(s)+1)
	b = append(b, 0x01)
	b = append(b, []byte(s)...)
	//fmt.Printf("Encoded Gob UUID (%d) = %v\n\n", len(b), b)
	return b, nil
}

// GobDecode is used by the Gob interface to read this type.
func (u *UUID) GobDecode(b []byte) (err error) {
	u.ustr = "" // Assume failure.
	if len(b) == 0 {
		// Blank and empty types okay.  Treat as zero.
		return nil
	}
	if len(b) <= 32 {
		return fmt.Errorf("Gob decode error. Not enough bytes. Expected >32, got %d.", len(b))
	}
	if b[0] != 0x01 {
		return fmt.Errorf("Gob decode error.  Unknown version of uuid object. Found version=%d", b[0])
	}
	s := strings.TrimSpace(string(b[1:]))
	// Blank types are okay, and treated as zero.
	if s == "" {
		u.ustr = ""
		return nil
	}
	if !IsUuidString(s) {
		return fmt.Errorf("Gob decode error. String not a 32 char Hex Key: \"%s\"", s)
	}
	u.ustr = strings.ToUpper(s)
	if u.ustr == UuidZeroString {
		// Disallow storage of zero string.
		u.ustr = ""
	}
	return nil
}

// IsUuidString checks to make sure the input is a 32 char hex
// string, suitable for building UUIDs.
func IsUuidString(s string) bool {
	hexdigits := "0123456789ABCDEF"
	s = strings.ToUpper(s)
	if utf8.RuneCountInString(s) != 32 {
		return false
	}
	for _, c := range s {
		c, _ := strconv.Unquote(strconv.QuoteRune(c))
		if !strings.Contains(hexdigits, string(c)) {
			return false
		}
	}
	return true
}
