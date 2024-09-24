# ginx-server
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

```bash
$ ginx-server -f config.toml
2024-09-24 21:36:51 INF [ginx-server] logging in level: INFO
2024-09-24 21:36:52 INF [ginx-server] message queue is listening
2024-09-24 21:36:52 INF [ginx-server] created 0 cron jobs
2024-09-24 21:36:52 INF [ginx-server] server is listiening at 127.0.0.1:8080
```

