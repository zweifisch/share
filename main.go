package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("./")
	ret := ""
	for _, f := range files {
		ret += fmt.Sprintf("<li><a href=\"/%s\">%s</a></li>", f.Name(), f.Name())
	}
	tmpl := `<html>
	<link rel="stylesheet" href="/public/style.css">
	<body>
	%s
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, fmt.Sprintf(tmpl, ret))
}

func entry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	content, err := ioutil.ReadFile(params["entry"])
	if err != nil {
		http.Error(w, "Not Found", 404)
		return
	}

	tmpl := `<html>
	<link rel="stylesheet" href="/public/styles/tomorrow.css">
	<pre><code>%s</code></pre>
	<script src="/public/highlight.pack.js"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, fmt.Sprintf(tmpl, html.EscapeString(string(content))))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/{entry}", entry)
	http.Handle("/", r)
	http.Handle("/public/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}))
	http.ListenAndServe(":5000", nil)
}
