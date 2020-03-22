// --------------------------------------------------------------------
// pages.go -- Handle all pages
//
// Created 2018-09-23 DLB For Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

type InvokeType string

const (
	Invoke_GET  = "GET"
	Invoke_POST = "POST"
)

type Page struct {
	Route    string
	Invoke   InvokeType
	Handlers []gin.HandlerFunc
}

var gPages []*Page

// Registor a web page
func RegisterPage(route string, invoke InvokeType, handlers ...gin.HandlerFunc) {
	if gPages == nil {
		gPages = make([]*Page, 0, 50)
	}
	p := &Page{route, invoke, handlers}
	gPages = append(gPages, p)
}

func GetAllPages() []*Page {
	return gPages
}
