package main

import (
	"flag"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

var (
	listenAddr = flag.String("port", ":8080", "http listen address")
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("public"))

	m.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello, gophers!"))
		rw.Header().Set("Content-Type", "text/plain")
	})

	log.Printf("Starting server on 127.0.0.1:%s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, m))
}
