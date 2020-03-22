// --------------------------------------------------------------------
// main.go -- Starting Page for m1
//
// Created 2020-03-15 DLB
// --------------------------------------------------------------------

package pages

import (
	"dbe/lib/log"
	"dbe/lib/util"
	"dbe/m1/users"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type AccountsDefaults struct {
	User string
	Year string
}

type AccountsData struct {
	*HeaderData
	Defaults  *AccountsDefaults
	UsersJson string
}

func init() {
	RegisterPage("/Accounts", Invoke_GET, authorizer, handle_accounts)
	RegisterPage("/SubmitAccounts", Invoke_POST, authorizer, handle_accounts_post)
}

func handle_accounts(c *gin.Context) {
	handle_accounts_with_error(c, "")
}

func handle_accounts_with_error(c *gin.Context, errmsg string) {
	data := &AccountsData{}
	data.HeaderData = GetHeaderData(c)
	data.PageTitle = "Accounts"
	data.Instructions = ""
	data.StyleSheets = []string{"accounts"}
	data.OnLoadFuncJS = "startUp"
	data.ErrorMessage = errmsg

	ulst := users.GetUsers()
	u_bytes, err := json.MarshalIndent(ulst, "", "  ")
	if err != nil {
		SendErrorPagef(c, "Unabel to convert user list to json. <br>Err=%v", err)
		return
	}
	data.UsersJson = string(u_bytes)

	//var err error
	// data.SelectionBoxData, err = GetSelectionBoxData()
	// if err != nil {
	// 	SendErrorPage(c, err)
	// 	return
	// }

	var sd *AccountsDefaults
	ses := GetSession(c)
	t, ok := ses.Data["AccountsDefaults"]
	if !ok {
		sd = &AccountsDefaults{}
	} else {
		sd, ok = t.(*AccountsDefaults)
		if !ok {
			log.Errorf("Unable to type convert AccountsDefaults in handle_accounts_with_error.")
			sd = &AccountsDefaults{}
		}
	}
	if util.Blank(sd.User) {
		sd.User = data.HeaderData.Designer
	}
	data.Defaults = sd
	SendPage(c, data, "header", "menubar", "accounts", "footer")
}

func handle_accounts_post(c *gin.Context) {
	handle_accounts_with_error(c, "Page not fully implemented yet.")
}
