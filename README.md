# shared

share code snips, logs and files with teammate

## install

```
go get github.com/zweifisch/shared
```

## usage

start server

```sh
shared --port 8888
```

client usage

from stdin

```sh
$ grep error /var/log/app.log | share --as error
http://localhost:8888/error
```

from clipboard

```sh
$ share --as password
http://localhost:8888/password
```

file

```sh
$ share assets.tgz
http://localhost:8888/assets.tgz
```

connect to a remote server

```sh
$ cat source.py | share --server http://192.168.11.16:8888
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

