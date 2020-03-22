// --------------------------------------------------------------------
// types.go -- Types for the olddata package
//
// Created 2020-03-16 DLB
// --------------------------------------------------------------------

package olddata

import (
	"time"
)

type OldTransaction struct {
	Amount      int // In cents
	Description string
	DatePosted  time.Time
	DateSettled time.Time
	Account     string    // Required
	Vendor      string    // blnak = unknown
	Month       time.Time // Statement Month
	BankInfo    string
	Location    string
	CheckNum    string
	Flag        string
	Category    string
	Receipt     string
}

func (t *OldTransaction) Year() int {
	d := t.Date()
	return d.Year()
}

func (t *OldTransaction) Date() time.Time {
	if !t.DatePosted.IsZero() {
		return t.DatePosted
	}
	if !t.DateSettled.IsZero() {
		return t.DateSettled
	}
	var d time.Time
	return d
}

func (t *OldTransaction) HasDate() bool {
	if t.DateSettled.IsZero() && t.DatePosted.IsZero() {
		return false
	}
	return true
}
