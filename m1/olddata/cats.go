// --------------------------------------------------------------------
// cats.go -- Loads all old data about cats from csv files.
//
// Created 2020-03-20 DLB
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

var cat_file string = "/home/dal/dev/src/dbe/m1/olddata/Categories.csv"

type CatItem struct {
	Alias string
	Name  string
}

type catheadermap struct {
	iAlias    int
	iName     int
	nMaxField int
}

func GetCatData(c *util.Context) ([]*CatItem, error) {
	records := make([]*CatItem, 0, 1000)
	var err error
	fi, err := os.Open(cat_file)
	if err != nil {
		c.Printf("Error opening %q. Err=%v\n", cat_file, err)
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
	m, err := map_cat_header(r)
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
			r[ii] = strings.TrimSpace(util.CleanForCafe(r[ii]))
		}
		// If all fields are blank, continue
		if all_fields_blank(r) {
			continue
		}
		rd, err := convert_cat_record(c, m, r, ilinenum)
		if err != nil {
			return records, err
		}
		records = append(records, rd)
	}
	return records, nil
}

func convert_cat_record(c *util.Context, m catheadermap, r []string, ilinenum int) (rout *CatItem, err error) {
	if len(r) <= m.nMaxField {
		return nil, fmt.Errorf("Only %d columnes found in line %d. Need %d columns.", len(r), ilinenum, m.nMaxField)
	}
	rd := &CatItem{}
	rd.Alias = r[m.iAlias]
	rd.Name = r[m.iName]
	return rd, nil
}

func map_cat_header(r []string) (m catheadermap, err error) {
	m = catheadermap{iAlias: -1, iName: -1}
	for i := 0; i < len(r); i++ {
		col := strings.TrimSpace(strings.ToLower(r[i]))
		col = util.CleanStr(col, "")
		if col == "alias" {
			m.iAlias = i
		}
		if col == "name" {
			m.iName = i
		}
	}
	if m.iAlias < 0 || m.iName < 0 {
		return m, fmt.Errorf("A requried column is missing (%d, %d)", m.iAlias, m.iName)
	}
	m.nMaxField = findmax(m.iAlias, m.iName)

	return m, nil
}
