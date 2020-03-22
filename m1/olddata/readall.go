// --------------------------------------------------------------------
// readall.go -- Loads all old data from csv files.
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

package olddata

import (
	"dbe/lib/util"
)

var folder string = "/home/dal/dev/src/dbe/m1/olddata/"
var inputfiles []string = []string{
	"CarolsVisa.csv",
	"BofAChecking.csv",
	"DalsVisa.csv",
	"FRBChecking.csv",
	"MLChecking.csv",
}

func ReadAll(c *util.Context) (records []*OldTransaction, err error) {
	t := make([]*OldTransaction, 0, 20000)
	n := 0
	for _, fn := range inputfiles {
		t, err = LoadTransactions(c, folder+fn, t)
		if err != nil {
			c.Printf("Error in file %s.\nErr=%v\nOperation hatled.\n", fn, err)
			break
		}
		c.Printf("Read %s. Records = %d\n", fn, len(t)-n)
		n = len(t)
	}
	return t, nil
}
