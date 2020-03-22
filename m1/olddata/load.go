// --------------------------------------------------------------------
// load.go -- Reads old data from csv files and loads it into the
// database.
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
	"strconv"
	"strings"
)

type headermap struct {
	iAccount     int
	iMonth       int
	iPostDate    int
	iSettleDate  int
	iBankInfo    int
	iLocation    int
	iVendor      int
	iDescription int
	iCategory    int
	iFlag        int
	iAmount      int
	iReceipt     int
	nMaxField    int
}

func LoadTransactions(c *util.Context, filename string, records []*OldTransaction) ([]*OldTransaction, error) {
	//records = make([]*OldTransaction, 0, 10)
	var err error
	fi, err := os.Open(filename)
	if err != nil {
		c.Printf("Error opening %q. Err=%v\n", filename, err)
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
	m, err := map_header(r)
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
		rd, err := convert_record(c, m, r, ilinenum)
		if err != nil {
			return records, err
		}
		records = append(records, rd)
	}
	return records, nil
}

func convert_record(c *util.Context, m headermap, r []string, ilinenum int) (rout *OldTransaction, err error) {
	if len(r) <= m.nMaxField {
		return nil, fmt.Errorf("Only %d columnes found in line %d. Need %d columns.", len(r), ilinenum, m.nMaxField)
	}
	rd := &OldTransaction{}
	if m.iAccount >= 0 {
		rd.Account = r[m.iAccount]
	}
	if m.iMonth >= 0 {
		rd.Month, err = util.ParseGenericTime(r[m.iMonth])
		if err != nil {
			return nil, fmt.Errorf("Unable to convert statement month in line %d. (%q)", ilinenum, r[m.iMonth])
		}
	}
	if m.iPostDate >= 0 {
		s := r[m.iPostDate]
		if !util.Blank(s) {
			rd.DatePosted, err = util.ParseGenericTime(s)
			if err != nil {
				return nil, fmt.Errorf("Unable to convert post date in line %d. (%q)", ilinenum, s)
			}
		}
	}
	if m.iSettleDate >= 0 {
		s := r[m.iSettleDate]
		if !util.Blank(s) {
			rd.DateSettled, err = util.ParseGenericTime(s)
			if err != nil {
				return nil, fmt.Errorf("Unable to convert settle date in line %d. (%s)", ilinenum, s)
			}
		}
	}
	if m.iBankInfo >= 0 {
		rd.BankInfo = r[m.iBankInfo]
	}
	if m.iLocation >= 0 {
		rd.Location = r[m.iLocation]
	}
	if m.iVendor >= 0 {
		rd.Vendor = r[m.iVendor]
	}
	if m.iDescription >= 0 {
		rd.Description = r[m.iDescription]
	}
	if m.iCategory >= 0 {
		rd.Category = r[m.iCategory]
	}
	if m.iFlag >= 0 {
		rd.Flag = r[m.iFlag]
	}
	if m.iReceipt >= 0 {
		rd.Receipt = r[m.iReceipt]
	}
	if m.iAmount > 0 {
		t := remove_commas(r[m.iAmount])
		t = strings.TrimSpace(t)
		if len(t) <= 0 {
			c.Printf("Blank value found for amount in line %d. Using zero.\n", ilinenum)
			rd.Amount = 0
		} else {
			f, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse the amount in line %d. (%q)", ilinenum, r[m.iAmount])
			}
			rd.Amount = int(f * 100.0)
		}
	}
	return rd, nil
}

func map_header(r []string) (m headermap, err error) {
	//No.,Account,Stmnt,Trans Date,Settle Date,Bank Info,Location,Vendor,Description,Cat,Flag,Amount,,Receipt,
	m = headermap{iAccount: -1, iMonth: -1, iPostDate: -1, iSettleDate: -1, iBankInfo: -1, iLocation: -1,
		iVendor: -1, iDescription: -1, iCategory: -1, iFlag: -1, iAmount: -1, iReceipt: -1}
	for i := 0; i < len(r); i++ {
		col := strings.ToLower(r[i])
		if col == "account" || col == "acc" {
			m.iAccount = i
		}
		if col == "stmnt" || col == "statement" || col == "month" {
			m.iMonth = i
		}
		if col == "post" || col == "postdate" || col == "post date" || col == "trans date" || col == "transdate" {
			m.iPostDate = i
		}
		if col == "settle" || col == "settle date" || col == "settledate" {
			m.iSettleDate = i
		}
		if col == "bank" || col == "bankinfo" || col == "bank info" || col == "bank infomation" {
			m.iBankInfo = i
		}
		if col == "loc" || col == "location" {
			m.iLocation = i
		}
		if col == "vendor" {
			m.iVendor = i
		}
		if col == "desc" || col == "description" || col == "comment" || col == "notes" {
			m.iDescription = i
		}
		if col == "cat" || col == "category" {
			m.iCategory = i
		}
		if col == "flag" {
			m.iFlag = i
		}
		if col == "amount" {
			m.iAmount = i
		}
		if col == "receipt" {
			m.iReceipt = i
		}
	}
	if m.iAmount < 0 {
		return m, fmt.Errorf("A column for Amount is missing.")
	}
	if m.iSettleDate < 0 {
		return m, fmt.Errorf("A column for Settle Date is missing.")
	}
	m.nMaxField = findmax(m.iAccount, m.iMonth, m.iPostDate, m.iSettleDate, m.iBankInfo, m.iLocation,
		m.iVendor, m.iDescription, m.iCategory, m.iFlag, m.iAmount, m.iReceipt)

	return m, nil
}
