// --------------------------------------------------------------------
// util.go -- Utilities for Console
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"fmt"
	"os"
	"path"
	"strings"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func strlist(ss []string) string {
	sout := ""
	for i, s := range ss {
		if i != 0 {
			sout += ", "
		}
		sout += s
	}
	return sout
}

func ParseCmdLine(input string, params map[string]string) (args []string, err error) {
	switches := make([]string, len(params))
	for k, _ := range params {
		switches = append(switches, k)
	}
	cwords := util.SplitArgs(input)
	args = make([]string, 0, len(cwords))
	for _, w := range cwords {
		ww := w
		if len(ww) >= 2 && strings.HasPrefix(ww, "\"") && strings.HasSuffix(ww, "\"") {
			ww = strings.TrimPrefix(ww, "\"")
			ww = strings.TrimSuffix(ww, "\"")
		} else {
			ww = strings.TrimSpace(w)
		}
		isswitch := false
		for _, sw := range switches {
			if ww == sw {
				isswitch = true
				break
			}
		}
		if isswitch {
			params[ww] = "true"
			continue
		}
		p2 := strings.Split(ww, "=")
		if len(p2) == 1 {
			args = append(args, ww)
		} else if len(p2) == 2 {
			s := p2[1]
			if len(s) >= 2 && strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
				s = strings.TrimPrefix(s, "\"")
				s = strings.TrimSuffix(s, "\"")
			}
			params[p2[0]] = s
		} else {
			return []string{}, fmt.Errorf("Parameter syntax error for %q.", w)
		}
	}
	return args, nil
}

func CreateHistoryFile() (*os.File, error) {
	dirpart := path.Dir(gHistoryFile)
	if dirpart != "." && dirpart != "/" {
		if !util.DirExists(dirpart) {
			err := os.Mkdir(dirpart, 0775)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Create(gHistoryFile)
}

func ParseActive(params map[string]string) (active, doactive bool, err error) {
	sactive, doactive := util.MapAlias(params, "Active", "active")
	if doactive {
		active, err = util.StrToBool(sactive, false)
		if err != nil {
			return false, false, fmt.Errorf("Value for active unrecognizable (%q)\n", sactive)
		}
		return active, true, nil
	}
	return false, false, nil
}

func fixstrfor_e(s string) string {
	sout := ""
	for _, y := range s {
		if y < 0 {
			continue
		}
		if y < 127 {
			sout = sout + string(y)
		}
		if y == 65533 {
			sout = sout + "e" // Fixes the 'cafe' problem introduced by excell
		}
	}
	return sout
}
