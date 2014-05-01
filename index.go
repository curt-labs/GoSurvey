package main

import (
	"flag"
	"github.com/curt-labs/GoSurvey/controllers/api/warranty"
	"github.com/curt-labs/GoSurvey/models/warranties"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"html/template"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("port", ":8080", "http listen address")
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))
	m.Use(render.Renderer(render.Options{
		Directory:       "views",
		Layout:          "layout",
		Extensions:      []string{".tmpl", ".html"},
		Funcs:           []template.FuncMap{},
		Delims:          render.Delims{"{{", "}}"},
		Charset:         "UTF-8",
		IndentJSON:      true,
		HTMLContentType: "text/html",
	}))

	m.Group("/api/warranty", func(r martini.Router) {
		r.Get("", warranty.All)
		r.Get("/:id", warranty.Get)
		r.Put("", binding.Bind(warranties.Warranty{}), warranty.Add)
		r.Delete("/:id", warranty.Delete)
	})

	log.Printf("Starting server on 127.0.0.1:%s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, m))
}
