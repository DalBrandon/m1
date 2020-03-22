// --------------------------------------------------------------------
// admin.go -- Runs the admin page
//
// Created 2018-10-03 DLB for epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"dbe/lib/log"
	"dbe/m1/console"
	"encoding/base64"
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("Admin", Invoke_GET, authorizer, handle_admin)
	RegisterPage("AdminCommand", Invoke_POST, authorizer, handle_command)
}

func handle_admin(c *gin.Context) {
	if !IsAdmin(c) {
		SendMessagePagef(c, "How did you get here?<br><br>This page is only for Admin Users!<br>")
		return
	}
	data := GetHeaderData(c)
	data.StyleSheets = []string{"admin"}

	SendPage(c, data, "header", "menubar", "admin", "footer")
}

type CmdResponse struct {
	ErrorMessage  string
	CommandOutput string
}

func handle_command(c *gin.Context) {
	data := GetHeaderData(c)
	if !data.IsAdmin {
		log.Infof("Non Admin user attemting to use admin command!  Returning nothing.")
		sendCmdResponse(c, "")
		return
	}
	if c.ContentType() != "text/plain" {
		log.Infof("Wrong request type (%s) for admin command. Returning nothing.", c.ContentType())
		sendCmdResponse(c, "")
		return
	}

	raw, err := c.GetRawData()
	if err != nil {
		log.Infof("Unable to get raw data for admin command. Err=%v", err)
		sendCmdResponse(c, "")
		return
	}
	cmd_bytes, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		log.Infof("Unable to decode base64 for admin command.  Recieved=%q. Err=%v", string(raw), err)
		sendCmdResponse(c, "")
	}
	cmd := string(cmd_bytes)
	log.Infof("Admin command from %s received: %s", data.Designer, cmd)

	sout := console.ExecuteCommand(cmd)
	sendCmdResponse(c, sout)
}

func sendCmdResponse(c *gin.Context, txt string) {
	txt64 := base64.StdEncoding.EncodeToString([]byte(txt))
	response := &CmdResponse{}
	response.ErrorMessage = ""
	response.CommandOutput = txt64
	c.JSON(200, response)
}
