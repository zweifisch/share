package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"
)

func index() string {
	files, _ := ioutil.ReadDir("./")
	ret := ""
	for _, f := range files {
		ret += fmt.Sprintf("<li><a href=\"/%s\">%s</a></li>", f.Name(), f.Name())
	}
	tmpl := `<html>
	<link rel="stylesheet" href="/style.css">
	<body>
	%s
	</body>
	</html>
	`
	return fmt.Sprintf(tmpl, ret)
}

func entry(params martini.Params) (int, string) {
	content, err := ioutil.ReadFile(params["entry"])
	if err != nil {
		return 404, "not found"
	}

	tmpl := `<html>
	<link rel="stylesheet" href="/styles/tomorrow.css">
	<pre><code>%s</code></pre>
	<script src="/highlight.pack.js"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	</html>
	`
	return 200, fmt.Sprintf(tmpl, html.EscapeString(string(content)))
}

func main() {
	m := martini.Classic()
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
	})
	m.Get("/", index)
	m.Get("/:entry", entry)
	m.Run()
}
