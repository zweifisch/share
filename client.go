package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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

func fromStdin() []byte {
	content, _ := ioutil.ReadAll(os.Stdin)
	return content
}

type Client struct {
	url string
}

func (c Client) post(content []byte) {
	req, _ := http.NewRequest("PUT", c.url, bytes.NewReader(content))
	client := &http.Client{}
	resp, err := client.Do(req)
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

func (c Client) get(entry string) []byte {
	req, _ := http.NewRequest("GET", c.url+"/"+entry, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return body
}
