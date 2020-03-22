// --------------------------------------------------------------------
// cmd_list_transactions.go -- Lists Transactions
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
	"sort"
	"strconv"
)

var gTopic_list_transactions string = `
The list-transactions command is used to list the transactions in the database.
The format of the command is:

  list-transactions max=nnn skip=nnn

where nnn is the max number of transactions listed. The default for max is 100.
The skip parameter is optional, and if given, the first nnn records will be skipped.

`

func init() {
	RegistorCmd("list-transactions", "", "List the tranactions in the database.", handle_list_transactions)
	RegistorTopic("list-transactions", gTopic_list_transactions)
}

func handle_list_transactions(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	maxlst := 100
	smax, ok := util.MapAlias(params, "max")
	if ok {
		maxlst, err = strconv.Atoi(smax)
		if err != nil {
			c.Printf("Invalid paramger for max. (%s), Err=%v\n", smax, err)
			return
		}
	}
	iSkip := 0
	sskip, ok := util.MapAlias(params, "skip", "start")
	if ok {
		iSkip, err = strconv.Atoi(sskip)
		if err != nil {
			c.Printf("Invalid parameter for skip (%s), Err=%v\n", sskip, err)
			return
		}
	}

	tlst := m1.GetTransactions()
	sort.Slice(tlst, func(i, j int) bool {
		return tlst[j].Date().Format("06-01-02") >
			tlst[i].Date().Format("06-01-02")
	})

	tbl := util.NewTable("Date", "Account", "Vendor", "Description", "Cat", "Amount")

	icnt := 0
	nrows := 0
	for _, t := range tlst {
		if icnt >= iSkip {
			sscat := ""
			if len(t.Cats) > 0 {
				sscat = t.Cats[0].Category
			}
			samt := util.StrLeft(util.CentsToStr(t.Amount), 14)
			tbl.AddRow(t.Date().Format("06-01-02"), t.Account, t.Vendor,
				t.Description, sscat, samt)
			nrows += 1
		}
		if nrows > maxlst {
			break
		}
		icnt += 1
	}
	c.Printf("%s\n", tbl.Text())
}
