package importer

import (
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/m1/types"
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

func ImportTransactions(filename string) (records []*types.ImportedTransaction, err error) {
	records = make([]*types.ImportedTransaction, 0, 10)
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %q.  (%v)", filename, err)
		return records, err
	}
	defer fi.Close()

	rdr := csv.NewReader(fi)
	rdr.LazyQuotes = true
	rdr.TrimLeadingSpace = true
	rdr.ReuseRecord = false
	ilinenum := 0
	r, err := rdr.Read()
	if err != nil {
		return records, fmt.Errorf("Header Record unreadable. %q", err)
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
		rd, err := convert_record(m, r, ilinenum)
		if err != nil {
			return records, err
		}
		records = append(records, rd)
	}
	return records, err
}

func convert_record(m headermap, r []string, ilinenum int) (rout *types.ImportedTransaction, err error) {
	if len(r) <= m.nMaxField {
		return nil, fmt.Errorf("Only %d columnes found in line %d. Need %d columns.", len(r), ilinenum, m.nMaxField)
	}
	rd := &types.ImportedTransaction{}
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
		rd.DatePosted, err = util.ParseGenericTime(r[m.iPostDate])
		if err != nil {
			return nil, fmt.Errorf("Unable to convert post date in line %d.", ilinenum)
		}
	}
	if m.iSettleDate >= 0 {
		rd.DateSettled, err = util.ParseGenericTime(r[m.iSettleDate])
		if err != nil {
			return nil, fmt.Errorf("Unable to convert settle date in line %d.", ilinenum)
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
			fmt.Printf("Blank value found for amount in line %d. Using zero.\n", ilinenum)
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
		return m, fmt.Errorf("A column for settle date is missing.")
	}
	m.nMaxField = findmaxfiled(m.iAccount, m.iMonth, m.iPostDate, m.iSettleDate, m.iBankInfo, m.iLocation,
		m.iVendor, m.iDescription, m.iCategory, m.iFlag, m.iAmount, m.iReceipt)

	return m, nil
}

func findmaxfiled(ifields ...int) int {
	n := -1
	for i := 0; i < len(ifields); i++ {
		if ifields[i] > n {
			n = ifields[i]
		}
	}
	return n
}

// remove_commas delete any commas found in the input string.
func remove_commas(s string) string {
	n := len(s)
	sout := ""
	for i := 0; i < n; i++ {
		c := s[i]
		if c != ',' {
			sout += string(c)
		}
	}
	return sout
}
