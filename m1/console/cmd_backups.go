// --------------------------------------------------------------------
// cmd_backups.go -- List the backups found in the backup folder.
//
// Created 2020-03-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	m1 "dbe/m1/m1data"
	"sort"
)

var gTopic_list_backups string = `
The list-backups command lists the backup files avaliable for
loading. The format of the command is:

  list-backups

`
var gTopic_make_backup string = `
The make-backup command saves a snapshot of the database to
the backup directory.  The format of the command is:

  make-backup fname

where fname is the root name of a file without any path or
extension.  The fname argument can be omitted in whichcase
a file named with the current time will be used.

`
var gTopic_load_backup string = `
The list-backups command lists the backup files that have been
written manually.  The format of the command is:

  load-backup fname

where fname is the name of a saved backup file, obtained with
list-backups command.  

WARNING: this command will overwrite the current database
with the backuped data, and can lead to a loss of data.
It is usually best to do a make-backup right before issuing
this command so that the current database can be restored.

`
var gTopic_delete_backup string = `
The delete-backup command will delete a backup. The format
of the command is:

  delete-backup fname

where fname is the name of the backup file to delete.  

WARNING: there is no recourse after the backup file
is deleted.

`

func init() {
	RegistorCmd("list-backups", "", "Lists the backup files.", handle_list_backups)
	RegistorCmd("make-backup", "", "Makes a backup.", handle_make_backup)
	RegistorCmd("load-backup", "", "Loads a backup.", handle_load_backup)
	RegistorCmd("delete-backup", "", "Deletes a backup.", handle_delete_backup)
	RegistorTopic("list-backups", gTopic_list_backups)
	RegistorTopic("make-backup", gTopic_make_backup)
	RegistorTopic("load-backup", gTopic_load_backup)
	RegistorTopic("deete-backup", gTopic_delete_backup)
}

func handle_list_backups(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	_, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}

	lst := m1.GetBackupFileList()
	if err != nil {
		c.Printf("Error = %v.\n", err)
		return
	}
	sort.Strings(lst)
	tbl := util.NewTable("File Name")
	for _, f := range lst {
		tbl.AddRow(f)
	}
	c.Printf("%s\n", tbl.Text())
}

func handle_make_backup(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	fname := ""
	if len(args) > 1 {
		fname = args[1]
	}
	if !util.Blank(fname) && !legalFileRootName(fname) {
		c.Printf("Input file name (%q) has illegal characters.\n")
		return
	}
	if !util.Blank(fname) {
		lst := m1.GetBackupFileList()
		if util.InStringSlice(lst, fname) {
			c.Printf("That backup file (%s) already exists.\n", fname)
			return
		}
	}
	fout, err := m1.SaveBackup(fname)
	if err != nil {
		c.Printf("Error: %v\n", err)
		return
	}
	c.Printf("Backup %s created.\n", fout)
	c.Printf("Success.\n")
}

func handle_load_backup(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Filename not provided.\n")
		return
	}
	fname := args[1]
	err = m1.LoadBackup(fname)
	if err != nil {
		c.Printf("Error: %v\n", err)
		return
	}
	c.Printf("Success.\n")
}

func handle_delete_backup(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Filename not provided.\n")
		return
	}
	fname := args[1]
	err = m1.DeleteBackup(fname)
	if err != nil {
		c.Printf("Error: %v\n", err)
		return
	}
	c.Printf("Success.\n")
}

func legalFileRootName(fname string) bool {
	for i, x := range fname {
		if x >= 'a' && x <= 'z' {
			continue
		}
		if x >= 'A' && x <= 'Z' {
			continue
		}
		if x >= '0' && x <= '9' {
			continue
		}
		if i > 0 && (x == '-' || x == '_') {
			continue
		}
		return false
	}
	return true
}
