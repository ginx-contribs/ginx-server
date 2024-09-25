# ginx-server

![Static Badge](https://img.shields.io/badge/go-1.23-blue)

```bash
$ ginx-server -f config.toml
2024-09-24 21:36:51 INF [ginx-server] logging in level: INFO
2024-09-24 21:36:52 INF [ginx-server] message queue is listening
2024-09-24 21:36:52 INF [ginx-server] server is listiening at 127.0.0.1:8080
```
ginx-server is a quickstart template for single http server project, features as bellow:

* ginx: integration with the ginx framework, supports graceful shutdown, hooks and more features.
* jwt: supports jwt authorization that contains access token and refresh token
* email: support register for email verification code
* ent: ent ORM framework, support datasource from mysql, postgresql, sqlite
* redis: supports redis cache
* mq: support message queue, default Redis Stream.
* wire: dependency injection with wire
* swagger: support generate swagger api document 
* makefile: build project with makefile


## commands

build project
```bash
$ make build
```
build project with all supports platforms
```bash
$ make build_all
```
generate swagger
```bash
$ make swag_gen
```
generate ent 
```bash
$ make ent_gen
```
generate wire
```bash
$ make wire
```

## how to use

clone this project
```bash
$ git clone git@github.com:ginx-contribs/ginx-server.git
```
checkout specify version
```bash
$ git checkout tags/v1.0.0
```
remove git dir
```bash
$ rm -rf .git
```
init your own git
```bash
$ git init
```