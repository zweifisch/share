package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/BurntSushi/toml"

	"os"
	"os/user"
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

func main() {
	cwd, _ := os.Getwd()
	conf := Config{CfgServer{8909, cwd, true},
		CfgClient{"http://localhost:8909"}}

	usr, _ := user.Current()
	configPath := path.Join(usr.HomeDir, ".share.toml")
	toml.DecodeFile(configPath, &conf)

	var startServer = flag.Bool("server", false, "run server")
	flag.BoolVar(startServer, "s", false, "ailas for server")
	var root = flag.String("root", conf.Server.Root, "root directory")
	var port = flag.Int("port", conf.Server.Port, "server port")
	var remote = flag.String("remote", conf.Client.Remote, "remote server")
	flag.Parse()

	if *startServer {
		server := Server{*root, 0, []string{}}
		server.serve(*port)
	} else {
		client := Client{*remote}
		if hasPipe() {
			client.post(fromStdin())
		} else if content := fromClipBoard(); len(content) > 0 {
			client.post([]byte(content))
		} else {
			fmt.Println("nothing in clipboard")
		}
	}
}
