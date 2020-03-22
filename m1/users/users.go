// --------------------------------------------------------------------
// users.go -- Manage users for m1
//
// Created 2020-03-15 DLB
// --------------------------------------------------------------------

package users

import (
	"dbe/lib/pwhash"
	pv "dbe/m1/privilege"
)

type UserWithPW struct {
	Name    string
	Account string
	Priv    pv.Privilege
	Pwhash  string
}

type User struct {
	Name    string
	Account string
}

var dal UserWithPW = UserWithPW{"Dal Brandon", "dal", pv.Admin, "JDJhJDA2JE9kRm9RUFpORXdiaDRDVXU3TWpyWU80NEI2ancvb09ySnBjMWtwM3ovNFJQaVhVazc0a0Eu"}
var carol UserWithPW = UserWithPW{"Carol Brandon", "carol", pv.Admin, "JDJhJDA2JE9kRm9RUFpORXdiaDRDVXU3TWpyWU80NEI2ancvb09ySnBjMWtwM3ovNFJQaVhVazc0a0Eu"}
var users = []UserWithPW{dal, carol}

func GetUsers() []User {
	ulst := make([]User, 0, 5)
	for _, u := range users {
		uu := User{Name: u.Name, Account: u.Account}
		ulst = append(ulst, uu)
	}
	return ulst
}

func GetUser(name string) *UserWithPW {
	for _, u := range users {
		if name == u.Name {
			uout := &UserWithPW{Name: u.Name, Account: u.Account, Priv: u.Priv}
			return uout
		}
	}
	return nil
}

func CheckPassword(name string, pw_in_clear string) bool {
	for _, u := range users {
		if name == u.Name {
			return pwhash.CheckPasswordHash(pw_in_clear, u.Pwhash)
		}
	}
	return false
}
