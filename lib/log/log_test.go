// --------------------------------------------------------------------
// log_test.go -- test some of the log features.
//
// Created 2018-09-21 DLB
// --------------------------------------------------------------------

package log

import (
	"fmt"
	"testing"
)

func Test_LogWriting(t *testing.T) {
	fmt.Printf("Writing some log messages.\n")
	UseConsole(true)
	Passf("This is a foreign message without a newline.")
	Passf("This is a foreign message with a newline.\n")
	Debugf("This is a debug message.")
	Infof("This is a info message.")
	Warnf("This is a warn message.")
	Errorf("This is a error message.")
}
