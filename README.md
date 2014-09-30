# shared

share urls, code snips with teammate

## install

```
go get github.com/zweifisch/shared
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
$ grep error /var/log/app.log | share --as error
http://localhost:8909/error
```

from clipboard

```sh
$ share --as password
http://localhost:8909/password
```

file

```sh
$ share assets.tgz
http://localhost:8909/assets.tgz
```

connect to a remote server

```sh
$ cat source.py | share --server http://192.168.11.16:8909
```

## config

client config

```sh
$ share --config server http://192.168.1.10:8080
```

server config

```sh
$ share --config disable put  # only client from localhost can publish content
```

## screenshots

TBD

## development

```sh
go-bindata public/...
```
