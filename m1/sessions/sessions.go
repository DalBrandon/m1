// --------------------------------------------------------------------
// sessions.go -- Manages logged in users and sessions...
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package sessions

import (
	"dbe/lib/log"
	"dbe/lib/uuid"
	pv "dbe/m1/privilege"
	"dbe/m1/users"
	"fmt"
	"strings"
	"sync"
	"time"
)

type SessionPrivilege string

type TSession struct {
	Name       string
	ClientIP   string
	Privilege  pv.Privilege
	LastAccess time.Time
	LoginTime  time.Time
	Data       map[string]interface{}
	AuthCookie string
}

var gSessions []*TSession
var gSessionLock sync.Mutex
var gTimeToLive int = 7200 // In Seconds

var gBypassDev string = "" // If this string is empty, developer bypass mode is disabled
var gAllowUniversalPasswords bool = false

// SetDeveloperBypass sets a mode whereby a given name is used as a Designer at login time
// and therefore bypasses the normal login procedure.  Use an empty string to
// disable this feature.
func SetDeveloperBypass(developer string) {
	gBypassDev = strings.TrimSpace(developer)
}

func GetDeveloperBypass() string {
	return gBypassDev
}

// SetAllowUniversalPasswords enables or disables a mode whereby all users can use the
// same "unervisal" password.
func SetAllowUniversalPasswords(enable bool) {
	gAllowUniversalPasswords = enable
}

func GetAllowUniversalPasswords() bool {
	return gAllowUniversalPasswords
}

func init() {
	gSessions = make([]*TSession, 0, 20)
	go func() {
		for {
			time.Sleep(30)
			clean_sessions()
		}
	}()
}

func clean_sessions() {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	t := time.Now()
	lst := make([]*TSession, 0, len(gSessions))
	for _, s := range gSessions {
		elapsed := t.Sub(s.LastAccess)
		if elapsed < time.Duration(gTimeToLive)*time.Second {
			lst = append(lst, s)
		} else {
			log.Infof("Logging out %s (%s) due to inactivity.", s.Name, s.ClientIP)
		}
	}
	gSessions = lst
}

func SetTimeToLive(secs int) {
	gTimeToLive = secs
}

func GetTimeToLive() int {
	return gTimeToLive
}

func NewSession(name, ClientIP string, Privilege pv.Privilege) *TSession {
	session := new(TSession)
	session.Name = name
	session.ClientIP = ClientIP
	session.Privilege = Privilege
	session.LoginTime = time.Now()
	session.LastAccess = time.Now()
	session.Data = make(map[string]interface{}, 30)
	session.AuthCookie = uuid.New().String()
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	gSessions = append(gSessions, session)
	return session
}

func NewGuestSession(ipaddr string) *TSession {
	session := new(TSession)
	session.Name = "Guest"
	session.ClientIP = ipaddr
	session.Privilege = pv.Guest
	session.LoginTime = time.Now()
	session.LastAccess = time.Now()
	session.Data = make(map[string]interface{}, 30)
	session.AuthCookie = ""
	return session
}

func KillSession(AuthCookie string) {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	i := -1
	for j := 0; j < len(gSessions); j++ {
		if gSessions[j].AuthCookie == AuthCookie {
			i = j
			break
		}
	}
	if i < 0 {
		return
	}
	copy(gSessions[i:], gSessions[i+1:])
	gSessions[len(gSessions)-1] = nil
	gSessions = gSessions[:len(gSessions)-1]
}

func GetSessionByAuth(AuthCookie string) (*TSession, error) {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	for _, ss := range gSessions {
		if ss.AuthCookie == AuthCookie {
			return ss, nil
		}
	}
	return nil, fmt.Errorf("Session not found.")
}

func GetAllSessions() []*TSession {
	gSessionLock.Lock()
	defer gSessionLock.Unlock()
	lst := make([]*TSession, 0, len(gSessions))
	for _, s := range gSessions {
		lst = append(lst, s)
	}
	return lst
}

func (ses *TSession) GetStringValue(key string) string {
	t, ok := ses.Data[key]
	if !ok {
		return ""
	}
	s, _ := t.(string)
	return s
}

func (ses *TSession) SetStringValue(key, value string) {
	ses.Data[key] = value
}

func (ses *TSession) IsAdmin() bool {
	return ses.Privilege.IsAdmin()
}

func (ses *TSession) IsUser() bool {
	return ses.Privilege == pv.User
}

func (ses *TSession) IsGuest() bool {
	return ses.Privilege == pv.Guest
}

func (ses *TSession) HasWritePrivilege() bool {
	return ses.Privilege.HasWritePrivilege()
}

func (ses *TSession) HasReadPrivilege() bool {
	return ses.Privilege.HasReadPrivilege()
}

func (ses *TSession) ResetTimeoutClock() {
	ses.LastAccess = time.Now()
}

// CheckPassword will check a cleartext password with the hash table to
// determine a match.  If there is one, the allowed privilege is returned.
// This function applies the policy for the allow-universal-password mode.
func CheckPassword(username, cleartextpw string) (pv.Privilege, bool) {
	u := users.GetUser(username)
	if u == nil {
		return pv.None, false
	}
	if users.CheckPassword(username, cleartextpw) {
		return u.Priv, true
	}
	if gAllowUniversalPasswords {
		// Do something here?
		return pv.None, false
	}
	return pv.None, false
}
