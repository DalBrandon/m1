package types

import (
	"dbe/lib/uuid"
	"time"
)

type Transaction struct {
	TID         uuid.UUID
	Amount      float64 // In dollars, rounded to nearest cent.
	Description string
	DatePosted  time.Time
	DateSettled time.Time
	Account     uuid.UUID // Required
	Vendor      uuid.UUID // 0  = unknown
	MonthNum    int       // Statement Month, 1-12 */
	BankInfo    string
	Location    string
	CheckNum    string
}
