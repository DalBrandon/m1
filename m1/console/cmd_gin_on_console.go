// --------------------------------------------------------------------
// cmd_gin_on_console.go -- show or hide gin output on the console
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/m1/config"
	"strings"
)

const gTopic_gin_on_console = `
The gin-on-console command enables or disables writing gin messages to the console.
The format of the command is:

    gin-on-console [true/false]

Note: this command can only be issued at a local connection.  It will not work
over http.  And, over a http connection, gin messages are never written to the terminal.

Background: gin messages are a class of log messages that are generated (by the gin 
package) at a lower level than most other log messages.  Gin messages are automattically sent
to the log files along with the normal log messages.  This command is only concerned
with showing these on the console.  Also note: if normal log messages are not being
shown on the console, then gin messages won't be shown either, no matter how gin-on-console
mode is set.

`

func init() {
	RegistorCmd("gin-on-console", "", "Shows or hides the gin output as it occurs.", handle_gin_on_console)
	RegistorTopic("gin-on-console", gTopic_gin_on_console)
}

func handle_gin_on_console(c *util.Context, cmdline string) {
	if c.IsExternal() {
		c.Printf("Cannot use this command from an external connection.\n")
		return
	}
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enough args.\n")
		return
	}
	mode := strings.ToLower(args[1])
	if mode == "true" {
		log.AllowPassOnConsole(true)
		config.SetBoolParam("gin_on_console", true)
		c.Printf("Success.\n")

	} else if mode == "false" {
		log.AllowPassOnConsole(false)
		config.SetBoolParam("gin_on_console", false)
		c.Printf("Success.\n")

	} else {
		c.Printf("True or False not found.\n")
		return
	}
}
