// --------------------------------------------------------------------
// accounts.go -- Manage the account and account alias tables
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package m1sql

import (
	"dbe/lib/log"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Alias struct {
	Name  string
	Notes string
}

type Account struct {
	Aid       int
	ShortName string
	DName     string
	FName     string
	Notes     string
	Active    bool
	Aliases   []Alias
}

type AccountAlias struct {
	Aid   int
	Alias string
	Notes string
}

var gAccountCasheLock sync.Mutex
var gAccountAliasCasheLock sync.Mutex

var gAccountCache []*Account
var gAccountAliases []*AccountAlias

func GetAccounts() []*Account {
	gAccountCasheLock.Lock()
	defer gAccountCasheLock.Unlock()
	if gAccountCache != nil {
		return gAccountCache
	}
	lstmap := make(map[int]*Account, 10)
	rows, err := m_db.Query("Select Aid, ShortName, DName, FName, Notes, Active from Accounts")
	if err != nil {
		log.Errorf("Err getting Accounts. Returning empty slice. Err=%v", err)
		return []*Account{}
	}
	for rows.Next() {
		var sname, dname, fname, notes string
		var aid, active int
		err = rows.Scan(&aid, &sname, &dname, &fname, &notes, &active)
		if err != nil {
			log.Errorf("Err during row scan in getAccounts. Shipping account. Err=%v.", err)
			continue
		}
		bactive := true
		if active == 0 {
			bactive = false
		}
		atmp := make([]Alias, 0, 10)
		lstmap[aid] = &Account{Aid: aid, DName: dname, ShortName: sname,
			FName: fname, Notes: notes, Active: bactive, Aliases: atmp}
	}
	aliases := GetAccountAliases()
	for _, a := range aliases {
		acc, ok := lstmap[a.Aid]
		if !ok {
			log.Errorf("Bad account alias. Aid=%d. Ignored.", a.Aid)
			continue
		}
		acc.Aliases = append(acc.Aliases, Alias{Name: a.Alias, Notes: a.Notes})
	}
	lst := make([]*Account, 0, len(lstmap))
	for _, v := range lstmap {
		lst = append(lst, v)
	}
	sort.Slice(lst, func(i, j int) bool { return lst[j].Aid > lst[i].Aid })
	gAccountCache = lst
	return gAccountCache
}

func GetAccountAliases() []*AccountAlias {
	gAccountAliasCasheLock.Lock()
	defer gAccountAliasCasheLock.Unlock()
	if gAccountAliases != nil {
		return gAccountAliases
	}
	lst := make([]*AccountAlias, 0, 10)
	rows, err := m_db.Query("Select Aid, Alias, Notes from AccountAlias")
	if err != nil {
		log.Errorf("Err getting AccountAlias. Returning empty slice. Err=%v", err)
		return lst
	}
	for rows.Next() {
		var aid int
		var alias, notes string
		err = rows.Scan(&aid, &alias, &notes)
		if err != nil {
			log.Errorf("Err during row scan in AccountAlias. Shipping alias. Err=%v.", err)
			continue
		}
		lst = append(lst, &AccountAlias{Aid: aid, Alias: alias, Notes: notes})
	}
	gAccountAliases = lst
	return gAccountAliases
}

func AddAccountAlias(Aid int, Name string, Notes string) error {
	var err error
	acclst := GetAccounts()
	bFound := false
	for _, a := range acclst {
		if a.Aid == Aid {
			bFound = true
			break
		}
	}
	if !bFound {
		return fmt.Errorf("Invalid account ID.")
	}
	aliaslst := GetAccountAliases()
	for _, a := range aliaslst {
		if a.Aid == Aid && strings.ToLower(strings.TrimSpace(a.Alias)) == strings.ToLower(strings.TrimSpace(Name)) {
			return fmt.Errorf("Alias already exists.")
		}
	}
	stmt, err := m_db.Prepare("insert AccountAlias set Aid=?, Alias=?, Notes=?")
	if err != nil {
		return fmt.Errorf("Err inserting into AccountAlias. Err=%v", err)
	}
	r, err := stmt.Exec(Aid, Name, Notes)
	if err != nil {
		return fmt.Errorf("Err inserting into AccountAlias. Err=%v", err)
	}
	if n, _ := r.RowsAffected(); n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when adding an alias to an account.", n)
	}

	gAccountAliasCasheLock.Lock()
	defer gAccountAliasCasheLock.Unlock()
	gAccountCasheLock.Lock()
	defer gAccountCasheLock.Unlock()
	gAccountCache = nil
	gAccountAliases = nil
	return err
}

func DeleteAccountAlias(Aid int, Name string) (int, error) {
	stmt, err := m_db.Prepare("delete from AccountAlias where Aid=? and Alias=?")
	if err != nil {
		return 0, fmt.Errorf("Err deleting from AccountAlias. Err=%v", err)
	}
	r, err := stmt.Exec(Aid, Name)
	if err != nil {
		return 0, fmt.Errorf("Err deleting from AccountAlias. Err=%v", err)
	}
	nn, _ := r.RowsAffected()
	n := int(nn)
	if n != 0 && n != 1 {
		err = fmt.Errorf("Wrong number of rows affected (%d) when deleting an alias from an account.", n)
	}
	gAccountAliasCasheLock.Lock()
	defer gAccountAliasCasheLock.Unlock()
	gAccountCasheLock.Lock()
	defer gAccountCasheLock.Unlock()
	gAccountCache = nil
	gAccountAliases = nil
	return n, err
}
