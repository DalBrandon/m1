// --------------------------------------------------------------------
// types.go -- Types for the templates
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

type HeaderData struct {
	PageTabTitle  string   // Title of the browser tab.
	PageTitle     string   // Title of the page -- not always used.
	StyleSheets   []string // Extra style sheets
	BrowserWidth  int      // Width of browser window. 0 = default
	ContentWidth  int      // Content width of screen area, 0 = default
	IsLoggedIn    bool     // True if logged in
	UserName      string   // Name of user that is logged in. Normally same as Designer.
	Designer      string   // Name of designer logged in
	OnLoadFuncJS  string   // Function to execute in js on page load
	IsAdmin       bool     // True if user is admin
	Instructions  string   // Instructions for the page -- not always used.
	HideLoginLink bool     // If true, Login link won't be shown in header
	HideAboutLink bool     // If true, about link won't be shown in footer
	Message       string   // For generic messages and stuff.
	ErrorMessage  string   // For postback pages, when error occurs.
	DesignerHint  string   // For setting sel boxes for desiger
}
