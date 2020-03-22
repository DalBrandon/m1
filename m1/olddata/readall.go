// --------------------------------------------------------------------
// readall.go -- Loads all old data from csv files.
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

package olddata

import (
	"dbe/lib/log"
	"dbe/lib/util"
)

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
		t, err = LoadTransactions(c, olddata_folder+fn, t)
		if err != nil {
			log.Errorf("Unable to read oldata file %s. Err=%v", olddata_folder+fn, err)
			c.Printf("Error in file %s.\nErr=%v\nOperation hatled.\n", fn, err)
			break
		}
		c.Printf("Read %s. Records = %d\n", fn, len(t)-n)
		n = len(t)
	}
	return t, nil
}
