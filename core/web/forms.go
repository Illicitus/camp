package web

import (
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
)

func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func ParseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)

	if err := dec.Decode(dst, values); err != nil {
		return err
	}
	return nil
}
