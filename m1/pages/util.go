// --------------------------------------------------------------------
// util.go -- Utilities for pages and template processing.
//
// Created 2018-09-23 DLB for Epic
// Copied  2020-03-15 DLB for m1
// --------------------------------------------------------------------

package pages

import (
	"bytes"
	"dbe/lib/log"
	"dbe/lib/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"text/template"
)

func GetTemplate(name string) (*template.Template, error) {
	fn := "./static/templates/" + name + ".tmpl"
	t_bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Errorf("Missing template %q. (file %s). Err=%v",
			name, fn, err)
		return nil, err
	}
	tmpl, err := template.New(name).Parse(string(t_bytes))
	if err != nil {
		log.Errorf("Invalid template %q. Err=%v", name, err)
		return nil, err
	}
	return tmpl, nil
}

func MakePage(data interface{}, template_names ...string) ([]byte, error) {

	tmpls := make([]*template.Template, 0, len(template_names))
	for _, n := range template_names {
		t, err := GetTemplate(n)
		if err != nil {
			return []byte{}, err
		}
		tmpls = append(tmpls, t)
	}
	html := new(bytes.Buffer)
	for i, t := range tmpls {
		err := t.Execute(html, data)
		if err != nil {
			log.Errorf("Error execution template %q. Err=%v", template_names[i], err)
			return html.Bytes(), err
		}
	}
	return html.Bytes(), nil
}

func SendPage(c *gin.Context, data interface{}, template_names ...string) {
	html, err := MakePage(data, template_names...)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.Data(200, "text/html", html)
}

func FixWebLinkAddr(slink string) string {
	sout := slink
	if !util.Blank(sout) {
		s := strings.ToLower(sout)
		if !strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "http://") {
			sout = "http://" + sout
		}
	}
	return sout
}

func SendMessagePagef(c *gin.Context, f string, args ...interface{}) {
	data := GetHeaderData(c)
	data.Message = fmt.Sprintf(f, args...)
	SendPage(c, data, "header", "menubar", "message", "footer")
}

func SendErrorPage(c *gin.Context, err error) {
	SendErrorPagef(c, "%v", err)
}

func SendErrorPagef(c *gin.Context, f string, args ...interface{}) {
	data := GetHeaderData(c)
	data.ErrorMessage = fmt.Sprintf(f, args...)
	log.Errorf("Error Page Sent. Err: %s", data.ErrorMessage)
	SendPage(c, data, "header", "menubar", "error", "footer")
}
