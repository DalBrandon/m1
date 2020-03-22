// --------------------------------------------------------------------
// cmd_pwhash.go -- Prints a password hash
//
// Created 2018-09-29 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/pwhash"
	"dbe/lib/util"
)

func init() {
	RegistorCmd("pwhash", "pw", "Computes password hash (for development).", handle_pwhash)
}

func handle_pwhash(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enough args.\n")
		return
	}
	pw := args[1]
	hash, err := pwhash.HashPassword(pw)
	if err != nil {
		c.Printf("Err: %v\n", err)
		return
	}
	c.Printf("Hash=%q\n", hash)
}
