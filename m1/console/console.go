// --------------------------------------------------------------------
// console.go -- Console loop for server
//
// Created 2018-09-21 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"sort"
	"strings"
)

type Command struct {
	CmdName  string
	ArgLine  string
	HelpLine string
	Handler  func(c *util.Context, cmd string)
}

var gCmds []*Command
var gConsole *liner.State
var gPrompt string = "> "
var gHistoryFile string = "./tmp/console_history.txt"

func cmdsorter(i, j int) bool {
	s1, s2 := gCmds[i].CmdName, gCmds[j].CmdName
	if s1 == "?" {
		return false
	}
	if s2 == "?" {
		return true
	}
	return s1 < s2
}

// Registor a command that will be avaliable at the console.
func RegistorCmd(name, argline, helpline string, handler func(*util.Context, string)) {
	if gCmds == nil {
		gCmds = make([]*Command, 0, 10)
	}
	cmd := &Command{name, argline, helpline, handler}
	gCmds = append(gCmds, cmd)
	sort.Slice(gCmds, cmdsorter)
}

func SetPrompt(prompt string) {
	gPrompt = prompt
}

// SetHistoryFile sets the filename that will be used to store perminent history.
func SetHistoryFile(filename string) {
	gHistoryFile = filename
}

// ExecuteCommand will execute a command outside a command loop.  Suitable for
// a web interface.  No history is maintained.
func ExecuteCommand(cmdline string) string {
	c := util.NewContext(util.Context_External)
	execute_cmd(c, cmdline)
	return c.Output()
}

// Start a console loop.  Runs forever.  Should be called as a go func.
func ConsoleLoop() {
	gConsole = liner.NewLiner()

	// Note, these two defers will only get executed if some other part of
	// the program calls os.Exit()...  If we exit below by calling os.Exit()
	// they will be ignored, since these are in the same nested call tread.
	defer gConsole.Close()
	defer write_history()

	gConsole.SetCtrlCAborts(false)
	gConsole.SetCompleter(func(line string) (lst []string) {
		lst = make([]string, 0, 10)
		for _, c := range gCmds {
			if strings.HasPrefix(c.CmdName, line) {
				lst = append(lst, c.CmdName)
			}
		}
		return lst
	})
	if f, err := os.Open(gHistoryFile); err == nil {
		gConsole.ReadHistory(f)
		f.Close()
	}
	if !liner.TerminalSupported() {
		fmt.Printf("Terminal not supported. Editting commands won't work.\n")
	}
	fmt.Printf("Use 'help' for a list of commands.\n")
	for {
		cmdline, err := gConsole.Prompt(gPrompt)
		if err == liner.ErrPromptAborted {
			continue
		}
		if err != nil {
			fmt.Printf("Input error: %v\n", err)
			continue
		}
		if !util.Blank(cmdline) {
			gConsole.AppendHistory(cmdline)
			c := util.NewContext(util.Context_Internal)
			c.SetFlusher(func() {
				fmt.Printf("%s", c.Output())
				c.Reset()
			})

			execute_cmd(c, cmdline)
			fmt.Printf("%s", c.Output())
		}
	}
}

func write_history() {
	f, err := CreateHistoryFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create console history file.\nErr=%v\n", err)
		return
	}
	gConsole.WriteHistory(f)
	f.Close()
}

func execute_cmd(c *util.Context, cmdline string) {

	cmdline = strings.TrimSpace(cmdline)
	if cmdline == "" {
		return
	}
	cmdwords := strings.Split(cmdline, " ")
	if len(cmdwords) <= 0 {
		return
	}
	cmd := strings.TrimSpace(cmdwords[0])
	if cmd == "exit" || cmd == "quit" {
		if c.IsExternal() {
			c.Printf("Cannot use this command from an external connection.\n")
			return
		}
		write_history()
		gConsole.Close()
		os.Exit(0)
		return
	}
	for _, x := range gCmds {
		if cmd == x.CmdName {
			x.Handler(c, cmdline)
			return
		}
	}
	c.Printf("Unknown Command. Use ? or help.\n")
	return
}
