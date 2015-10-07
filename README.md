# Binder Registry

CRUD for Binder Templates.

## Go get it.

```bash
$ go get github.com/binder-project/binder-registry
```

## Build it.

```bash
$ cd $GOPATH/src/github.com/binder-project/binder-registry/simpleregistry
$ go build
$ BINDER_API_KEY=THISISMYTOKEN ./simpleregistry
2015/10/07 09:52:48 Serving on :8080
```

## Use it.

```
$ curl 127.0.0.1:8080
{"status": "Binder Registry Live!"}
$ curl 127.0.0.1:8080/templates
[]
$ curl -X POST 127.0.0.1:8080/templates -d '{"name": "myenv", "image-name":"jupyter/demo"}'
{"message":"Authorization header not set. Should be of format 'Authorization: token key'"}
$ # Need the token for POST and PUT
$ curl -X POST 127.0.0.1:8080/templates -d '{"name": "myenv", "image-name":"jupyter/demo"}' -H "Authorization: token THISISMYTOKEN"
{"name":"myenv","image-name":"jupyter/demo","command":"","limits":{},"time-created":"2015-10-07T14:54:55.782549453Z","time-modified":"2015-10-07T14:54:55.782549453Z"}
```
