package main

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"

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

	app := cli.NewApp()
	app.Name = "share"
	app.Usage = "share code snips"

	clientFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "remote",
			Value: conf.Client.Remote,
			Usage: "remote server",
		},
	}

	app.Flags = clientFlags

	app.Action = func(c *cli.Context) {
		client := Client{c.String("remote")}
		if hasPipe() {
			client.post(fromStdin())
		} else if content := fromClipBoard(); len(content) > 0 {
			client.post([]byte(content))
		} else {
			fmt.Println("nothing in clipboard")
		}
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "start the server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "root",
					Value: conf.Server.Root,
					Usage: "server root directory",
				},
				cli.IntFlag{
					Name:  "port",
					Value: conf.Server.Port,
					Usage: "port to listen",
				},
			},
			Action: func(c *cli.Context) {
				server := Server{c.String("root"), 0, []string{}}
				server.serve(c.Int("port"))
			},
		},
		{
			Name:  "get",
			Usage: "get an entry",
			Flags: clientFlags,
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					client := Client{c.String("remote")}
					name := c.Args()[0]
					fmt.Print(string(client.get(name)))
				}
			},
		},
	}

	app.Run(os.Args)
}
