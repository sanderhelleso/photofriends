package views

import "html/template"

// note the "...string" syntax, it means that the function
// can take in "n" amount of arguments of the type string
func NewView(layout string, files ...string) *View {
	files = append(files, 
		"views/layouts/layout.gohtml",
		"views/layouts/footer.gohtml",
	)

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
