// --------------------------------------------------------------------
// cmd_help.go -- Commands for help.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"fmt"
	"sort"
	"strings"
)

var gArgs map[string]string
var gTopics map[string]string

func init() {
	RegistorCmd("help", "[topic]", "Gives this help or help on a topic.", handle_help)
	RegistorCmd("?", "", "Condensed help.", handle_help_condensed)
	RegistorArg("topic", "The topic for help. Use 'help topics' to get a list of topics.")
}

func RegistorArg(arg, desc string) {
	if gArgs == nil {
		gArgs = make(map[string]string, 100)
	}
	gArgs[arg] = desc
}

func RegistorTopic(topic, desc string) {
	if gTopics == nil {
		gTopics = make(map[string]string, 100)
	}
	gTopics[topic] = desc
}

func SayHelp() {
	fmt.Printf("%s/n", gHelpText)
}

func show_topic_list(c *util.Context) {
	tlst := make([]string, 0, len(gTopics))
	for k, _ := range gTopics {
		tlst = append(tlst, k)
	}
	sort.Strings(tlst)
	sout := "Topic list: "
	for i, t := range tlst {
		if i != 0 {
			sout += ", "
		}
		if len(sout)+len(t) > 80 {
			c.Printf("%s\n", sout)
			sout = ""
		}
		sout += t
	}
	if len(sout) > 0 {
		c.Printf("%s\n", sout)
	}
}

func handle_help(c *util.Context, cmdline string) {
	cwrds := strings.Split(cmdline, " ")
	if len(cwrds) > 1 {
		top := strings.TrimSpace(cwrds[1])
		if top == "topics" {
			show_topic_list(c)
			return
		}
		desc, ok := gTopics[top]
		if !ok {
			c.Printf("Unknown topic. ")
			show_topic_list(c)
			return
		}
		c.Printf("%s\n%s\n", top, desc)
		return
	}
	if c.IsInternal() {
		c.Printf("You are now in an interactive command loop. The commnds are:\n\n")
	}
	if c.IsExternal() {
		c.Printf("Help Info.\n")
	}
	w1 := 10
	for _, c := range gCmds {
		w1 = max(w1, len(c.CmdName+" "+c.ArgLine+" "))
	}
	w1 = max(w1, len("exit | quit  "))
	for _, x := range gCmds {
		s1 := util.FixStrLen(x.CmdName+" "+x.ArgLine+" ", w1, "")
		c.Printf("    %s -- %s\n", s1, x.HelpLine)
	}
	c.Printf("    %s -- Exits this program.\n", util.FixStrLen("exit | quit  ", w1, ""))
	c.Printf("\n")

	w2 := 10
	argnames := make([]string, 0, len(gArgs))
	for aname, _ := range gArgs {
		w2 = max(w2, len(aname))
		argnames = append(argnames, aname)
	}
	sort.Strings(argnames)
	c.Printf("Where the arguments are:\n")
	for _, aname := range argnames {
		desc := gArgs[aname]
		sn := util.FixStrLen(aname, w2, "")
		c.Printf("    %s -- %s\n", sn, desc)
	}
	c.Printf("\n")
}

func handle_help_condensed(c *util.Context, cmdline string) {
	sout := ""
	w := 0
	i := 0
	for _, c := range gCmds {
		if c.CmdName == "?" {
			continue
		}
		if i != 0 {
			sout += ", "
			w += 2
		}
		if len(c.CmdName)+w > 80 {
			sout += "\n"
			w = 0
		}
		sout += c.CmdName
		w += len(c.CmdName)
		i += 1
	}
	sout += ", exit"
	c.Printf("%s\n", sout)
}

var gHelpText string = `
    m1 -- Brandon's m1 Accounting System -- March 2020
    This server runs on a Droplet in the clowd and maintains
    accounting records for Dal and Carol Brandon.
`
