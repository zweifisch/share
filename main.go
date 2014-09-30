package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func hasPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

func fromStdin() string {
	content, _ := ioutil.ReadAll(os.Stdin)
	return string(content)
}

type Client struct {
	url string
}

func (c Client) post(content string, name string) {
	resp, err := http.PostForm(c.url+"/"+name, url.Values{"content": {content}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main() {
	var startServer = flag.Bool("server", false, "run server")
	flag.BoolVar(startServer, "s", false, "ailas for server")

	wd, _ := os.Getwd()
	var root = flag.String("root", wd, "root directory")

	var name = flag.String("as", "entry", "name")

	flag.Parse()

	if *startServer {
		server := Server{*root}
		server.serve(8909)
	} else {
		client := Client{"http://localhost:8909"}
		if hasPipe() {
			client.post(fromStdin(), *name)
		}
	}
}
