// --------------------------------------------------------------------
// config.go -- Manages the configuration file
//
// Created 2018-10-04 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package config

import (
	"bytes"
	"dbe/lib/util"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

const (
	ConfigFileName = "config.txt"
)

var gConfigParams map[string]string
var gError bool = true

func SetParam(key, value string) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return
		}
	}
	gConfigParams[key] = value
}

func GetParam(key string) (value string, ok bool) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return "", false
		}
	}
	v, ok := gConfigParams[key]
	return v, ok
}

func GetStringParam(key string, defaultval string) (value string, ok bool) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return defaultval, false
		}
	}
	v, ok := gConfigParams[key]
	if ok {
		return v, true
	} else {
		return defaultval, false
	}

}

func GetBoolParam(key string, defaultval bool) (value bool, ok bool) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return defaultval, false
		}
	}
	v, ok := gConfigParams[key]
	if !ok {
		return defaultval, false
	}
	v = strings.ToLower(v)
	if v == "true" || v == "t" || v == "yes" || v == "y" {
		return true, true
	}
	if v == "false" || v == "f" || v == "no" || v == "n" {
		return false, true
	}
	return defaultval, false
}

func SetBoolParam(key string, val bool) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return
		}
	}
	if val {
		gConfigParams[key] = "true"
	} else {
		gConfigParams[key] = "false"
	}
}

func RemoveParam(key string) {
	if gConfigParams == nil {
		_, _ = GetConfig()
		if gConfigParams == nil {
			return
		}
	}
	_, ok := gConfigParams[key]
	if !ok {
		return
	}
	delete(gConfigParams, key)
}

func GetConfig() (map[string]string, error) {
	if gConfigParams != nil {
		return gConfigParams, nil
	}
	gConfigParams = make(map[string]string, 10)
	data, err := ioutil.ReadFile(ConfigFileName)
	if err != nil {
		gError = true
		return gConfigParams, err
	}
	lines := strings.Split(string(data), "\n")
	ilinenum := 0
	for _, ln := range lines {
		ilinenum++
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "//") {
			continue
		}
		if util.Blank(ln) {
			continue
		}
		wrds := strings.Split(ln, "=")
		if len(wrds) != 2 {
			gError = true
			return gConfigParams, fmt.Errorf("Bad syntax on line %d. One equal char not found.\n", ilinenum)
		}
		key := strings.TrimSpace(wrds[0])
		val := strings.TrimSpace(wrds[1])
		gConfigParams[key] = val
	}
	gError = false
	return gConfigParams, nil
}

func WriteConfig() error {
	if gError {
		return fmt.Errorf("Unable to write config because it was not correctly read at startup.")
	}
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "// Config file writen by pnserver, on console command.\n")
	fmt.Fprintf(buf, "// Time writen: %s.\n\n", time.Now().Format("2006-01-02 15:04:05"))

	keys := make([]string, len(gConfigParams))
	for k, _ := range gConfigParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(buf, "%s=%s\n", k, gConfigParams[k])
	}
	fmt.Fprintf(buf, "\n")
	err := ioutil.WriteFile(ConfigFileName, buf.Bytes(), 0775)
	return err
}
