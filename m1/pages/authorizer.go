// --------------------------------------------------------------------
// authorizer.go -- Authorizer: checks login and fills header data.
//
// Created 2018-09-29 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"dbe/lib/log"
	pv "dbe/m1/privilege"
	"dbe/m1/sessions"
	"fmt"
	"github.com/gin-gonic/gin"
)

func authorizer(c *gin.Context) {
	// If not logged in, redirect to login page.
	cookie, err := c.Cookie("Cred")
	ses, err := sessions.GetSessionByAuth(cookie)
	if err != nil {
		// We are not logged in!
		if bypass := sessions.GetDeveloperBypass(); bypass != "" {
			// Use developer bypass...
			log.Infof("Developer ByPass Login Mode Triggered...")
			ses = setup_login(c, bypass, pv.Admin)
		} else {
			log.Infof("Invalid auth cookie.. redirecting to Login.")
			c.Redirect(303, "/Login")
			c.Abort()
			return
		}
	}
	// Reset Session Timeout clock here...
	ses.ResetTimeoutClock()

	// Session found. Fill up data.
	c.Set("Session", ses)
	hdrdata := &HeaderData{}
	hdrdata.PageTabTitle = "Brandon m1"
	hdrdata.IsLoggedIn = true
	hdrdata.UserName = ses.Name
	hdrdata.Designer = ses.Name
	hdrdata.IsAdmin = ses.IsAdmin()
	c.Set("HeaderData", hdrdata)
}

func guest_auth(c *gin.Context) {
	cookie, err := c.Cookie("Cred")
	ses, err := sessions.GetSessionByAuth(cookie)
	if err != nil {
		c.Set("Session", sessions.NewGuestSession(c.ClientIP()))
		// We are not logged in!
		// Provide a guest level version of the header.
		hdrdata := &HeaderData{}
		hdrdata.PageTabTitle = "Brandon m1"
		hdrdata.IsLoggedIn = false
		hdrdata.Designer = ""
		hdrdata.UserName = ""
		hdrdata.IsAdmin = false
		c.Set("HeaderData", hdrdata)
		return
	}
	// Reset Session Timeout clock here...
	ses.ResetTimeoutClock()

	// Session found. Fill up data.
	c.Set("Session", ses)
	hdrdata := &HeaderData{}
	hdrdata.PageTabTitle = "Brandon m1"
	hdrdata.IsLoggedIn = true
	hdrdata.Designer = ses.Name
	hdrdata.UserName = ses.Name
	hdrdata.IsAdmin = ses.IsAdmin()
	c.Set("HeaderData", hdrdata)
}

// Note: GetHeaderDate should NEVER return nil for any
// page that has used authorizer or guest_auth.  Therefore
// don't bother checking for nil.  If nil is returned, allow
// your code to painic.
func GetHeaderData(c *gin.Context) *HeaderData {
	hdrdata, ok := c.Get("HeaderData")
	if !ok {
		err := fmt.Errorf("Header Data not avaliable.")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	hdr, ok := hdrdata.(*HeaderData)
	if !ok {
		err := fmt.Errorf("Header Data assert error! Programming BUG.")
		log.Errorf("%v\n", err)
		c.AbortWithError(400, err)
		return nil
	}
	return hdr
}

// GetSession returns session data for the current session.
// This call should always return a valid session for any
// page that used authorizer or guest_auth.
func GetSession(c *gin.Context) *sessions.TSession {
	x, ok := c.Get("Session")
	if !ok {
		err := fmt.Errorf("Session data not avaliable after authentication!?")
		log.Errorf("%v\n", err)
		log.Infof("Returning an empty session.")
		return &sessions.TSession{}
	}
	ses, ok := x.(*sessions.TSession)
	if !ok {
		err := fmt.Errorf("Session Data assert error! Programming BUG.")
		log.Errorf("%v\n", err)
		log.Infof("Returning an empty session.")
		return &sessions.TSession{}
	}
	return ses
}

// GetUser returns the currently logged in user if possible, otherwise
// and empty string is returned.
func GetUser(c *gin.Context) string {
	ses := GetSession(c)
	return ses.Name
}

// IsLoggedIn returns the login condition.
func IsLoggedIn(c *gin.Context) bool {
	ses := GetSession(c)
	if ses.IsAdmin() || ses.IsUser() || ses.IsGuest() {
		return true
	}
	return false
}

// IsAdmin returns true if the current context has admin rights.
func IsAdmin(c *gin.Context) bool {
	return GetSession(c).IsAdmin()
}

// IsUser return true if the current context has user rights.
func IsUser(c *gin.Context) bool {
	return GetSession(c).IsUser()
}

// IsGuest returns true if the current context has guest rights.
func IsGuest(c *gin.Context) bool {
	return GetSession(c).IsGuest()
}

// HasWritePrivilege returns true if current context can write
// to the database.
func HasWritePrivilege(c *gin.Context) bool {
	return GetSession(c).HasWritePrivilege()
}

// HasReadPrivilege returns true if current context can read
// from the database.
func HasReadPrivilege(c *gin.Context) bool {
	return GetSession(c).HasReadPrivilege()
}
