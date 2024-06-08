package portal

import (
	"html/template"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	rootTemplate, err = template.ParseFiles(
		"E:\\Go_Code\\Distribute\\portal\\students.html",
		"E:\\Go_Code\\Distribute\\portal\\student.html")

	if err != nil {
		return err
	}

	return nil
}
