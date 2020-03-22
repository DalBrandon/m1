// --------------------------------------------------------------------
// about.go -- Help page
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("About", Invoke_GET, guest_auth, handle_about)
}

func handle_about(c *gin.Context) {
	data := GetHeaderData(c)
	data.HideAboutLink = true
	SendPage(c, data, "header", "menubar", "about", "footer")
}
