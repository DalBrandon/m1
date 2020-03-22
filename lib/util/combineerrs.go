// --------------------------------------------------------------------
// combinerrs.go -- Combine errors.
//
// Created 2018-09-23 DLB
// --------------------------------------------------------------------

package util

import (
	"fmt"
)

func CombineErrs(errs ...error) error {
	if len(errs) <= 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}

	var err1 error
	icnt := 0
	for _, e := range errs {
		if e != nil {
			icnt += 1
			if icnt == 1 {
				err1 = e
			}
		}
	}
	if icnt <= 0 {
		return nil
	}
	err := fmt.Errorf("Mult errs (%d). First Err: %v\n", icnt, err1)
	return err
}
