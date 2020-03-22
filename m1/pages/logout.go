// --------------------------------------------------------------------
// logout.go -- Logout Page
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("Logout", Invoke_GET, authorizer, handle_logout)
}

func handle_logout(c *gin.Context) {
	kill_session(c)
	data := &HeaderData{}
	data.PageTabTitle = "Brandon M1"
	data.HideLoginLink = false
	data.HideAboutLink = false
	SendPage(c, data, "header", "logout", "footer")
}
