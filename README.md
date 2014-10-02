# share

share code snips with teammate

## install

```
go get github.com/zweifisch/share
```

## usage

start server

```sh
$ share --server
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
$ cat source.py | share --to http://192.168.11.16:8909
```

get entry

```sh
$ share --get 17
```

## config

client config

```sh
$ share --config remote http://192.168.1.10:8909
```

server config

```sh
$ share --config disable write # only client from localhost can publish content
```

## development

```sh
go-bindata public/...
```
