// --------------------------------------------------------------------
// help.go -- Help page
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("Help", Invoke_GET, guest_auth, handle_help)
}

func handle_help(c *gin.Context) {
	data := GetHeaderData(c)
	SendPage(c, data, "header", "menubar", "help", "footer")
}
