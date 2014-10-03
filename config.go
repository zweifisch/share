package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
)

type CfgServer struct {
	Port  int    `toml:"port"`
	Root  string `toml:"root"`
	Write bool   `toml:"write"`
}

type CfgClient struct {
	Remote string `toml:"server"`
}

type Config struct {
	Server CfgServer `toml:"server"`
	Client CfgClient `toml:"client"`
}

func expandTilda(path string) string {
	usr, _ := user.Current()
	if path[:2] == "~/" {
		return strings.Replace(path, "~/", usr.HomeDir+"/", 1)
	} else {
		return path
	}
}

func (c Config) dump(path string) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(expandTilda(path), buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Config) load(path string) {
	configPath := expandTilda(path)
	toml.DecodeFile(configPath, &c)
}
