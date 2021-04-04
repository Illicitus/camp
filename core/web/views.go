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
	TemplateExt = ".gohtml"
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
	if err := v.Render(w, r, nil); err != nil {
		panic(err)
	}
}

func IsAuthenticated(r *http.Request, d *Data) {
	cookie, _ := r.Cookie("remember_token")
	if cookie != nil && cookie.Value != "" {
		d.IsAuthenticated = true
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		IsAuthenticated(r, &vd)
	default:
		vd = Data{
			Yield: data,
		}
		IsAuthenticated(r, &vd)
	}
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, please email us.", http.StatusInternalServerError)
		return nil
	}
	if _, err := io.Copy(w, &buf); err != nil {
		return err
	}
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
