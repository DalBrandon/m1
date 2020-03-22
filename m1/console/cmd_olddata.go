// --------------------------------------------------------------------
// cmd_olddata.go -- Analyzes Old Data
//
// Created 2020-03-17 DLB
// --------------------------------------------------------------------

package console

import (
	"dbe/lib/util"
	//"dbe/m1/m1sql"
	m1 "dbe/m1/m1data"
	"dbe/m1/olddata"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var gTopic_olddata string = `
The olddata command is used to analyze old data.  The format of
the command is:

  olddata report yyyy  

where report is the type of report desired and yyyy is the year
that is covered by the report.  If yyyy is omited, all the old 
data is included in the report.  Possible reports are:

  cat      -- list the categories
  vendors  -- lists the vendors
  accounts -- lists the accounts
  ckeck    -- Checks for errors between old and new
  load     -- Loads the old data into the database
`

func init() {
	RegistorCmd("olddata", "", "Analyze Old Data.", handle_olddata)
	RegistorTopic("olddata", gTopic_olddata)
}

func handle_olddata(c *util.Context, cmdline string) {
	params := make(map[string]string, 10)
	args, err := ParseCmdLine(cmdline, params)
	if err != nil {
		c.Printf("%v\n", err)
		return
	}
	if len(args) < 2 {
		c.Printf("Not enough args.\n")
		return
	}
	rpt := strings.ToLower(args[1])
	year := 0
	if len(args) >= 3 {
		year, err = strconv.Atoi(args[2])
		if err != nil {
			c.Printf("Bad year input (%s).\n", args[2])
			return
		}
		if year < 2000 || year > 2030 {
			c.Printf("Year is out of range (2000-2030). \n")
			return
		}
	}
	if rpt == "cat" || rpt == "category" {
		oldata_cat_report(c, year)
		return
	}
	if rpt == "accounts" {
		olddata_accounts_report(c, year)
		return
	}
	if rpt == "vendors" {
		olddata_vendors_report(c, year)
		return
	}
	if rpt == "load" {
		olddata_load(c)
		return
	}
	// if rpt == "check" {
	// 	olddata_check_report(c, year)
	// 	return
	// }
	c.Printf("Error -- unknown operation (%s).\n", rpt)
}

func oldata_cat_report(c *util.Context, year int) {
	type catinfo struct {
		Name   string
		Amount int
		Count  int
	}
	newcontext := util.NewContext(util.Context_Internal)
	t0 := time.Now()
	transactions, err := olddata.ReadAll(newcontext)
	tload := time.Now().Sub(t0)
	c.Printf("Time to Load Data: %d ms.\n", tload.Milliseconds())
	if err != nil {
		c.Printf("Error getting data: %q\n", err)
		return
	}
	if len(transactions) <= 0 {
		c.Printf("No transactions found.\n")
		return
	}
	ncnt := 0
	m := make(map[string]*catinfo, 50)
	for _, x := range transactions {
		if x.Year() == year || year == 0 {
			ncnt += 1
			cat := x.Category
			mcat, ok := m[cat]
			if !ok {
				m[cat] = &catinfo{Name: cat, Amount: x.Amount, Count: 1}
			} else {
				mcat.Count = mcat.Count + 1
				mcat.Amount = mcat.Amount + x.Amount
			}
		}
	}
	if year == 0 {
		c.Printf("Number of transactions (all years) = %d\n", ncnt)
	} else {
		c.Printf("Number of transactions (year=%d) = %d\n", year, ncnt)
	}
	// Sort
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tbl := util.NewTable("Category", "Item Count", "Total")
	for _, cat := range keys {
		scnt := fmt.Sprintf("%5d", m[cat].Count)
		samt := util.StrLeft(util.CentsToStr(m[cat].Amount), 14)
		tbl.AddRow(cat, scnt, samt)
	}
	c.Printf("%s", tbl.Text())
}

func olddata_accounts_report(c *util.Context, year int) {
	type accinfo struct {
		Name   string
		Amount int
		Count  int
	}
	newcontext := util.NewContext(util.Context_Internal)
	t0 := time.Now()
	transactions, err := olddata.ReadAll(newcontext)
	tload := time.Now().Sub(t0)
	c.Printf("Time to Load Data: %d ms.\n", tload.Milliseconds())
	if err != nil {
		c.Printf("Error getting data: %q\n", err)
		return
	}
	if len(transactions) <= 0 {
		c.Printf("No transactions found.\n")
		return
	}
	ncnt := 0
	m := make(map[string]*accinfo, 50)
	for _, x := range transactions {
		if x.Year() == year || year == 0 {
			ncnt += 1
			acc := x.Account
			macc, ok := m[acc]
			if !ok {
				m[acc] = &accinfo{Name: acc, Amount: x.Amount, Count: 1}
			} else {
				macc.Count = macc.Count + 1
				macc.Amount = macc.Amount + x.Amount
			}
		}
	}
	if year == 0 {
		c.Printf("Number of transactions (all years) = %d\n", ncnt)
	} else {
		c.Printf("Number of transactions (year=%d) = %d\n", year, ncnt)
	}
	// Sort
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tbl := util.NewTable("Account", "Item Count", "Total")
	for _, acc := range keys {
		scnt := fmt.Sprintf("%5d", m[acc].Count)
		samt := util.StrLeft(util.CentsToStr(m[acc].Amount), 14)
		tbl.AddRow(acc, scnt, samt)
	}
	c.Printf("%s", tbl.Text())
}

func olddata_vendors_report(c *util.Context, year int) {
	type vendorinfo struct {
		Name   string
		Amount int
		Count  int
	}
	newcontext := util.NewContext(util.Context_Internal)
	t0 := time.Now()
	transactions, err := olddata.ReadAll(newcontext)
	tload := time.Now().Sub(t0)
	c.Printf("Time to Load Data: %d ms.\n", tload.Milliseconds())
	if err != nil {
		c.Printf("Error getting data: %q\n", err)
		return
	}
	if len(transactions) <= 0 {
		c.Printf("No transactions found.\n")
		return
	}
	ncnt := 0
	m := make(map[string]*vendorinfo, 50)
	for _, x := range transactions {
		if x.Year() == year || year == 0 {
			ncnt += 1
			vendor := x.Vendor
			mvend, ok := m[vendor]
			if !ok {
				m[vendor] = &vendorinfo{Name: vendor, Amount: x.Amount, Count: 1}
			} else {
				mvend.Count = mvend.Count + 1
				mvend.Amount = mvend.Amount + x.Amount
			}
		}
	}
	if year == 0 {
		c.Printf("Number of transactions (all years) = %d\n", ncnt)
	} else {
		c.Printf("Number of transactions (year=%d) = %d\n", year, ncnt)
	}
	// Sort
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tbl := util.NewTable("Account", "Item Count", "Total")
	for _, vendor := range keys {
		scnt := fmt.Sprintf("%5d", m[vendor].Count)
		samt := util.StrLeft(util.CentsToStr(m[vendor].Amount), 14)
		tbl.AddRow(vendor, scnt, samt)
	}
	c.Printf("%s", tbl.Text())
}

// func olddata_check_report(c *util.Context, year int) {
// 	type accinfo struct {
// 		Name   string
// 		Amount int
// 		Count  int
// 	}
// 	newcontext := util.NewContext(util.Context_Internal)
// 	t0 := time.Now()
// 	transactions, err := olddata.ReadAll(newcontext)
// 	tload := time.Now().Sub(t0)
// 	c.Printf("Time to Load Data: %d ms.\n", tload.Milliseconds())
// 	if err != nil {
// 		c.Printf("Error getting data: %q\n", err)
// 		return
// 	}
// 	if len(transactions) <= 0 {
// 		c.Printf("No transactions found.\n")
// 		return
// 	}
// 	// Check that each transaction has one and only one account match
// 	newaccounts := m1sql.GetAccounts()
// 	nerrcnt := 0
// 	ncheckcnt := 0
// 	for _, t := range transactions {
// 		nfind := 0
// 		taccount := strings.TrimSpace(strings.ToLower(t.Account))
// 		ncheckcnt += 1
// 		for _, a := range newaccounts {
// 			found := false
// 			if strings.TrimSpace(strings.ToLower(a.ShortName)) == taccount {
// 				found = true
// 			}
// 			if strings.TrimSpace(strings.ToLower(a.DName)) == taccount {
// 				found = true
// 			}
// 			if strings.TrimSpace(strings.ToLower(a.FName)) == taccount {
// 				found = true
// 			}
// 			for _, aa := range a.Aliases {
// 				if strings.TrimSpace(strings.ToLower(aa.Name)) == taccount {
// 					found = true
// 				}
// 			}
// 			if found {
// 				nfind += 1
// 			}
// 		}
// 		if nfind != 1 {
// 			if nfind == 0 {
// 				c.Printf("\nNo Account Found for Transaction:\n ")
// 			} else {
// 				c.Printf("\nMore than one account found for Transaction (%d):\n", nfind)
// 			}
// 			c.Printf("Date: %s, Account: %s, Ammount: %s\n", t.Date().Format("2006-01-02"),
// 				t.Account, util.CentsToStr(t.Amount))
// 			c.Printf("Cat: %s, Description: %s\n", t.Category, t.Description)
// 			nerrcnt += 1
// 		}
// 	}
// 	c.Printf("Number of transactions checked: %d\n", ncheckcnt)
// 	if nerrcnt == 0 {
// 		c.Printf("No errors found.\n")
// 	} else {
// 		c.Printf("Number of transcations with invalid accounts: %d\n", nerrcnt)
// 	}
// }

func olddata_load(c *util.Context) {
	newcontext := util.NewContext(util.Context_Internal)
	catdata, err := olddata.GetCatData(newcontext)
	if err != nil {
		c.Printf("Error reading cat data. Err=%v. Aborting.\n")
		return
	}
	vendata, err := olddata.GetVendorData(newcontext)
	if err != nil {
		c.Printf("Error reading vendor data. Err=%v. Aborting.\n")
		return
	}

	// ---------------   Accounts
	if len(m1.GetAccounts()) <= 0 {
		err1 := m1.AddAccount(&m1.Account{ShortName: "ML", DName: "ML Checking", Active: false,
			FName: "Merrill Lynch Checking", Aliases: []string{"ML Check"}})
		err2 := m1.AddAccount(&m1.Account{ShortName: "FMB", DName: "FM Checking", Active: true,
			FName: "Farmers and Merchants Bank", Aliases: []string{"FMB Check", "FM Check", "FRB Check"}})
		err3 := m1.AddAccount(&m1.Account{ShortName: "BofA", DName: "BofA Checking", Active: true,
			FName: "Bank of America Checking", Aliases: []string{"BofA Check"}})
		err4 := m1.AddAccount(&m1.Account{ShortName: "DVisa", DName: "Dal's Visa", Active: true,
			FName: "Dal's Visa at ML", Aliases: []string{"Visa 3232", "Visa 3627",
				"Visa 6974", "Visa 7759", "Visa 9389"}})
		err5 := m1.AddAccount(&m1.Account{ShortName: "CVisa", DName: "Carol's Visa", Active: true,
			FName: "Carols's Visa at ML", Aliases: []string{"Visa 1638", "Visa 2859",
				"Visa 3799", "Visa 4007", "Visa 4907", "Visa 9713", "Visa 9836"}})
		err6 := m1.AddAccount(&m1.Account{ShortName: "HDept", DName: "Home Dept", Active: true,
			FName: "Home Dept", Aliases: []string{}})
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
			c.Printf("Error loading accounts.\n")
			return
		}
		c.Printf("Accounts Loaded.\n")
		c.Flush()
	} else {
		c.Printf("Skipping loadding of accounts. They seem to be there.\n")
		c.Flush()
	}

	// ---------------   Categories
	mcats := make(map[string]*m1.Category, len(catdata))
	for _, cx := range catdata {
		if util.Blank(cx.Name) {
			continue
		}
		mc, ok := mcats[cx.Name]
		if !ok {
			cnew := &m1.Category{Name: cx.Name, Notes: "", Aliases: make([]string, 0, 5)}
			cnew.Aliases = append(cnew.Aliases, cx.Alias)
			mcats[cnew.Name] = cnew
		} else {
			if !util.InStringSlice(mc.Aliases, cx.Alias) {
				mc.Aliases = append(mc.Aliases, cx.Alias)
			}
		}
	}
	c.Printf("Number of old categories: %d\n", len(mcats))
	ncatcnt := 0
	for k, v := range mcats {
		// If the cat exists, don't add it.
		ce := m1.GetCategory(k)
		if ce == nil {
			err := m1.AddCategory(v)
			if err != nil {
				c.Printf("Unable to add cat %s. Err=%v. Aborting.\n", k, err)
				return
			}
			ncatcnt += 1
		}
	}
	c.Printf("Number of new categories added: %d\n", ncatcnt)
	c.Flush()

	// ---------------   Vendors
	mvens := make(map[string]*m1.Vendor, len(vendata))
	for _, vx := range vendata {
		if util.Blank(vx.FName) {
			continue
		}
		mv, ok := mvens[vx.FName]
		if !ok {
			vnew := &m1.Vendor{FName: vx.FName, DName: vx.DName, Aliases: make([]string, 0, 5)}
			vnew.Aliases = append(vnew.Aliases, vx.Alias)
			mvens[vnew.FName] = vnew
		} else {
			if !util.InStringSlice(mv.Aliases, vx.Alias) {
				mv.Aliases = append(mv.Aliases, vx.Alias)
			}
		}
	}
	c.Printf("Number of old vendors: %d\n", len(mvens))
	nvencnt := 0
	for k, v := range mvens {
		// If the vendor exists, don't add it.
		ve := m1.GetVendor(k)
		if ve == nil {
			err := m1.AddVendor(v)
			if err != nil {
				c.Printf("Unable to add vendor %s. Err=%v. Aborting.\n", k, err)
				return
			}
			nvencnt += 1
		}
	}
	c.Printf("Number of new vendors added: %d\n", nvencnt)
	c.Flush()

	// ------------- Transactions
	temp_accounts := m1.GetAccounts()
	c.Printf("Number of Accounts in Database: %d\n", len(temp_accounts))
	c.Flush()
	temp_vendors := m1.GetVendors()
	c.Printf("Number of Vendors in Database: %d\n", len(temp_vendors))
	c.Flush()
	temp_catetories := m1.GetCategories()
	c.Printf("Number of Categories in Database: %d\n", len(temp_catetories))
	c.Flush()
	temp_transactions := m1.GetTransactions()
	c.Printf("Number of Transactions in Database: %d\n", len(temp_transactions))
	c.Flush()

	tlst, err := olddata.ReadAll(newcontext)
	if err != nil {
		c.Printf("Error reading old tranactions. Err=%v\n", err)
		return
	}
	c.Printf("Number of Old Transactions to Process: %d\n", len(tlst))
	c.Flush()

	nSkips := 0
	nErrs := 0
	nCnt := 0
	t0 := time.Now()
	for _, t := range tlst {
		svendor, err1 := getbestvendor(temp_vendors, t.Vendor)
		saccount, err2 := getbestaccount(temp_accounts, t.Account)
		scat, err3 := getbestcategory(temp_catetories, t.Category)
		if err1 != nil || err2 != nil || err3 != nil {
			nSkips += 1
			if nSkips < 20 {
				c.Printf("Skipping Tranactions. %v, %v, %v\n", err1, err2, err3)
				c.Printf("Tranaction: Amount: %d, Date: %s, Account: %s, Description: %s\n",
					t.Amount, t.Date().Format("06-01-02"), t.Account, t.Description)
			}
			continue
		}
		bSkip := false
		for _, tt := range temp_transactions {
			if tt.Date() == t.Date() && t.Description == tt.Description &&
				tt.Account == saccount && tt.Amount == t.Amount {
				nSkips += 1
				if nSkips < 20 {
					c.Printf("Duplicate Transaction? Skipping...\n")
					c.Printf("Date: %s, Account: %s, Description: %s\n",
						t.Date().Format("06-01-02"), t.Account, t.Description)
				}
				bSkip = true
			}
		}
		if bSkip {
			continue
		}
		tnew := &m1.Transaction{}
		tnew.Amount = t.Amount
		tnew.DatePosted = t.DatePosted
		tnew.DateSettled = t.DateSettled
		tnew.Month = t.Month
		tnew.Description = t.Description
		tnew.Month = t.Month
		tnew.BankInfo = t.BankInfo
		tnew.Location = t.Location
		tnew.CheckNum = t.CheckNum
		tnew.Flag = t.Flag
		tnew.Account = saccount
		tnew.Vendor = svendor
		if util.Blank(scat) {
			tnew.Cats = []m1.CatItem{}
		} else {
			tnew.Cats = []m1.CatItem{m1.CatItem{Category: scat, Amount: t.Amount}}
		}
		tnew.Notes = fmt.Sprintf("From Old Data... Converted on %s\n")
		tnew.Notes += fmt.Sprintf("Old Account: %s\nOld Vendor: %s\nOld Cat: %s.\n",
			time.Now().Format("06-01-02"), t.Account, t.Vendor, t.Category)
		if !util.Blank(t.Flag) {
			tnew.Notes += fmt.Sprintf("Flags: %s\n", t.Flag)
		}
		if !util.Blank(t.CheckNum) {
			tnew.Notes += fmt.Sprintf("CheckNum: %s\n", t.CheckNum)
		}
		if !util.Blank(t.Receipt) {
			tnew.Notes += fmt.Sprintf("Receipt: %s\n", t.Receipt)
		}
		err := m1.AddTransaction(tnew)
		if err != nil {
			nErrs += 1
			if nErrs < 20 {
				c.Printf("Unable to add transaction. Err=%v\nTrans = %v\n", err, tnew)
			}
		} else {
			nCnt += 1
		}
		if time.Now().Sub(t0).Seconds() > 10.0 {
			c.Printf("Number of Transactions Processed so far = %d\n", nCnt)
			c.Flush()
			t0 = time.Now()
		}
	}
	if nSkips > 0 {
		c.Printf("Number of transactions skipped: %d (Only frist 20 shown).\n", nSkips)
	}
	if nErrs > 0 {
		c.Printf("Number of Errors: %d (Only first 20 shown).\n", nErrs)
	}
	c.Printf("Number of transcations added: %d\n", nCnt)
	c.Printf("Success.\n")
}

