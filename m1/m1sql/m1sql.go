// --------------------------------------------------------------------
// m1sql.go -- Access to the M1Data database.
//
// Created 2020-03-14 DLB
// --------------------------------------------------------------------

package m1sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var m_db *sql.DB

func OpenDatabase(pw string) error {
	var err error
	connection := fmt.Sprintf("root:%s@/M1Data", pw)
	m_db, err = sql.Open("mysql", connection)
	if err != nil {
		err := fmt.Errorf("Unable to open database. Err=%v\n", err)
		return err
	}
	return nil
}
