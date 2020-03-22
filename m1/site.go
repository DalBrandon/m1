// --------------------------------------------------------------------
// site.go -- Start of the m1 project
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------
package main

// import (
// 	"dbe/lib/util"
// 	"dbe/m1/m1sql"
// 	"dbe/m1/olddata"
// 	"fmt"
// 	"os"
// )

// func main() {
// 	fmt.Printf("Brandon's Money Manger (M1)\n")
// 	fmt.Printf("Opening database.\n")
// 	err := m1sql.OpenDatabase("m1fun")
// 	if err != nil {
// 		fmt.Printf("Error: %q\n", err)
// 	} else {
// 		fmt.Printf("Database successfully opened.\n")
// 	}

// 	t, err := olddata.ReadAll()
// 	if err != nil {
// 		fmt.Printf("Error: %q\n", err)
// 		os.Exit(0)
// 	}
// 	fmt.Printf("Number of records imported = %d\n", len(t))
// 	mcat := CatList(t)
// 	tbl := util.NewTable("Category", "Items")
// 	for c, i := range mcat {
// 		s := fmt.Sprintf("%5d", i)
// 		tbl.AddRow(c, s)
// 	}
// 	fmt.Printf("%s", tbl.Text())
// }

// func CatList(r []*olddata.OldTransaction) map[string]int {
// 	m := make(map[string]int, 50)
// 	for _, x := range r {
// 		c := x.Category
// 		m[c] = m[c] + 1
// 	}
// 	return m
// }
