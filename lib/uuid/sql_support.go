// --------------------------------------------------------------------
// sql_support.go -- support uuid in sql
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package uuid

import (
	"database/sql/driver"
	"encoding/hex"
	"epic/lib/util"
	"fmt"
	"strings"
)

// Value implements the Valuer DB driver interface for UUID type
func (uuid UUID) Value() (driver.Value, error) {
	bytes, err := hex.DecodeString(uuid.String())
	if err != nil {
		return nil, err
	}
	return driver.Value(bytes), nil
}

// Scan implements the Scanner DB driver interface for UUID type
func (uuid *UUID) Scan(src interface{}) error {
	if src == nil {
		uuid.ustr = "" // THis was changed by DLB to be complient with meaning of zero.
		return nil
	}

	switch srcType := src.(type) {
	case string:
		str := src.(string)
		if util.Blank(str) {
			uuid.ustr = ""
		} else if !IsUuidString(str) {
			return fmt.Errorf("String '%s' is not a valid UUID", str)
		}
		uuid.ustr = strings.ToUpper(str)

	case []byte:
		bytes := src.([]byte)
		if len(bytes) != 16 {
			return fmt.Errorf("Bytes '%x' is not a valid UUID", bytes)
		}
		uuid.ustr = fmt.Sprintf("%X", bytes)
		if uuid.ustr == UuidZeroString {
			uuid.ustr = "" // Make sure zero is encoded correctly.
		}

	default:
		return fmt.Errorf("Incompatible type for UUID %T: %v, %v", srcType, srcType, src)
	}

	return nil
}
