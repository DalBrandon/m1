// --------------------------------------------------------------------
// vendors.go -- Loads all old data about vendors from csv files.
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

package olddata

import (
	"dbe/lib/util"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var vendor_file string = "Vendors.csv"

type VendorItem struct {
	Alias string
	FName string
	DName string
}

type vendorheadermap struct {
	iAlias    int
	iFName    int
	iDName    int
	nMaxField int
}

func GetVendorData(c *util.Context) ([]*VendorItem, error) {
	records := make([]*VendorItem, 0, 2500)
	var err error
	fn := olddata_folder + vendor_file
	fi, err := os.Open(fn)
	if err != nil {
		c.Printf("Error opening %q. Err=%v\n", fn, err)
		return records, err
	}
	defer fi.Close()

	rdr := csv.NewReader(fi)
	rdr.LazyQuotes = true
	rdr.TrimLeadingSpace = true
	rdr.ReuseRecord = false
	ilinenum := 1
	r, err := rdr.Read()
	if err != nil {
		c.Printf("Header record unreadable.\n")
		return records, fmt.Errorf("Header record unreadable. %q", err)
	}
	m, err := map_vendor_header(r)
	if err != nil {
		return records, err
	}

	for {
		ilinenum++
		r, err := rdr.Read()
		if err == io.EOF {
			err = nil
			break
		}
		// Clean all data
		for ii := 0; ii < len(r); ii++ {
			r[ii] = util.CleanForCafe(r[ii])
		}
		// If all fields are blank, continue
		if all_fields_blank(r) {
			continue
		}
		rd, err := convert_vendor_record(c, m, r, ilinenum)
		if err != nil {
			return records, err
		}
		records = append(records, rd)
	}
	return records, nil
}

func convert_vendor_record(c *util.Context, m vendorheadermap, r []string, ilinenum int) (rout *VendorItem, err error) {
	if len(r) <= m.nMaxField {
		return nil, fmt.Errorf("Only %d columnes found in line %d. Need %d columns.", len(r), ilinenum, m.nMaxField)
	}
	rd := &VendorItem{}
	rd.Alias = r[m.iAlias]
	rd.FName = r[m.iFName]
	rd.DName = r[m.iDName]
	return rd, nil
}

func map_vendor_header(r []string) (m vendorheadermap, err error) {
	//Alias,Vendor Name,Display Name
	m = vendorheadermap{iAlias: -1, iFName: -1, iDName: -1}
	for i := 0; i < len(r); i++ {
		col := strings.TrimSpace(strings.ToLower(r[i]))
		col = util.CleanStr(col, "")
		//fmt.Printf("col = %q\n", col)
		if col == "alias" {
			m.iAlias = i
		}
		if col == "name" || col == "vendor" || col == "vendor name" {
			m.iFName = i
		}
		if col == "display" || col == "display name" {
			m.iDName = i
		}
	}
	if m.iAlias < 0 || m.iFName < 0 || m.iDName < 0 {
		return m, fmt.Errorf("A requried column is missing (%d, %d, %d)", m.iAlias, m.iFName, m.iDName)
	}
	m.nMaxField = findmax(m.iAlias, m.iFName, m.iDName)

	return m, nil
}
