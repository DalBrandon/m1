// --------------------------------------------------------------------
// manager.go -- Manages the data
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package m1data

import (
	"dbe/lib/util"
	"dbe/lib/uuid"
	"fmt"
	"sync"
)

var dblock sync.Mutex
var db *Database

func init() {
	if db == nil {
		db = &Database{}
		db.Accounts = make(map[string]*Account, 10)
		db.Vendors = make(map[string]*Vendor, 4000)
		db.Categories = make(map[string]*Category, 1000)
		db.Transactions = make(map[uuid.UUID]*Transaction, 30000)
	}
}

// GetVendors returns all the vendors in the database.
func GetVendors() []*Vendor {
	dblock.Lock()
	defer dblock.Unlock()
	copylst := make([]*Vendor, 0, len(db.Vendors))
	for _, v := range db.Vendors {
		vc := *v
		copylst = append(copylst, &vc)
	}
	return copylst
}

// GetVendor returns a vendor given is full name.
func GetVendor(name string) *Vendor {
	dblock.Lock()
	defer dblock.Unlock()
	vv, ok := db.Vendors[name]
	if !ok {
		return nil
	}
	return vv
}

// GetAccounts returns all the accounts in the database.
func GetAccounts() []*Account {
	dblock.Lock()
	defer dblock.Unlock()
	copylst := make([]*Account, 0, len(db.Accounts))
	for _, v := range db.Accounts {
		vc := *v
		copylst = append(copylst, &vc)
	}
	return copylst
}

// GetCategories returns all the Categories in the database.
func GetCategories() []*Category {
	dblock.Lock()
	defer dblock.Unlock()
	copylst := make([]*Category, 0, len(db.Categories))
	for _, v := range db.Categories {
		vc := *v
		copylst = append(copylst, &vc)
	}
	return copylst
}

// GetCategory returns the category given its name.
func GetCategory(name string) *Category {
	dblock.Lock()
	defer dblock.Unlock()
	cc, ok := db.Categories[name]
	if !ok {
		return nil
	}
	return cc
}

// GetTransactions returns all the transactions in the database.
func GetTransactions() []*Transaction {
	dblock.Lock()
	defer dblock.Unlock()
	copylst := make([]*Transaction, 0, len(db.Transactions))
	for _, v := range db.Transactions {
		vc := *v
		copylst = append(copylst, &vc)
	}
	return copylst
}

// AddVendor will either add a new vendor to the vendor list, or
// replace an existing vendor with a updated version.
func AddVendor(v *Vendor) error {
	if util.Blank(v.FName) {
		return fmt.Errorf("Vendor FName cannot be blank.")
	}
	dblock.Lock()
	defer dblock.Unlock()
	db.Vendors[v.FName] = v
	return nil
}

// AddCategory will add a category to the category list.  If the
// category already exists, it will be updated.
func AddCategory(c *Category) error {
	dblock.Lock()
	defer dblock.Unlock()
	for k, v := range db.Categories {
		if k == c.Name {
			continue
		}
		for _, n := range v.Aliases {
			for _, n2 := range c.Aliases {
				if n == n2 {
					return fmt.Errorf("Dublicate Aliases (%s) found in existing cat (%s).", n, k)
				}
			}
		}
	}
	if c.Aliases == nil {
		c.Aliases = make([]string, 0, 1)
	}
	// Add Name to Aliases if not already there.
	if !util.InStringSlice(c.Aliases, c.Name) {
		c.Aliases = append(c.Aliases, c.Name)
	}
	db.Categories[c.Name] = c
	return nil
}

// AddAccount will add an account to the accout list, or it will
// update an existing account
func AddAccount(a *Account) error {
	if util.Blank(a.FName) {
		return fmt.Errorf("Account FName cannot be blank.")
	}
	dblock.Lock()
	defer dblock.Unlock()
	db.Accounts[a.FName] = a
	return nil
}

// AddTransaction will either add a new transaction or update an
// existing transaction.  For updating, the Vid must be provided,
// or a new transaction is assumed.
func AddTransaction(t *Transaction) error {
	dblock.Lock()
	defer dblock.Unlock()
	// Make a copy...
	tc := *t
	_, ok := db.Accounts[tc.Account]
	if !ok {
		return fmt.Errorf("No Account (%s) for transaction.  Add Account first.", tc.Account)
	}
	if !util.Blank(tc.Vendor) {
		_, ok = db.Vendors[tc.Vendor]
		if !ok {
			return fmt.Errorf("No Vendor (%s) for transaction.  Add Vendor first.", tc.Vendor)
		}
	}
	if tc.Cats == nil {
		tc.Cats = make([]CatItem, 0, 1)
	}
	if tc.Tid.IsZero() {
		tc.Tid = uuid.New()
	}
	db.Transactions[tc.Tid] = &tc
	return nil
}
