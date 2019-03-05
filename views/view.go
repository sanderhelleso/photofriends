package views

import (
	"html/template"
	"path/filepath"
	"net/http"
)

var (
	layoutDir 	= "views/layouts/"
	templateExt = ".gohtml"
)

// note the "...string" syntax, it means that the function
// can take in "n" amount of arguments of the type string
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...) // spread strings from arr
	if err != nil { panic(nil) }

	return &View {
		Layout: layout,
		Template: t,
	}
}

// View structure to initialize "n" amount
// of passed inn template paths aswell as
// a layout string resulting in a yield to render
type View struct {
	Layout string
	Template *template.Template
}

// Render is used to render the view with the predefined layout
func (v *View) Render(res http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(res, v.Layout, data)
}

// layout files return a slice of strings
// representing the layout files used in
// our application using globbing
func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil { panic(err) }
	
	return files;
}
