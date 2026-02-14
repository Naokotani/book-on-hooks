package template

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"booksonhooks.ca/internal/domain"
)

type ServerErrorFunc func(w http.ResponseWriter, err error)

type Templates struct {
	TemplateCache map[string]*template.Template
	serverError   ServerErrorFunc
}

type TemplateData struct {
	CurrentYear int
	Row         int
	Col         int
	Form        any
	Book        *domain.Book
	Books       []domain.Book
	Machine     *domain.Machine
	Machines    []domain.Machine
	MachineRow  MachineRow
}

type MachineRow struct {
	books []domain.Book
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func New(serverError ServerErrorFunc) (*Templates, error) {
	tpl, err := newTemplateCache()
	if err != nil {
		return nil, fmt.Errorf("Failed to initalize template cache:\n%s", err)
	}
	return &Templates{TemplateCache: tpl, serverError: serverError}, nil
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.gotmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gotmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.gotmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (t *Templates) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func (t *Templates) Render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	ts, ok := t.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("Template %s does not exist in the cache", page)
		t.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		t.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
