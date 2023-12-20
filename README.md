# go-http-server

> A HTTP server for sharing a local folder

# Features

- serve a plain folder to the web
- vue.js history mode for single page apps
  - dynamic HTML templating
- alternate webserver root `.ghs.yaml` based config

# Installation

install with Go

```sh
go install gitlab.com/bobymcbobs/go-http-server@latest

go-http-server
```

launch with Podman/Docker

```sh
podman run -it --rm -p 8080:8080 -v "$PWD:$PWD" --workdir "$PWD" registry.gitlab.com/bobymcbobs/go-http-server:latest
```

verify a container image with cosign

```sh
cosign verify \
    --certificate-identity-regexp 'https://gitlab.com/BobyMCbobs/go-http-server//.gitlab-ci.yml@(refs/heads/main|refs/tags/.*)' \
    --certificate-oidc-issuer-regexp 'https://gitlab.com' \
    -o text \
    registry.gitlab.com/bobymcbobs/go-http-server:latest
```

# Use cases

## Local web development

launch `go-http-server` in the directory of the built/rendered/source of a website locally.

## As a base layer

serving a website as a container image

```dockerfile
FROM registry.gitlab.com/bobymcbobs/go-http-server:latest
COPY public /var/run/ko
```

see [deployment](./docs/deployment.md).

## Serving a single page app

given a folder to serve with an *index.html*, rewrite all requests except assets to *index.html* with `APP_VUEJS_HISTORY_MODE` set to `true`.

Check out [templating configuration](./docs/configuration.md#templating).

## Self-serve dotfile config

a `.ghs.yaml` may be written to the serve folder, to configure a small subset of the server functions. Check out [dotfile configuration](./docs/configuration.md#dotfile-configuration).

# Documentation

Docs are located in the [docs](./docs/) folder, as well as [on GitLab pages](https://BobyMCbobs.gitlab.io/go-http-server).

This is a hard-fork of [https://gitlab.com/safesurfer/go-http-server](https://gitlab.com/safesurfer/go-http-server).

# History

Some time ago at Safe Surfer, there was a need to pass settings values from the infrastructure deployment and backend to the frontend so it can behave correctly depending on the environment, and to have a minimal and secure base image with a fast server.
Existing web servers don't provide this functionality, such as NGINX or Apache2.
At the time, there was a major rewrite for almost everything to be in Go and this fit into the ecosystem very well.
In an application using it, the functionality allowed values passed all the way from `helm install` to be passed into the frontend, if the plumbing is set in place.
A pre-configured server with all the needed features fit the purpose well.

# License

Copyright 2020-2021 Safe Surfer, 2020-2023 BobyMCbobs.
This project is licensed under the [AGPL-3.0](http://www.gnu.org/licenses/agpl-3.0.html) and is [Free Software](https://www.gnu.org/philosophy/free-sw.en.html).
This program comes with absolutely no warranty.

