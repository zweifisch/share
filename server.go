package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

type Server struct {
	root string
}

func (s Server) index(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(s.root)
	ret := ""
	for _, f := range files {
		name := f.Name()
		if !strings.HasPrefix(name, ".") {
			ret += fmt.Sprintf("<li><a href=\"/%s\">%s</a></li>", name, name)
		}
	}
	tmpl := `<html>
	<link rel="stylesheet" href="/public/style.css">
	<ol>
	%s
	</ol>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, fmt.Sprintf(tmpl, ret))
}

func (s Server) createEntry(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)
	params := mux.Vars(r)
	entry := params["entry"]
	realpath := path.Join(s.root, entry)
	content := r.PostFormValue("content")
	if _, err := os.Stat(realpath); err == nil {
		now := time.Now()
		entry = fmt.Sprintf("%s-%02d%02d%02d", entry, now.Hour(), now.Minute(), now.Second())
		realpath = fmt.Sprintf("%s-%02d%02d%02d", realpath, now.Hour(), now.Minute(), now.Second())
	}
	ioutil.WriteFile(realpath, []byte(content), 0644)
	io.WriteString(w, "http://"+r.Host+"/"+entry)
}

func (s Server) getEntry(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)
	params := mux.Vars(r)
	realpath := path.Join(s.root, params["entry"])
	content, err := ioutil.ReadFile(realpath)
	if err != nil {
		http.Error(w, "Not Found", 404)
		return
	}
	tmpl := `<html>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="/public/styles/tomorrow.css">
	<pre><code>%s</code></pre>
	<script src="/public/highlight.pack.js"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, fmt.Sprintf(tmpl, html.EscapeString(string(content))))
}

func (s Server) serve(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/", s.index)
	router.HandleFunc("/{entry}", s.getEntry).Methods("GET")
	router.HandleFunc("/{entry}", s.createEntry).Methods("POST")
	http.Handle("/", router)
	http.Handle("/public/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
