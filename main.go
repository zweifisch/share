package main

import (
	"fmt"
	"reflect"
	"strings"

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
				if len(c.Args()) > 1 {
					keys := strings.Split(c.Args()[0], ".")
					value := c.Args()[1]
					if len(keys) == 2 {
						key := strings.ToUpper(keys[1][:1]) + keys[1][1:]
						switch keys[0] {
						case "client":
							reflect.ValueOf(&conf.Client).Elem().FieldByName(key).SetString(value)
						case "server":
							reflect.ValueOf(&conf.Server).Elem().FieldByName(key).SetString(value)
						default:
							fmt.Println("invalide key %s", keys[0])
							return
						}
						conf.dump("~/.share.toml")
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
