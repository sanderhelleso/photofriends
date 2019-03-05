package views

import "html/template"

// note the "...string" syntax, it means that the function
// can take in "n" amount of arguments of the type string
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...) // spread strings from arr
	if err != nil { panic(nil) }

	return &View{
		Template: t,
	}
}

// View structure to initialize "n" amount
// of passed inn template paths
type View struct {
	Template *template.Template
}
