// --------------------------------------------------------------------
// log.go -- Custom logger stuff for epic web sites
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package log

import (
	"epic/lib/util"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

var m_console bool = true
var m_allowPass bool = true
var m_file bool = true
var m_errcnt int
var m_filelock sync.Mutex

func init() {
	fi, err := os.Stat("./logs")
	if err != nil {
		err := os.Mkdir("./logs", 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create logs directory. Err=%v\n", err)
			panic(err)
		}
		fi, err = os.Stat("./logs")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get info about logs dir. Err=%v\n", err)
			panic(err)
		}
	}
	if !fi.IsDir() {
		err := fmt.Errorf("Unable to create logs directory -- is file instead.")
		fmt.Fprintf(os.Stderr, "%v\n", err)
		panic(err)
	}
}

// UseConsole sets or unsets a mode whereby messages are not sent to the
// terminal as well as the log file.
func UseConsole(mode bool) {
	m_console = mode
}

// AllowPassOnConsole sets a mode whereby messages that don't use the normal
// Debugf, Infof, Warnf, and Errorf are filtered out from going to the
// terminal.  In practice this means that the raw GIN messages may be
// prevented from going to the terminal.
func AllowPassOnConsole(mode bool) {
	m_allowPass = mode
}

func Passf(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("", s)
}

func Debugf(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("debug", s)
}

func Infof(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("info", s)
}

func Warnf(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("warn", s)
}

func Errorf(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("error", s)
}

func Fatalf(ft string, args ...interface{}) {
	s := fmt.Sprintf(ft, args...)
	Logit("fatal", s)
	os.Exit(1)
}

func Logit(level, msg string) {
	t := time.Now()
	var msgout string
	if level == "" {
		msgout = "!" + msg
		if !strings.HasSuffix(msgout, "\n") {
			msgout += "\n"
		}
	} else {
		msec := int(t.Nanosecond() / 1000000)
		msgout = fmt.Sprintf(">%s.%03d [%s] %s\n",
			t.Format("06-01-02 15:04:05"), msec, util.FixStrLen(level, 5, " "), msg)
	}
	if m_console {
		if level != "" || m_allowPass {
			fmt.Fprintf(os.Stderr, msgout)
		}
	}
	if !m_file {
		return
	}

	m_filelock.Lock()
	defer m_filelock.Unlock()
	fn := fmt.Sprintf("./logs/Log_%s.txt", t.Format("06-01-02"))
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		m_errcnt += 1
		if m_errcnt == 1 {
			fmt.Fprintf(os.Stderr, "Unable to write to log file!\nErr=%v\n", err)
			return
		} else if m_errcnt%50 == 0 {
			fmt.Fprintf(os.Stderr, "Still unable to write to log file after 50 more attemps.\nErr=%v\n", err)
			return
		}
	}
	if _, err := f.Write([]byte(msgout)); err != nil {
		m_errcnt += 1
		if m_errcnt == 1 {
			fmt.Fprintf(os.Stderr, "Unable to write to log file!\nErr=%v\n", err)
		} else if m_errcnt%50 == 0 {
			fmt.Fprintf(os.Stderr, "Still unable to write to log file after 50 more attemps.\nErr=%v\n", err)
		}
	}
	if err := f.Close(); err != nil {
		m_errcnt += 1
		if m_errcnt == 1 || m_errcnt%50 == 0 {
			fmt.Fprintf(os.Stderr, "Unable to close log file.\n")
		}
	}
}

func ReadLogFile(date time.Time) (string, error) {
	fn := fmt.Sprintf("./logs/Log_%s.txt", date.Format("06-01-02"))
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
