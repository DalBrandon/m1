// --------------------------------------------------------------------
// vendors.go -- Manage the vendors and vendor alias tables
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------

package m1sql

import (
	"database/sql"
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/lib/uuid"
	"fmt"
	"sort"
)

type Vendor struct {
	Vid            uuid.UUID
	DName          string
	FName          string
	DefaultCat     string // If populated, use it.
	BusinessType   string // "Retail", "Major Vendor", "Small Vendor", "Individual", "Misc"
	PrimaryProduct string
	Notes          string
	Aliases        []string
}

type VendorAlias struct {
	Vid   uuid.UUID
	Alias string
}

// GetAllVendors returns all vendors in the database
func GetAllVendors() ([]*Vendor, error) {
	lstmap := make(map[uuid.UUID]*Vendor, 2000)
	scmd := "Select Vendors.Vid, DName, FName, DefaultCat, BusinessType, PrimaryProduct, Notes, Alias from Vendors "
	scmd += "Left Join VendorAlias on Vendors.Vid=VendorAlias.Vid"
	rows, err := m_db.Query(scmd)
	if err != nil {
		log.Errorf("Err getting Vendors. Returning empty slice. Err=%v", err)
		return []*Vendor{}, fmt.Errorf("Err getting vendors. Err=%v", err)
	}
	defer rows.Close()
	for rows.Next() {
		vraw, err := scan_vendor(rows)
		if err != nil {
			log.Errorf("Bad row in Vendors: %v.", err)
			return []*Vendor{}, fmt.Errorf("Bad row in vendors: %v", err)
		}
		if vraw.Vid.IsZero() {
			log.Errorf("Vendor with zero uuid found. Skipping.", err)
			continue
		}
		v, ok := lstmap[vraw.Vid]
		if !ok {
			if len(vraw.Aliases) > 0 {
				alias := vraw.Aliases[0]
				vraw.Aliases = make([]string, 0, 5)
				vraw.Aliases = append(vraw.Aliases, alias)
			} else {
				vraw.Aliases = make([]string, 0, 5)
			}
			lstmap[vraw.Vid] = vraw
		} else {
			if len(vraw.Aliases) > 0 {
				alias := vraw.Aliases[0]
				if !util.Blank(alias) {
					v.Aliases = append(v.Aliases, alias)
				}
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Errorf("Database failure after iterating rows on Vendor table. Err=%v", err)
		return []*Vendor{}, fmt.Errorf("Error during row interation on Vendor Table. Err=%v", err)
	}
	lst := make([]*Vendor, 0, len(lstmap))
	for _, v := range lstmap {
		lst = append(lst, v)
	}
	sort.Slice(lst, func(i, j int) bool { return lst[j].FName > lst[i].FName })
	return lst, nil
}

// UpdateVendor will either change a vendor to match the input or make a new vendor
// if one doesn't already exist.
func UpdateVendor(v *Vendor) error {
	tx, err := m_db.Begin()
	if err != nil {
		log.Errorf("Unable to begin transaction for UpdateVendor. Err=%v", err)
		return err
	}
	if !v.Vid.IsZero() {
		// Does the vendor exist?
		vcurrent, err := get_vendor_by_uuid(tx, v.Vid)
		if err != nil {
			attempt_rollback(tx, "Err from get_vendor_by_uuid")
			return fmt.Errorf("Unable to get vendor by id. Err=%v", err)
		}
		if vcurrent != nil {
			// Vendor exists. Update it here.
			res, err := tx.Exec("Update Vendors set Dname=?, FName=?, DefaultCat=?, BusinessType=?, PrimaryProduct=?, Notes=? Where Vid=?",
				v.DName, v.FName, v.DefaultCat, v.BusinessType, v.PrimaryProduct, v.Notes, v.Vid.String())
			if err != nil {
				attempt_rollback(tx, "Update Vendors table failed.")
				return fmt.Errorf("Unable to update vendor table. Err=%v", err)
			}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				attempt_rollback(tx, "Getting RowsAffected failed after Update Vendors")
				return fmt.Errorf("Unable to retrieve rows affected after Update Vendors. Err=%v", err)
			}
			if rowCnt != 1 {
				attempt_rollback(tx, "Rowcount incorrect after Update Vendors.")
				return fmt.Errorf("Unable to update vendors. Row affected not one (was %d).", rowCnt)
			}
			err = delete_vendor_aliases(tx, vcurrent.Vid, len(vcurrent.Aliases))
			if err != nil {
				attempt_rollback(tx, "Unable to delete vendor aliases.")
				return fmt.Errorf("Unable to delete vendor aliases. Err=%v", err)
			}
			err = add_vendor_aliases(tx, v.Vid, v.Aliases)
			if err != nil {
				attempt_rollback(tx, "Unable to add vendor aliases.")
				return fmt.Errorf("Unable to add vendor aliases. Err=%v", err)
			}
			err = tx.Commit()
			if err != nil {
				log.Errorf("Commit failed on Update Vendors. Err=%v", err)
				return fmt.Errorf("Commit failed on Update Vendors. Err=%v", err)
			}
			return nil
		}
	}
	if v.Vid.IsZero() {
		v.Vid = uuid.New()
	}
	res, err := tx.Exec("Insert into Vendors(Vid, DName, FName, DefaultCat, BusinessType, PrimaryProduct, Notes)"+
		" values(?, ?, ?, ?, ?, ?, ?)",
		v.Vid.String(), v.DName, v.FName, v.DefaultCat, v.BusinessType, v.PrimaryProduct, v.Notes)
	if err != nil {
		attempt_rollback(tx, "Insert into Vendors failed.")
		return fmt.Errorf("Unable to insert into Vendors. Err=%v", err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		attempt_rollback(tx, "Unable to get Rows Affected.")
		return fmt.Errorf("Unable to get Rows Affected. Err=%v", err)
	}
	if rowCnt != 1 {
		attempt_rollback(tx, "Wrong rowcount after Insert")
		return fmt.Errorf("Wrong rowcount (%d), after Insert.", rowCnt)
	}
	err = add_vendor_aliases(tx, v.Vid, v.Aliases)
	if err != nil {
		attempt_rollback(tx, "Unable to add vendor aliases.")
		return fmt.Errorf("Unable to add vendor aliases. Err=%v", err)
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("Commit failed on Update Vendors. Err=%v", err)
		return fmt.Errorf("Commit failed on Update Vendors. Err=%v", err)
	}
	return nil
}

func attempt_rollback(tx *sql.Tx, reason string) {
	err := tx.Rollback()
	if err != nil {
		log.Errorf("Unable to rollback. Err=%v. Reason for Rollback: %s", err, reason)
	}
}

func delete_vendor_aliases(tx *sql.Tx, vid uuid.UUID, nexpected int) error {
	res, err := tx.Exec("Delete from VendorAlias where Vid = ?", vid.String())
	if err != nil {
		return fmt.Errorf("Unable to delete from VendorAlias. Err=%v", err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to get results from VendorAlias table delete op. Err=%v", err)
	}
	if int(rowCnt) != nexpected {
		return fmt.Errorf("Delete failed in VendorAlias table. Row affected should be %d, was %d.",
			nexpected, rowCnt)
	}
	return nil
}

func add_vendor_aliases(tx *sql.Tx, vid uuid.UUID, aliases []string) error {
	if len(aliases) <= 0 {
		return nil
	}
	stmt, err := tx.Prepare("Insert into VendorAlias(Vid, Alias) values(?, ?)")
	if err != nil {
		return fmt.Errorf("Unable to Prepare for insert into VendorAlias. Err=%v", err)
	}
	for _, a := range aliases {
		_, err = stmt.Exec(vid.String(), a)
		if err != nil {
			return fmt.Errorf("Failed to insert into VendorAlias. Err=%v", err)
		}
	}
	return nil
}

func get_vendor_by_uuid(tx *sql.Tx, vid uuid.UUID) (*Vendor, error) {
	var v Vendor

	scmd := "Select Vendors.Vid, DName, FName, DefaultCat, BusinessType, PrimaryProduct, Notes, Alias from Vendors "
	scmd += "Left Join VendorAlias on Vendors.Vid=VendorAlias.Vid where Vendors.Vid = ?"
	rows, err := tx.Query(scmd, vid.String())
	if err != nil {
		return nil, err
	}
	ncnt := 0
	for rows.Next() {
		vraw, err := scan_vendor(rows)
		if err != nil {
			return nil, fmt.Errorf("Scan failure in get_vendor_by_uuid. Err=%v.", err)
		}
		if !vraw.Vid.IsZero() {
			return nil, fmt.Errorf("Vendor with zero uuid returned from select statement!")
		}
		if vraw.Vid != vid {
			return nil, fmt.Errorf("Mismatched uuid (%s!=%s) returned from select statement!", vraw.Vid, vid)
		}
		alias := ""
		if len(vraw.Aliases) > 0 {
			alias = vraw.Aliases[0]
		}
		if ncnt == 0 {
			v = *vraw
			v.Aliases = make([]string, 0, 5)
		}
		if !util.Blank(alias) {
			v.Aliases = append(v.Aliases, alias)
		}
		ncnt += 1
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Database failure after iterating rows on Vendor table. Err=%v", err)
	}
	if ncnt <= 0 {
		return nil, nil
	}
	return &v, nil
}

func scan_vendor(rows *sql.Rows) (*Vendor, error) {
	var v Vendor
	var svid_null, dname_null, fname_null, defaultcat_null sql.NullString
	var businesstype_null, primaryproduct_null, notes_null, alias_null sql.NullString
	err := rows.Scan(&svid_null, &dname_null, &fname_null, &defaultcat_null,
		&businesstype_null, &primaryproduct_null, &notes_null, &alias_null)
	if err != nil {
		return &v, fmt.Errorf("Err during row scan in GetAllVendors. Shipping account. Err=%v.", err)
	}
	if svid_null.Valid {
		v.Vid, err = uuid.FromString(svid_null.String)
		if err != nil {
			return &v, fmt.Errorf("Invalid uuid (%q) found for vendor. Err=%v", svid_null.String, err)
		}
	} else {
		return &v, fmt.Errorf("NULL uuid found for vendor.")
	}
	if dname_null.Valid {
		v.DName = dname_null.String
	}
	if fname_null.Valid {
		v.FName = fname_null.String
	}
	if defaultcat_null.Valid {
		v.DefaultCat = defaultcat_null.String
	}
	if businesstype_null.Valid {
		v.BusinessType = businesstype_null.String
	}
	if primaryproduct_null.Valid {
		v.PrimaryProduct = primaryproduct_null.String
	}
	if notes_null.Valid {
		v.Notes = notes_null.String
	}
	if alias_null.Valid {
		v.Aliases = []string{alias_null.String}
	} else {
		v.Aliases = []string{}
	}
	return &v, nil
}

// func GetAllVendorAliases() []*VendorAlias {
// 	gVendorAliasCasheLock.Lock()
// 	defer gVendorAliasCasheLock.Unlock()
// 	if gVendorAliases != nil {
// 		return gVendorAliases
// 	}
// 	lst := make([]*VendorAlias, 0, 10)
// 	rows, err := m_db.Query("Select Vid, Alias from VendorAlias")
// 	if err != nil {
// 		log.Errorf("Err getting VendorAlias. Returning empty slice. Err=%v", err)
// 		return lst
// 	}
// 	for rows.Next() {
// 		var svid, alias string
// 		err = rows.Scan(&svid, &alias)
// 		if err != nil {
// 			log.Errorf("Err during row scan in VendorAlias. Shipping alias. Err=%v.", err)
// 			continue
// 		}
// 		uvid, err := uuid.FromString0(svid)
// 		if err != nil {
// 			log.Errorf("Invalid uuid (%s) in VendorAlias with Alias = %q. Shipping. Err=%v.",
// 				svid, alias, err)
// 			continue
// 		}
// 		lst = append(lst, &VendorAlias{Vid: uvid, Alias: alias})
// 	}
// 	gVendorAliases = lst
// 	return gVendorAliases
// }

// func GetVendorAliases() []*VendorAlias {
// 	gVendorAliasCasheLock.Lock()
// 	defer gVendorAliasCasheLock.Unlock()
// 	if gVendorAliases != nil {
// 		return gVendorAliases
// 	}
// 	lst := make([]*VendorAlias, 0, 10)
// 	rows, err := m_db.Query("Select Vid, Alias from VendorAlias")
// 	if err != nil {
// 		log.Errorf("Err getting VendorAlias. Returning empty slice. Err=%v", err)
// 		return lst
// 	}
// 	for rows.Next() {
// 		var svid, alias string
// 		err = rows.Scan(&svid, &alias)
// 		if err != nil {
// 			log.Errorf("Err during row scan in VendorAlias. Shipping alias. Err=%v.", err)
// 			continue
// 		}
// 		uvid, err := uuid.FromString0(svid)
// 		if err != nil {
// 			log.Errorf("Invalid uuid (%s) in VendorAlias with Alias = %q. Shipping. Err=%v.",
// 				svid, alias, err)
// 			continue
// 		}
// 		lst = append(lst, &VendorAlias{Vid: uvid, Alias: alias})
// 	}
// 	gVendorAliases = lst
// 	return gVendorAliases
// }

// func AddVendorAlias(Vid uuid.UUID, Alias string) error {
// 	var err error
// 	acclst := GetAccounts()
// 	bFound := false
// 	for _, a := range acclst {
// 		if a.Aid == Aid {
// 			bFound = true
// 			break
// 		}
// 	}
// 	if !bFound {
// 		return fmt.Errorf("Invalid account ID.")
// 	}
// 	aliaslst := GetAccountAliases()
// 	for _, a := range aliaslst {
// 		if a.Aid == Aid && strings.ToLower(strings.TrimSpace(a.Alias)) == strings.ToLower(strings.TrimSpace(Name)) {
// 			return fmt.Errorf("Alias already exists.")
// 		}
// 	}
// 	stmt, err := m_db.Prepare("insert AccountAlias set Aid=?, Alias=?, Notes=?")
// 	if err != nil {
// 		return fmt.Errorf("Err inserting into AccountAlias. Err=%v", err)
// 	}
// 	r, err := stmt.Exec(Aid, Name, Notes)
// 	if err != nil {
// 		return fmt.Errorf("Err inserting into AccountAlias. Err=%v", err)
// 	}
// 	if n, _ := r.RowsAffected(); n != 1 {
// 		err = fmt.Errorf("Wrong number of rows affected (%d) when adding an alias to an account.", n)
// 	}

// 	gVendorAliasCasheLock.Lock()
// 	defer gVendorAliasCasheLock.Unlock()
// 	gVendorCasheLock.Lock()
// 	defer gVendorCasheLock.Unlock()
// 	gVendorCache = nil
// 	gVendorAliases = nil
// 	return err
// }

// func DeleteVendorAlias(Aid int, Name string) (int, error) {
// 	stmt, err := m_db.Prepare("delete from AccountAlias where Aid=? and Alias=?")
// 	if err != nil {
// 		return 0, fmt.Errorf("Err deleting from AccountAlias. Err=%v", err)
// 	}
// 	r, err := stmt.Exec(Aid, Name)
// 	if err != nil {
// 		return 0, fmt.Errorf("Err deleting from AccountAlias. Err=%v", err)
// 	}
// 	nn, _ := r.RowsAffected()
// 	n := int(nn)
// 	if n != 0 && n != 1 {
// 		err = fmt.Errorf("Wrong number of rows affected (%d) when deleting an alias from an account.", n)
// 	}
// 	gVendorAliasCasheLock.Lock()
// 	defer gVendorAliasCasheLock.Unlock()
// 	gVendorCasheLock.Lock()
// 	defer gVendorCasheLock.Unlock()
// 	gAccountCache = nil
// 	gAccountAliases = nil
// 	return n, err
// }
