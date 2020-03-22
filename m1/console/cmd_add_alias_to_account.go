// --------------------------------------------------------------------
// cmd_add_account_alias.go -- Add an alias to an account
//
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"dbe/m1/m1sql"
	"strconv"
)

var gTopic_add_alias_to_account string = `
The add-alias-to-accont command is used to add an alias to an account.
The format of the command is:

  add-alias-to-account aid alias="xxxx" notes="xxxx"

where aid is the account ID, and the alias argument is the name of the
alias.  The note argument is optional. 
`

func init() {
	RegistorCmd("add-alias-to-account", "", "Add an alias to an account.", handle_add_alias_to_account)
	RegistorTopic("add-alias-to-account", gTopic_add_alias_to_account)
}

func handle_add_alias_to_account(c *util.Context, cmdline string) {
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
	aid, err := strconv.Atoi(args[1])
	name, ok := util.MapAlias(params, "Alias", "Name", "alias", "name")
	if !ok || util.Blank(name) {
		c.Printf("An Alias must be provided.\n")
		return
	}
	notes, _ := util.MapAlias(params, "Notes", "notes")
	err = m1sql.AddAccountAlias(aid, name, notes)
	if err != nil {
		c.Printf("Error -- %v\n", err)
		return
	}
	c.Printf("Success.\n")
}
