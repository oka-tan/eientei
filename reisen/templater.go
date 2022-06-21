package main

import (
	"fmt"
	"io"

	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
)

//Templater glues together mustache and echo.
type Templater struct {
	layout    *mustache.Template
	templates map[string]*mustache.Template
}

func mustCompile(filename string) *mustache.Template {
	template, err := mustache.ParseFile(filename)

	if err != nil {
		panic(fmt.Sprintf("Error loading %s: %s", filename, err))
	}

	return template
}

//NewTemplater initiates a Templater object.
//Templates need to be added in by hand.
func NewTemplater() *Templater {
	layout := mustCompile("templates/layout.html.mustache")

	templates := map[string]*mustache.Template{
		"index":              mustCompile("templates/index.html.mustache"),
		"board":              mustCompile("templates/board.html.mustache"),
		"board-error":        mustCompile("templates/board-error.html.mustache"),
		"board-empty":        mustCompile("templates/board-empty.html.mustache"),
		"board-search":       mustCompile("templates/board-search.html.mustache"),
		"board-search-error": mustCompile("templates/board-search-error.html.mustache"),
		"board-thread":       mustCompile("templates/board-thread.html.mustache"),
		"contact":            mustCompile("templates/contact.html.mustache"),
	}

	return &Templater{
		layout:    layout,
		templates: templates,
	}
}

//Render implements a method echo needs for template rendering.
func (t *Templater) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].FRenderInLayout(w, t.layout, data)
}
