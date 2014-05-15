package main

import (
	"flag"
	"github.com/curt-labs/GoSurvey/controllers/api/prizes"
	"github.com/curt-labs/GoSurvey/controllers/api/surveys"
	"github.com/curt-labs/GoSurvey/controllers/api/warranty"
	"github.com/curt-labs/GoSurvey/models/warranties"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/logstasher"
	"github.com/martini-contrib/render"
	"github.com/mipearson/rfw"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	listenAddr  = flag.String("port", ":8080", "http listen address")
	environment = flag.String("env", "", "current app environment")
)

func main() {
	flag.Parse()
	m := martini.Classic()
	m.Use(gzip.All())
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

	if *environment == "prod" {
		logFile, err := rfw.Open("requests.log", 0644)
		if err == nil {
			m.Use(logstasher.Logger(logFile))
		}
	}

	m.Group("/api/warranty", func(r martini.Router) {
		r.Get("", warranty.All)
		r.Get("/:id", warranty.Get)
		r.Put("", binding.Bind(warranties.Warranty{}), warranty.Add)
		r.Delete("/:id", warranty.Delete)
	})

	m.Group("/api/survey", func(r martini.Router) {
		r.Get("", surveys.All)
		r.Get("/:id", surveys.Get)
		r.Post("/:id", surveys.Submit)
	})

	m.Group("/api/prize", func(r martini.Router) {
		r.Get("", prizes.All)
		r.Get("/current", prizes.Current)
		r.Get("/:id", prizes.Get)
	})

	m.Get("/**", func(rw http.ResponseWriter, req *http.Request, r render.Render) {
		bag := make(map[string]interface{}, 0)
		bag["Host"] = req.URL.Host
		r.HTML(200, "index", bag)
	})

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      m,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on 127.0.0.1:%s\n", *listenAddr)
	log.Fatal(srv.ListenAndServe())
}
