package web

import (
	"bytes"
	"errors"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

var (
	TemplateExt string = ".gohtml"
)

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(templateDir string, layoutDir string, layout string, files ...string) *View {
	addTemplatePath(templateDir, files)
	addTemplateExt(files)

	files = append(files, layoutFiles(layoutDir)...)

	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
	}).ParseFiles(files...)

	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, please email us.", http.StatusInternalServerError)
		return nil
	}
	io.Copy(w, &buf)
	return nil
}

func layoutFiles(layoutDir string) []string {
	files, err := filepath.Glob(layoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}

	return files
}

func addTemplatePath(templateDir string, files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
