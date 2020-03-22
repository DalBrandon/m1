// --------------------------------------------------------------------
// cmd_list_cats.go -- Lists Categories
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
	"sort"
)

var gTopic_list_cats string = `
The list-cats command is used to list the categories in the database.
The format of the command is:

  list-cats

`

func init() {
	RegistorCmd("list-cats", "", "List the categories in the database.", handle_list_cats)
	RegistorTopic("list-cats", gTopic_list_cats)
}

func handle_list_cats(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	cats := m1.GetCategories()
	sort.Slice(cats, func(i, j int) bool { return cats[j].Name > cats[i].Name })

	tbl := util.NewTable("Name", "Aliases")
	for _, cx := range cats {
		saliases := util.FormatStrSlice(cx.Aliases)
		tbl.AddRow(cx.Name, saliases)
	}
	c.Printf("%s\n", tbl.Text())
}
