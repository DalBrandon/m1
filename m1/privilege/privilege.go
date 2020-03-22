// --------------------------------------------------------------------
// privilege.go -- package to define privilege levels.
//
// Created 2018-10-04 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

// The privilege package defines the constants used to store
// privilege levels.  It is done here so that other packages can
// share them and not be dependent on one another.
package privilege

import (
	"fmt"
	"strings"
)

type Privilege int

const (
	Admin Privilege = 3
	User  Privilege = 2
	Guest Privilege = 1
	None  Privilege = 0
)

func (priv Privilege) IsAdmin() bool {
	return priv == Admin
}

func (priv Privilege) HasWritePrivilege() bool {
	if priv == Admin || priv == User {
		return true
	}
	return false
}

func (priv Privilege) HasReadPrivilege() bool {
	if priv == Admin || priv == User || priv == Guest {
		return true
	}
	return false
}

func StrToPrivilege(s string) (Privilege, error) {

	s = strings.ToLower(s)
	if s == "admin" || s == "administration" {
		return Admin, nil
	}
	if s == "user" {
		return User, nil
	}
	if s == "guest" {
		return Guest, nil
	}
	if s == "none" || s == "" {
		return None, nil
	}
	return None, fmt.Errorf("Unknown privilege.")
}

// String returns a short string which describes the privilege level.
func (priv Privilege) String() string {
	switch priv {
	case Admin:
		return "Admin"
	case User:
		return "User"
	case Guest:
		return "Guest"
	case None:
		return "None"
	default:
		return "??"
	}
}
