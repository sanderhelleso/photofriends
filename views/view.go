package views

import (
	"html/template"
	"path/filepath"
	"net/http"
)

var (
	layoutDir 	= "views/layouts/"
	templateDir = "views/"
	templateExt = ".gohtml"
)

// note the "...string" syntax, it means that the function
// can take in "n" amount of arguments of the type string
func NewView(layout string, files ...string) *View {

	// prepeare paths
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...) // spread strings from slice
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

func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if err := v.Render(res, nil); err != nil {
		panic(nil)
	}
}

// Render is used to render the view with the predefined layout
func (v *View) Render(res http.ResponseWriter, data interface{}) error {
	res.Header().Set("Content-Type", "text/html")
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

// addTemplatePath takes in a slice of string
// representing file paths for templates, and
// it prepends the "templateDir" directory
// to each string in the slice
//
// eg: {"home"} -> {"views/home"} if templateDir == "views/"
func addTemplatePath(files[]string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and it
// appends the "templateExt" to each string in the slice
//
// eg: {"home"} -> {"home.gohtml"} of templateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}