// --------------------------------------------------------------------
// cmd_load_data.go -- Loads the database.
//
// Created 2020-03-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
)

var gTopic_load_data string = `
The load-data command loads the database from a saved file.
The format of the command is:

  load-data

Without any arguments, this command loads the latest saved
data.  Note: this WILL overwrite changes that have been
made since the last save.
`

func init() {
	RegistorCmd("load-data", "", "Loads a saved database.", handle_load_data)
	RegistorTopic("load-data", gTopic_load_data)
}

func handle_load_data(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	err = m1.LoadData()
	if err != nil {
		c.Printf("Error = %v.\n", err)
		return
	}
	c.Printf("Success.\n")
}
