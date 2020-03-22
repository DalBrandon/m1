// --------------------------------------------------------------------
// types.go -- Types for the m1 project
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package m1data

import (
	"dbe/lib/uuid"
	"time"
)

// Database is the entire database for the m1 project
type Database struct {
	Accounts     map[string]*Account
	Vendors      map[string]*Vendor
	Categories   map[string]*Category
	Transactions map[uuid.UUID]*Transaction
}

// Transaction is the basic data item for m1
type Transaction struct {
	Tid         uuid.UUID // To uniquely id this transaction
	Amount      int       // In cents
	Account     string    // Points to FName in Accounts map
	Vendor      string    // Points to FName in Vendors map
	Cats        []CatItem // Can be empty but not nil
	Description string
	DatePosted  time.Time
	DateSettled time.Time
	Month       time.Time // Statement Month-Year
	BankInfo    string
	Location    string
	CheckNum    string
	Flag        string
	Receipts    []string // Urls to receipt files (images, pdfs, etc.)
	Notes       string
}

// CatItem is use to categorize transactions.  Note that
// a transaction should have at least one CatItem and all the
// CatItems in a transaction should add to the ammount in
// the transaction.
type CatItem struct {
	Amount   int
	Category string // Points to name in Category map
	Notes    string
}

// Vendor describes the primary party for a transaction.
type Vendor struct {
	FName          string
	DName          string
	Aliases        []string
	PrimaryProduct string
	BusinessType   string
	DefaultCat     string // blank, or points to name in Category map
	Notes          string
}

// Category is used to organize transactions.
type Category struct {
	Name    string
	Aliases []string // Must contain the Name.
	Notes   string
}

// Account is the basic bucket where money flows in or out.
type Account struct {
	ShortName string
	DName     string
	FName     string
	Notes     string
	Active    bool
	Aliases   []string
}

// Year returns the year (as a 4 digit int) in which the
// transaction belongs for booking -- which is normally
// the year that the transaction took place.
func (t *Transaction) Year() int {
	d := t.Date()
	return d.Year()
}

// Date returns the date to use for the transaction.
func (t *Transaction) Date() time.Time {
	if !t.DatePosted.IsZero() {
		return t.DatePosted
	}
	if !t.DateSettled.IsZero() {
		return t.DateSettled
	}
	var d time.Time
	return d
}

// HasDate returns true if the transaction has a date.
func (t *Transaction) HasDate() bool {
	if t.DateSettled.IsZero() && t.DatePosted.IsZero() {
		return false
	}
	return true
}
