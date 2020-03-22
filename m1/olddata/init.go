// --------------------------------------------------------------------
// init.go -- Finds the location of the olddata, if possible.
//
// Created 2020-03-22 DLB
// --------------------------------------------------------------------

package olddata

import (
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/m1/config"
	"strings"
)

var olddata_folder string = ""

func init() {
	log.Infof("Running Loader...")
	var ok bool
	olddata_folder, ok = config.GetParam("olddata_folder")
	if !ok {
		log.Infof("Configuration parameter 'olddata_folder' not provided.")
		log.Infof("Olddata will not be avaliable.")
		return
	}
	if !strings.HasSuffix(olddata_folder, "/") {
		olddata_folder += "/"
	}
	if !util.DirExists(olddata_folder) {
		log.Fatalf("Old data folder (%q) doesn't exist.", olddata_folder)
		log.Infof("Olddata will not be avaliable.")
		return
	}
	log.Infof("Location of olddata: %s", olddata_folder)
}
