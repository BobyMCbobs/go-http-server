- [Container build](#sec-1)
  - [Simple](#sec-1-1)
  - [Simple + headers](#sec-1-2)
  - [Vuejs](#sec-1-3)
  - [Vuejs + template map](#sec-1-4)
  - [Vuejs + template map + headers](#sec-1-5)

> Packaging your site with safesurfer/go-http-server

# Container build<a id="sec-1"></a>

## Simple<a id="sec-1-1"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
env APP_SERVE_FOLDER=./site
COPY site /app/site
```

## Simple + headers<a id="sec-1-2"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
env APP_SERVE_FOLDER=./site \
  APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY site /app/site
COPY headers.yaml /app/headers.yaml
```

## Vuejs<a id="sec-1-3"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
env APP_VUEJS_HISTORY_MODE=true
COPY dist /app/dist
```

## Vuejs + template map<a id="sec-1-4"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
env APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml
COPY dist /app/dist
COPY templatemap.yaml /app/map.yaml
```

## Vuejs + template map + headers<a id="sec-1-5"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
env APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml \
  APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY dist /app/dist
COPY headers.yaml /app/headers.yaml
COPY templatemap.yaml /app/map.yaml
```
