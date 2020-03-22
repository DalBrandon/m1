// --------------------------------------------------------------------
// pwhash.go -- Password hashing algorithm
//
// Created 2018-09-29 DLB
// --------------------------------------------------------------------

// Uses the current gold web standard: bcrypt.  However, since not
// sure what bcrypt will return, does a base64 encode on hash to make
// sure it can be stored and transmitted as a string.

package pwhash

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashbytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}
	hash := base64.StdEncoding.EncodeToString(hashbytes)
	return hash, nil
}

func CheckPasswordHash(password, hash string) bool {
	hashbytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(hashbytes, []byte(password))
	return err == nil
}
