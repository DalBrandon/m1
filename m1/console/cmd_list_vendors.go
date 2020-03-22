// --------------------------------------------------------------------
// cmd_list_vendors.go -- Lists Vendors
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
	"fmt"
	"sort"
)

var gTopic_list_vendors string = `
The list-vendors command is used to list the vendors in the database.
The format of the command is:

  list-vendors

`

func init() {
	RegistorCmd("list-vendors", "", "List the vendors in the database.", handle_list_vendors)
	RegistorTopic("list-vendors", gTopic_list_vendors)
}

func handle_list_vendors(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	vens := m1.GetVendors()
	sort.Slice(vens, func(i, j int) bool { return vens[j].FName > vens[i].FName })

	tbl := util.NewTable("Full Name", "DName", "Default Cat", "N Aliases", "Aliases")
	for _, v := range vens {
		saliases := util.FixStrLen(util.FormatStrSlice(v.Aliases), 50, "...")
		snaliases := util.StrLeft(fmt.Sprintf("%d", len(v.Aliases)), 10)
		tbl.AddRow(v.FName, v.DName, v.DefaultCat, snaliases, saliases)
	}
	c.Printf("%s\n", tbl.Text())
}
