// cmd_delete_account_alias.go -- Delete an alias from an account
//
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"dbe/m1/m1sql"
	"strconv"
)

var gTopic_delete_alias_from_account string = `
The delete-alias-from-accont command is used to delete an alias from an account.
The format of the command is:

  delete-alias-from-account aid alias="xxxx"

where aid is the account ID, and the alias argument is the name of the
alias. 
`

func init() {
	RegistorCmd("delete-alias-from-account", "", "Delete an alias from an account.", handle_delete_alias_from_account)
	RegistorTopic("delete-alias-from-account", gTopic_delete_alias_from_account)
}

func handle_delete_alias_from_account(c *util.Context, cmdline string) {
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
	n, err := m1sql.DeleteAccountAlias(aid, name)
	if err != nil {
		c.Printf("Error -- %v\n", err)
		return
	}
	if n == 0 {
		c.Printf("Success (no data deleted).\n")
	} else if n == 1 {
		c.Printf("Success (one record deleted).\n")
	} else {
		c.Printf("Success (%d records deleted).\n", n)
	}
}
