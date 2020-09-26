package tmplfunc

import (
	"html/template"
)

func SafeHTML(s string) template.HTML     { return template.HTML(s) }
func SafeAttr(s string) template.HTMLAttr { return template.HTMLAttr(s) }
func SafeJS(s string) template.JS         { return template.JS(s) }
func SafeCSS(s string) template.CSS       { return template.CSS(s) }
