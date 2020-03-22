// --------------------------------------------------------------------
// cmd_save_data.go -- Saves the database to a disk file.
//
// Created 2020-03-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
)

var gTopic_save_data string = `
The save-data command saves the database to a disk file.
The format of the command is:

  save-data

Without any arguments, the data saved will become the 
default data to load the next time the server starts.
`

func init() {
	RegistorCmd("save-data", "", "Saves data to disk.", handle_save_data)
	RegistorTopic("save-data", gTopic_save_data)
}

func handle_save_data(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	err = m1.SaveData()
	if err != nil {
		c.Printf("Error = %v.\n", err)
		return
	}
	c.Printf("Success.\n")
}
