## :dash: :dash: **The Binder Project is moving to a [new repo](https://github.com/jupyterhub/binderhub).** :dash: :dash:

:books: Same functionality. Better performance for you. :books:

Over the past few months, we've been improving Binder's architecture and infrastructure. We're retiring this repo as it will no longer be actively developed. Future development will occur under the [JupyterHub](https://github.com/jupyterhub/) organization.

* All development of the Binder technology will occur in the [binderhub repo](https://github.com/jupyterhub/binderhub)
* Documentation for *users* will occur in the [jupyterhub binder repo](https://github.com/jupyterhub/binder) 
* All conversations and chat for users will occur in the [jupyterhub binder gitter channel](https://gitter.im/jupyterhub/binder)

Thanks for updating your bookmarked links.

## :dash: :dash: **The Binder Project is moving to a [new repo](https://github.com/jupyterhub/binderhub).** :dash: :dash:

---

# Binder Registry

[![Build Status](https://travis-ci.org/binder-project/binder-registry.svg?branch=master)](https://travis-ci.org/binder-project/binder-registry)

CRUD for Binder Templates.

:warning: Prototype :warning:

This is an initial implementation of the [registry API from the Binder API spec proposal](https://github.com/jupyter/enhancement-proposals/pull/5).

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
