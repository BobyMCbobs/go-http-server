# Environment variables

| Variable                            | Description                                                   | Default          |
|-------------------------------------|---------------------------------------------------------------|------------------|
| `APP_ENV_FILE`                      | The location of an env file to load in, during initialisation | `.env`           |
| `APP_HEALTH_PORT_ENABLED`           | Enable binding of a health port                               | `true`           |
| `APP_HEALTH_PORT`                   | The port to bind to for health checking                       | `:8081`          |
| `APP_PORT`                          | The port to serve traffic on                                  | `:8080`          |
| `APP_METRICS_ENABLED`               | Enable binding of a metrics port                              | `true`           |
| `APP_PORT_METRICS`                  | The port to bind for metrics traffic                          | `:2112`          |
| `APP_HTTP_REAL_IP_HEADER`           | The HTTP header to use for real IPs                           | `""`             |
| `APP_SERVE_FOLDER` / `KO_DATA_PATH` | The local folder path to serve                                | `./site`         |
| `APP_TEMPLATE_MAP_PATH`             | The path to a template map                                    | `./template-map.yaml`     |
| `APP_VUEJS_HISTORY_MODE`            | Enable Vuejs history mode path rewriting                      | `false`          |
| `APP_HEADER_SET_ENABLE`             | Enable header setting for requests                            | `false`          |
| `APP_HEADER_MAP_PATH`               | The path to the header map                                    | `./headers.yaml` |

# Templating

When Vuejs history mode is enabled, the `index.html` doc is able to be templated. This is useful for meta tags in HTML.

Example:

```yaml
siteName: my site
apiURL: https://api.example.com
```

```html
<html>
    <head>
        <title>{{ .siteName }}</title>
        <meta name="apiURL" content="{{ .apiURL }}">
    </head>
    <body>
        <h1>Example site</h1>
    </body>
</html>
```

Will render:

```html
<html>
    <head>
        <title>my site</title>
        <meta name="apiURL" content="https://api.example.com">
    </head>
    <body>
        <h1>Example site</h1>
    </body>
</html>
```

# Environment variables

Just like template map, environment variables can be highly dynamically set.

```yaml
  Referrer-Policy:
    - strict-origin
  X-Content-Type-Options:
    - "${APP_CONTENT_TYPE_OPTIONS}"
```

These headers will be added to the response. Headers configured this way can also be evaluated from environment variables.

# Dotfile configuration

When a file called `.ghs.yaml` exists in the serve folder, it will be loaded in for the web server configuration.
Current values for configuration are

```yaml
error404FilePath: string
headerMap:        map[string][]string
historyMode:      bool
redirectRoutes:   map[string]string
templateMap:      map[string]string
```

for overriding the value set by the server.

## Fields

The dotfile config supports a smaller and limited subset of the go-http-server settings. This is to ensure that in a self-service environment, certain configs cannot be set. The following fields are:

**error404FilePath**: the path to a html document to serve the file not found message
**headerMap**: a key+value-array pair to set headers. Values are env-evaluated (e.g: `X-Something-Important: ["Value-Here", "${SOME_ENV}"]`)
**historyMode**: when set, rewrites all requests with the exception of assets to _index.html_
**redirectRoutes**: a key+value pair to direct paths URLs to other URLs. (e.g: `/a: /b`)
**templateMap**: combined with `historyMode`, use Go html templating to replace Go templating expressions in an _index.html_

