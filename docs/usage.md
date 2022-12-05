> Packaging your site with safesurfer/go-http-server

# Container build

## Simple

Just serve your site

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.7.0
COPY site /app/site
```

## Simple + headers

Serve your site and set headers

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.7.0
env APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY site /app/site
COPY headers.yaml /app/headers.yaml
```

## Vuejs

Serve your site using [history mode](https://router.vuejs.org/guide/essentials/history-mode.html)

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.7.0
env APP_VUEJS_HISTORY_MODE=true
COPY dist /app/dist
```

## Vuejs + template map

Serve your site using history mode and Go html templating for index.html

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.7.0
env APP_SERVE_FOLDER=./dist \
  APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml
COPY dist /app/dist
COPY templatemap.yaml /app/map.yaml
```

## Vuejs + template map + headers

Serve your site using history mode and Go html templating for index.html, with setting headers

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.7.0
env APP_SERVE_FOLDER=./dist \
  APP_VUEJS_HISTORY_MODE=true \
  APP_TEMPLATE_MAP_PATH=/app/map.yaml \
  APP_HEADER_SET_ENABLE=true \
  APP_HEADER_MAP_PATH=./headers.yaml
COPY dist /app/dist
COPY headers.yaml /app/headers.yaml
COPY templatemap.yaml /app/map.yaml
```

# Import and use

Sometimes it is useful to implement go-http-server with importing it

```go
package main

import (
  common "gitlab.com/safesurfer/go-http-server/pkg/common"
  ghs "gitlab.com/safesurfer/go-http-server/pkg/httpserver"
)

func main() {
  ghs.NewWebServer().
    SetServeFolder(common.GetEnvOrDefault("KO_DATA_PATH", "./")).
    Listen()
}
```

# An instant HTTP server

A webserver, serving on port 8080, can be brought up instantly to share files.

```bash
go install gitlab.com/safesurfer/go-http-server@latest
go-http-server
```
