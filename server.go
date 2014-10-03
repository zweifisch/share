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
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
)

type Server struct {
	root    string
	next    int
	entries []string
}

func (s *Server) scan() {
	files, _ := ioutil.ReadDir(s.root)
	for _, f := range files {
		name := f.Name()
		if _, err := strconv.Atoi(name); err == nil {
			s.register(name)
		}
	}
}

func (s *Server) register(entry string) int {
	s.entries = append(s.entries, entry)
	num, _ := strconv.Atoi(entry)
	if num >= s.next {
		s.next = num + 1
	}
	return num
}

func (s *Server) incr() int {
	return s.register(strconv.Itoa(s.next))
}

func (s Server) index(w http.ResponseWriter, r *http.Request) {
	entries := ""
	for _, entry := range s.entries {
		entries += fmt.Sprintf("<li><a href=\"/entry/%s\">%s</a></li>", entry, entry)
	}
	tmpl := `<html>
	<link rel="stylesheet" href="/public/style.css">
	<ol>%s</ol>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, fmt.Sprintf(tmpl, entries))
}

func (s *Server) createEntry(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)
	entry := strconv.Itoa(s.incr())
	realpath := path.Join(s.root, entry)
	if _, err := os.Stat(realpath); err == nil {
		if _, err := os.Stat(realpath); err == nil {
			http.Error(w, "file exists", 409)
			return
		}
	}
	dest, err := os.Create(realpath)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to crete file", 500)
		return
	}
	_, err = io.Copy(dest, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to crete file", 500)
		return
	}
	io.WriteString(w, "http://"+r.Host+"/"+entry)
}

func (s Server) getEntry(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)
	entry := r.URL.Path[1:]
	realpath := path.Join(s.root, entry)
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

func (s *Server) serve(port int) {
	s.root = expandTilda(s.root)
	if fi, err := os.Stat(s.root); err != nil || !fi.IsDir() {
		fmt.Printf("illegal root \"%s\" \n", s.root)
		return
	}
	s.scan()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			s.createEntry(w, r)
		case "GET":
			if r.URL.Path == "/" {
				s.index(w, r)
			} else {
				s.getEntry(w, r)
			}
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})
	http.Handle("/public/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	fmt.Printf("%d entries found, listen on port %d\n", len(s.entries), port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
