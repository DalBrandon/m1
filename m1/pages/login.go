// --------------------------------------------------------------------
// login.go -- Handles Login Page
//
// Created 2018-09-22 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"dbe/lib/log"
	pv "dbe/m1/privilege"
	"dbe/m1/sessions"
	"dbe/m1/users"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	Name         string
	Year0        string
	Active       bool
	PassWordHash string
}

type LoginData struct {
	*HeaderData
	UsersJson string
}

func init() {
	RegisterPage("/Login", Invoke_GET, handle_login)
	RegisterPage("/LoginPost", Invoke_POST, handle_login_post)
}

func make_default_login_data() *LoginData {
	data := &LoginData{}
	data.HeaderData = &HeaderData{}
	data.PageTabTitle = "Brandon M1 Login"
	data.OnLoadFuncJS = "startUp"
	data.HideLoginLink = true
	return data
}

func handle_login(c *gin.Context) {
	kill_session(c)
	data := make_default_login_data()
	show_login_page(c, data)
}

type LoginSubmitData struct {
	User     string `form:"User"`
	Password string `form:"Password"`
}

func handle_login_post(c *gin.Context) {
	kill_session(c)
	data := make_default_login_data()

	var ld LoginSubmitData
	err := c.ShouldBind(&ld)
	if err != nil {
		log.Errorf("Login data failed to bind. Err=%v", err)
		data.ErrorMessage = "Web app failure!  Programming problem?"
		show_login_page(c, data)
		return
	}
	u := users.GetUser(ld.User)
	if u == nil {
		log.Infof("Attempt to login with unknown user (%q). Hacking attempt?", ld.User)
		data.ErrorMessage = "Login Failed."
		show_login_page(c, data)
		return
	}

	priv, ok := sessions.CheckPassword(ld.User, ld.Password)
	if ok {
		setup_login(c, ld.User, priv)
		c.Redirect(302, "/Accounts")
		return
	}

	log.Infof("Login failed for %s: bad password.", ld.User)
	data.ErrorMessage = "Login Failed."
	show_login_page(c, data)
}

func setup_login(c *gin.Context, user string, priv pv.Privilege) *sessions.TSession {

	ses := sessions.NewSession(user, c.ClientIP(), priv)
	ses.SetStringValue("UserHint", user)

	c.SetCookie("Cred", ses.AuthCookie, 0, "/", "", http.SameSiteLaxMode, false, true)
	log.Infof("New Login: %s (%s) with %s privilege.", user, c.ClientIP(), priv)
	return ses
}

func kill_session(c *gin.Context) {
	cookie, err := c.Cookie("Cred")
	if err == nil {
		ses, err := sessions.GetSessionByAuth(cookie)
		if err == nil {
			log.Infof("Logging off %s (%s).\n", ses.Name, ses.ClientIP)
			sessions.KillSession(cookie)
		}
	}
}

func show_login_page(c *gin.Context, data *LoginData) {
	ulst := users.GetUsers()
	u_bytes, err := json.MarshalIndent(ulst, "", "  ")
	if err != nil {
		SendErrorPagef(c, "Unabel to convert user list to json. <br>Err=%v", err)
		return
	}
	data.UsersJson = string(u_bytes)

	SendPage(c, data, "header", "login", "footer")
}
