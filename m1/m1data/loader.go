// --------------------------------------------------------------------
// loader.go -- Loads the data from disk
//
// Created 2020-03-20 DLB
// --------------------------------------------------------------------

package m1data

import (
	"bytes"
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/m1/config"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var holddisk sync.Mutex
var backupfolder string = ""
var datafile string = ""
var backupfile string = ""
var rootname string = "m1data"

func init() {
	log.Infof("Running Loader...")
	df, ok := config.GetParam("data_folder")
	if !ok {
		log.Fatalf("Configuration parameter 'data_folder' not provided.")
	}
	if !strings.HasSuffix(df, "/") {
		df += "/"
	}
	if !util.DirExists(df) {
		log.Fatalf("Data folder (%q) doesn't exist.", df)
	}
	datafile = df + rootname + ".dat"
	backupfile = df + rootname + ".bck"
	backupfolder = df + "backups/"
	if !util.DirExists(backupfolder) {
		err := os.Mkdir(backupfolder, 0775)
		if err != nil {
			log.Fatalf("Unable to make backup directory. Err=%v", err)
		}
	}
	log.Infof("Main data file: %s\n", datafile)
	if util.FileExists(datafile) {
		log.Infof("Attempting to Load Data.")
		LoadData()
	} else {
		log.Infof("No Database Exiits!  Be sure to correct with backup or oldata.")
	}
}

// LoadData reads the snapshot on the disk into the current database.
// It will try to load the primary file, and if that fails, it will
// try the backup file.
func LoadData() error {
	dblock.Lock()
	defer dblock.Unlock()
	holddisk.Lock()
	defer holddisk.Unlock()
	if util.FileExists(datafile) {
		t0 := time.Now()
		d, err := read_file(datafile)
		telp := time.Now().Sub(t0).Seconds() * 1000.0
		if err == nil {
			db = d
			log.Infof("Database loaded from primary file. (%8.2f ms)", telp)
			return nil
		}
		log.Errorf("Failed to load database from primary file. Err=%v", err)
		log.Errorf("Primary file name: %s", datafile)
	}
	if !util.FileExists(backupfile) {
		log.Errorf("No backup file found. Database empty!")
		log.Errorf("Backup file name: %s", backupfile)
		return fmt.Errorf("Unable to load database.")
	}
	t0 := time.Now()
	d, err := read_file(backupfile)
	telp := time.Now().Sub(t0).Seconds() * 1000.0
	if err != nil {
		log.Errorf("Failed to load database from backup file. Err=%v", err)
		log.Errorf("Backup file name: %s", backupfile)
		return fmt.Errorf("Unalbe to load database.")
	}
	db = d
	log.Infof("Database loaded from backup file. (%f8.2 ms)", telp)
	return nil
}

// SaveData writes a snapshot of the current database to the disk. It also saves
// the current snapshot to a different name.
func SaveData() error {
	holddisk.Lock()
	defer holddisk.Unlock()
	if util.FileExists(datafile) {
		err := os.Rename(datafile, backupfile)
		if err != nil {
			return fmt.Errorf("Unable to make way for new data. Err=%v", err)
		}
	}
	dblock.Lock()
	defer dblock.Unlock()
	t0 := time.Now()
	err := write_file(db, datafile)
	telp := time.Now().Sub(t0).Seconds() * 1000.0
	if err != nil {
		log.Errorf("Unable to write to database file (%s). Err=%v", datafile, err)
		return fmt.Errorf("Unable to save database.")
	}
	log.Infof("Database saved to disk. (%8.2f ms)", telp)
	log.Infof("Location: %s\n", datafile)
	return nil
}

// GetBackupFileList returns a list of saved backup files.
func GetBackupFileList() []string {
	holddisk.Lock()
	defer holddisk.Unlock()
	files, err := ioutil.ReadDir(backupfolder)
	if err != nil {
		log.Errorf("Unable to list file from backup directory. Err=%v", err)
		return []string{}
	}
	lst := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		s := f.Name()
		if !strings.HasSuffix(s, ".dat") {
			continue
		}
		s = strings.Replace(s, backupfolder, "", 1)
		s = strings.Replace(s, ".dat", "", 1)
		s = strings.TrimSpace(s)
		lst = append(lst, s)
	}
	sort.Strings(lst)
	return lst
}

// DeleteBackup deletes a given backup file.
func DeleteBackup(fname string) error {
	holddisk.Lock()
	defer holddisk.Unlock()
	fn := backupfolder + fname + ".dat"
	err := os.Remove(fn)
	if err != nil {
		return fmt.Errorf("Unable to remove backup file. Err=%v\n", err)
	}
	return nil
}

// LoadBackup loads a Backup file into the current database.
func LoadBackup(fname string) error {
	fn := backupfolder + fname + ".dat"
	if !util.FileExists(fn) {
		return fmt.Errorf("File doesn't exist (%s).", fn)
	}
	t0 := time.Now()
	d, err := read_file(fn)
	telp := time.Now().Sub(t0).Seconds() * 1000.0
	if err != nil {
		err = fmt.Errorf("Unable to load backup database file (%s). Err=%v", err)
		log.Errorf("%v", err)
		return err
	}
	dblock.Lock()
	defer dblock.Unlock()
	db = d
	log.Infof("Backup file %s loaded into database. (%8.2f ms)", fname, telp)
	return nil
}

// SaveBackup writes the database to a backup file and returns
// the name of the file.  The rootname can be blank, in which case
// a rootname with the current time will be create.  The actual
// rootname used will be returned upon success.
func SaveBackup(rootname string) (string, error) {
	t0 := time.Now()
	fname := rootname
	if util.Blank(fname) {
		fname = "Backup_" + t0.Format("2006-01-02-15-04-05")
	}
	fn := backupfolder + fname + ".dat"
	dblock.Lock()
	defer dblock.Unlock()
	err := write_file(db, fn)
	telp := time.Now().Sub(t0).Seconds() * 1000.0
	if err != nil {
		log.Errorf("Unable to write backup file. Err=%v", err)
		log.Errorf("Name of backup file: %s", fn)
		return fname, fmt.Errorf("Unalble to write backup file. Err=%v", err)
	} else {
		log.Infof("Backup file (%s) written to disk. (%8.2f ms)", fn, telp)
	}
	return fname, nil
}

func read_file(fn string) (*Database, error) {
	var d Database
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return &d, fmt.Errorf("Unalbe to read database file. Err=%v", err)
	}
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&d)
	if err != nil {
		return &d, fmt.Errorf("Unable to decode database. Err=%v", err)
	}
	return &d, nil
}

func write_file(d *Database, fn string) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return fmt.Errorf("Unable to encode the database. Err=%v", err)
	}
	err = ioutil.WriteFile(fn, buf.Bytes(), 0775)
	if err != nil {
		return fmt.Errorf("Unable to write data. Err=%v", err)
	}
	return nil
}