func getbestvendor(vendors []*m1.Vendor, v string) (string, error) {
	if util.Blank(v) {
		return "", nil
	}
	vv := strings.TrimSpace(v)
	for _, v := range vendors {
		if vv == strings.TrimSpace(v.FName) {
			return v.FName, nil
		}
		if vv == strings.TrimSpace(v.DName) {
			return v.FName, nil
		}
		for _, a := range v.Aliases {
			if vv == strings.TrimSpace(a) {
				return v.FName, nil
			}
		}
	}
	// Try lowercase.
	vv = strings.ToLower(v)
	for _, v := range vendors {
		if vv == strings.ToLower(strings.TrimSpace(v.FName)) {
			return v.FName, nil
		}
		if vv == strings.ToLower(strings.TrimSpace(v.DName)) {
			return v.FName, nil
		}
		for _, a := range v.Aliases {
			if vv == strings.ToLower(strings.TrimSpace(a)) {
				return v.FName, nil
			}
		}
	}
	return "", fmt.Errorf("No vendor for %s.", v)
}

func getbestaccount(accounts []*m1.Account, a string) (string, error) {
	aa := strings.ToLower(strings.TrimSpace(a))
	for _, account := range accounts {
		if aa == strings.ToLower(strings.TrimSpace(account.ShortName)) {
			return account.FName, nil
		}
		if aa == strings.ToLower(strings.TrimSpace(account.DName)) {
			return account.FName, nil
		}
		if aa == strings.ToLower(strings.TrimSpace(account.FName)) {
			return account.FName, nil
		}
		for _, alias := range account.Aliases {
			if aa == strings.ToLower(strings.TrimSpace(alias)) {
				return account.FName, nil
			}
		}
	}
	return "", fmt.Errorf("No account for %s.", a)
}

func getbestcategory(categories []*m1.Category, cat string) (string, error) {
	if util.Blank(cat) {
		return "", nil
	}
	catc := strings.ToLower(strings.TrimSpace(cat))
	for _, c := range categories {
		if catc == strings.ToLower(strings.TrimSpace(c.Name)) {
			return c.Name, nil
		}
		for _, a := range c.Aliases {
			if catc == strings.ToLower(strings.TrimSpace(a)) {
				return c.Name, nil
			}
		}
	}
	return "", fmt.Errorf("No category for %s.", cat)
}
