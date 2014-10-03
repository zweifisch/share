package main

import (
	"fmt"

	"github.com/codegangsta/cli"

	"os"
)

func main() {
	cwd, _ := os.Getwd()
	conf := Config{CfgServer{8909, cwd, true},
		CfgClient{"http://localhost:8909"}}
	conf.load("~/.share.toml")

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
		{
			Name:  "set",
			Usage: "global config",
			Action: func(c *cli.Context) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("can't set %s\n", c.Args()[0])
					}
				}()
				if len(c.Args()) > 1 {
					conf.set(c.Args()[0], c.Args()[1])
					conf.dump("~/.share.toml")
				}
			},
		},
	}

	app.Run(os.Args)
}
