package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/atotto/clipboard"
)

func hasPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Mode()&os.ModeNamedPipe != 0
}

func fromClipBoard() string {
	content, _ := clipboard.ReadAll()
	return content
}

func fromStdin() string {
	content, _ := ioutil.ReadAll(os.Stdin)
	return string(content)
}

type Client struct {
	url string
}

func (c Client) post(content string) {
	resp, err := http.PostForm(c.url+"/entry", url.Values{"content": {content}})
	if err != nil {
		fmt.Print(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(string(body))
}
