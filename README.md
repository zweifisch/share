# share

share code snips with teammate

## install

```
go get github.com/zweifisch/share
```

## usage

start server

```sh
$ share server
listen on 8909
```

client usage

from stdin

```sh
$ grep error /var/log/app.log | share
http://localhost:8909/1
```

from clipboard

```sh
$ share
http://localhost:8909/2
```

connect to a remote server

```sh
$ cat source.py | share --remote http://192.168.11.16:8909
```

get an entry

```sh
$ share get 17 --remote http://192.168.11.16:8909
```

## config

```sh
$ share set client.remote http://192.168.1.10:8909
$ share set server.port 8909
$ share set server.root ~/.share
```

## development

```sh
go-bindata public/...
```
