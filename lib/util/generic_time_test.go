// --------------------------------------------------------------------
// generic_time_test.go -- Generic Time format parser -- tests.
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package util

import (
	//"fmt"
	"testing"
	//"time"
)

func run_time_test(t *testing.T, tststr string, year, month, day, hour, min, sec, nano int) {
	tme, err := ParseGenericTime(tststr)
	if err != nil {
		t.Fatalf("Time test failed with str = %q.", tststr)
	}
	iy, im, id := tme.Date()
	if iy != year || int(im) != month || id != day {
		t.Fatalf("Time test failed with str = %q. Date not match %d-%d-%d.", tststr, iy, im, id)
	}
	if tme.Hour() != hour {
		t.Fatalf("Time test failed with str = %q. Hour not match (%d)", tststr, tme.Hour())
	}
	if tme.Minute() != min {
		t.Fatalf("Time test failed with str = %q. Minute not match (%d)", tststr, tme.Minute())
	}
	if tme.Second() != sec {
		t.Fatalf("Time test failed with str = %q. Second not match (%d)", tststr, tme.Second())
	}
	if tme.Nanosecond() != nano {
		t.Fatalf("Time test failed with str = %q. Second not match (%d)", tststr, tme.Nanosecond())
	}
}

func Test_ParseGenericTime(t *testing.T) {

	// ed := time.Date(2018, time.February, 27, 11, 18, 46, 321*1000*1000, time.UTC)
	// for _, ft := range GenericTimeFormats {
	// 	fmt.Printf("%q > %s\n", ft, ed.Format(ft))
	// }

	run_time_test(t, "2018", 2018, 1, 1, 0, 0, 0, 0)
	run_time_test(t, "18", 2018, 1, 1, 0, 0, 0, 0)
	run_time_test(t, "Apr 18", 2018, 4, 1, 0, 0, 0, 0)
	run_time_test(t, "April 18", 2018, 4, 1, 0, 0, 0, 0)
	run_time_test(t, "Apr 2018", 2018, 4, 1, 0, 0, 0, 0)
	run_time_test(t, "April 18", 2018, 4, 1, 0, 0, 0, 0)

	run_time_test(t, "2018-05-27", 2018, 5, 27, 0, 0, 0, 0)
	run_time_test(t, "18-05-27", 2018, 5, 27, 0, 0, 0, 0)
	run_time_test(t, "5/27/18", 2018, 5, 27, 0, 0, 0, 0)
	run_time_test(t, "5/27/2018", 2018, 5, 27, 0, 0, 0, 0)
	run_time_test(t, "Feb 27, 2018", 2018, 2, 27, 0, 0, 0, 0)
	run_time_test(t, "Feb 27, 18", 2018, 2, 27, 0, 0, 0, 0)
	run_time_test(t, "February 27, 2018", 2018, 2, 27, 0, 0, 0, 0)
	run_time_test(t, "February 27, 18", 2018, 2, 27, 0, 0, 0, 0)
	run_time_test(t, "2018-02-27T11:18:46", 2018, 2, 27, 11, 18, 46, 0)
	run_time_test(t, "2018-02-27 11:18:46", 2018, 2, 27, 11, 18, 46, 0)
	run_time_test(t, "2018-2-27T11:18:46", 2018, 2, 27, 11, 18, 46, 0)
	run_time_test(t, "2018-2-27 11:18:46", 2018, 2, 27, 11, 18, 46, 0)
	run_time_test(t, "2/27/18 11:18:46", 2018, 2, 27, 11, 18, 46, 0)
	run_time_test(t, "2/27/2018 11:18:46", 2018, 2, 27, 11, 18, 46, 0)

	run_time_test(t, "2018-02-27T11:6:46.321", 2018, 2, 27, 11, 6, 46, 321*1000*1000)
	run_time_test(t, "2018-02-27 11:6:46.321", 2018, 2, 27, 11, 6, 46, 321*1000*1000)
	run_time_test(t, "2018-2-27T11:6:46.321", 2018, 2, 27, 11, 6, 46, 321*1000*1000)
	run_time_test(t, "2018-2-27 11:6:46.321", 2018, 2, 27, 11, 6, 46, 321*1000*1000)
	run_time_test(t, "2/27/18 11:18:6.321", 2018, 2, 27, 11, 18, 6, 321*1000*1000)
	run_time_test(t, "2/27/2018 11:18:46.321", 2018, 2, 27, 11, 18, 46, 321*1000*1000)

}
