- [Container build](#org002cf11)
  - [Simple](#orgef5aa61)
  - [Simple + headers](#org55ea3d0)
  - [Vuejs](#org9c32e8c)
  - [Vuejs + template map](#org304b41c)
  - [Vuejs + template map + headers](#org0dacdb6)
- [An instant HTTP server](#org6f656df)

> Packaging your site with safesurfer/go-http-server


<a id="org002cf11"></a>

# Container build


<a id="orgef5aa61"></a>

## Simple

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.2.0
COPY site /app/site
```


<a id="org55ea3d0"></a>

## Simple + headers

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.2.0
env APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY site /app/site
COPY headers.yaml /app/headers.yaml
```


<a id="org9c32e8c"></a>

## Vuejs

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.2.0
env APP_VUEJS_HISTORY_MODE=true
COPY dist /app/dist
```


<a id="org304b41c"></a>

## Vuejs + template map

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.2.0
env APP_SERVE_FOLDER=./dist \
  APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml
COPY dist /app/dist
COPY templatemap.yaml /app/map.yaml
```


<a id="org0dacdb6"></a>

## Vuejs + template map + headers

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.2.0
env APP_SERVE_FOLDER=./dist \
  APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml \
  APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY dist /app/dist
COPY headers.yaml /app/headers.yaml
COPY templatemap.yaml /app/map.yaml
```


<a id="org6f656df"></a>

# An instant HTTP server

A webserver, serving on port 8080, can be brought up instantly to share files.

```bash
go get -u gitlab.com/safesurfer/go-http-server
go-http-server
```
