// --------------------------------------------------------------------
// cmd_load_vendors.go -- Loads the old vendors
//
// Created 2020-03-18 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	"dbe/m1/m1sql"
	"dbe/m1/olddata"
	"sort"
	"strings"
	"time"
)

var gTopic_load_vendors string = `
The load-old-vendors command is used to load the old vendor data into the new
database.  This command should be safe to run anytime, since data is only added,
not deleted.  The format of the command is:

  load-old-vendors update

where "update" is optional, and if not included the old data will simply be displayed
and the actual database will not be updated.

The old data is kept in a csv file that is in the development directory.

`

func init() {
	RegistorCmd("load-old-vendors", "", "Load old vendors into database.", handle_load_old_vendors)
	RegistorTopic("load-old-vendors", gTopic_load_vendors)
}

func handle_load_old_vendors(c *util.Context, cmdline string) {
	type Vendor struct {
		FName   string
		DName   string
		Aliases []string
	}

	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	doupdate := false
	if len(args) > 1 {
		if strings.ToLower(args[1]) == "update" {
			doupdate = true
		} else {
			c.Printf("Unknown argument (%s)\n", args[1])
			return
		}
	}

	newcontext := util.NewContext(util.Context_Internal)
	t0 := time.Now()
	vlst, err := olddata.GetVendorData(newcontext)
	tload := time.Now().Sub(t0)
	c.Printf("Time to Load Data: %d ms.\n", tload.Milliseconds())
	if err != nil {
		c.Printf("Error: %v\n", err)
		return
	}
	c.Printf("Number of Raw Records Read: %d\n", len(vlst))

	vmap := make(map[string]*Vendor, len(vlst))
	for _, vraw := range vlst {
		v, ok := vmap[vraw.FName]
		if !ok {
			vnew := &Vendor{FName: vraw.FName, DName: vraw.DName, Aliases: make([]string, 0, 10)}
			vnew.Aliases = append(vnew.Aliases, vraw.Alias)
			vmap[vnew.FName] = vnew
		} else {
			if !util.InStringSlice(v.Aliases, vraw.Alias) {
				v.Aliases = append(v.Aliases, vraw.Alias)
			} else {
				c.Printf("Duplicate Alias (%q) found for vendor %q.\n", vraw.Alias, v.FName)
			}
		}
	}
	klst := make([]string, 0, len(vmap))
	for k, _ := range vmap {
		klst = append(klst, k)
	}
	sort.Strings(klst)

	if doupdate {
		CurrentVendors, err := m1sql.GetAllVendors()
		if err != nil {
			c.Printf("Error: %v.\n", err)
			return
		}

		ncnt := 0
		for _, k := range klst {
			vinfo := vmap[k]
			bexists := false
			for _, cv := range CurrentVendors {
				if strings.TrimSpace(strings.ToLower(cv.FName)) == strings.TrimSpace(strings.ToLower(vinfo.FName)) {
					bexists = true
					break
				}
			}
			if bexists {
				continue
			}
			vnew := m1sql.Vendor{DName: vinfo.DName, FName: vinfo.FName, Aliases: make([]string, 0, len(vinfo.Aliases))}
			for _, a := range vinfo.Aliases {
				vnew.Aliases = append(vnew.Aliases, a)
			}
			err = m1sql.UpdateVendor(&vnew)
			if err != nil {
				c.Printf("Error updated vendor %s. Err=%v.\nAborting.\n", vinfo.FName, err)
				break
			}
			ncnt += 1
		}
		c.Printf("Number of new vendors added to the database: %d\n", ncnt)
		return
	}

	tbl := util.NewTable("Vendor", "Display Name", "Aliases")
	for _, k := range klst {
		ss := util.FormatStrSlice(vmap[k].Aliases)
		ss = util.FixStrLen(ss, 100, "...")
		tbl.AddRow(vmap[k].FName, vmap[k].DName, ss)
	}
	c.Printf("%s\n", tbl.Text())
	c.Printf("Number of Vendors: %d\n", len(vmap))
}
