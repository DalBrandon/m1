// --------------------------------------------------------------------
// cmd_list_accounts.go -- Lists Raw data for accounts
//
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"dbe/m1/m1data"
	"sort"
)

var gTopic_listaccounts string = `
The list-accounts command is used to list the accounts in the database.
The format of the command is:

  list-accounts  

`

func init() {
	RegistorCmd("list-accounts", "", "Analyze Old Data.", handle_listaccounts)
	RegistorTopic("list-accounts", gTopic_listaccounts)
}

func handle_listaccounts(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	accounts := m1data.GetAccounts()
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[j].ShortName > accounts[i].ShortName
	})

	tbl := util.NewTable("ShortName", "Display Name", "Full Name", "Active", "Aliases")
	for _, a := range accounts {
		sactive := "Yes"
		if !a.Active {
			sactive = "No"
		}
		saliases := util.FormatStrSlice(a.Aliases)
		tbl.AddRow(a.ShortName, a.DName, a.FName, sactive, saliases)
	}
	c.Printf("%s\n", tbl.Text())
}
