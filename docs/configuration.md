- [Environment variables](#sec-1)
- [Templating](#sec-2)
- [Environment variables](#sec-3)
- [Helm configuration](#sec-4)


# Environment variables<a id="sec-1"></a>

| Variable                  | Description                                                   | Default          |
|------------------------- |------------------------------------------------------------- |---------------- |
| `APP_ENV_FILE`            | The location of an env file to load in, during initialisation | `.env`           |
| `APP_HEALTH_PORT_ENABLED` | Enable binding of a health port                               | `true`           |
| `APP_HEALTH_PORT`         | The port to bind to for health checking                       | `:8081`          |
| `APP_PORT`                | The port to serve traffic on                                  | `:8080`          |
| `APP_METRICS_ENABLED`     | Enable binding of a metrics port                              | `true`           |
| `APP_PORT_METRICS`        | The port to bind for metrics traffic                          | `:2112`          |
| `APP_SERVE_FOLDER`        | The local folder path to serve                                | `./site`         |
| `APP_TEMPLATE_MAP_PATH`   | The path to a template map                                    | `./map.yaml`     |
| `APP_VUEJS_HISTORY_MODE`  | Enable Vuejs history mode path rewriting                      | `false`          |
| `APP_HEADER_SET_ENABLE`   | Enable header setting for requests                            | `false`          |
| `APP_HEADER_MAP_PATH`     | The path to the header map                                    | `./headers.yaml` |

# Templating<a id="sec-2"></a>

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

# Environment variables<a id="sec-3"></a>

Just like template map, environment variables can be highly dynamically set.

```yaml
Referrer-Policy:
  - strict-origin
X-Content-Type-Options:
  - "${APP_CONTENT_TYPE_OPTIONS}"
```

These headers will be added to the response. Headers configured this way can also be evaluated from environment variables.

# Helm configuration<a id="sec-4"></a>

| Parameter                                     | Description                                                                                       | Default                                         |
|--------------------------------------------- |------------------------------------------------------------------------------------------------- |----------------------------------------------- |
| serveFolder                                   | The local folder path to serve                                                                    | `/app/site`                                     |
| templateMap                                   | The template map of fields and environment variables to template in `$APP_SERVE_FOLER/index.html` | `{}`                                            |
| templateMapPath                               | The path to a template map                                                                        | `/app/map.yaml`                                 |
| vuejsHistoryMode                              | Enable Vuejs history mode path rewriting                                                          | `true`                                          |
| headerMap                                     | Custom headers to set on response                                                                 | `{}`                                            |
| headerMapPath                                 | The path to a header map                                                                          | `/app/headers.yaml`                             |
| labels                                        | Extra labels to add to all managed resources                                                      | `{}`                                            |
| extraEnv                                      | Declare extra environment variables                                                               |                                                 |
| image.repository                              | The repo where the image lives                                                                    | `registry.gitlab.com/safesurfer/go-http-server` |
| image.tag                                     | Specifies a tag of from the image to use                                                          | `latest`                                        |
| image.pullPolicy                              | Pod container pull policy                                                                         | `IfNotPresent`                                  |
| imagePullSecrets                              | References for the registry secrets to pull the container images in the Pod with                  | `[]`                                            |
| nameOverride                                  | Expand the name of the chart                                                                      | `""`                                            |
| fullNameOverride                              | Create a FQDN for the app name                                                                    | `""`                                            |
| serviceAccount.create                         | Whether a serviceAccount should be created for the Pod to use                                     | `false`                                         |
| serviceAccount.name                           | A name to give the servce account                                                                 | `nil`                                           |
| podAnnotations                                | Annotations to assign Pods                                                                        | `{}`                                            |
| podSecurityContext                            | Set a security context for the Pod                                                                | `{}`                                            |
| securityContext.readOnlyRootFilesystem        | Mount container filesytem as read only                                                            | `true`                                          |
| securityContext.runAsNonRoot                  | Don't allow the container in the Pod to run as root                                               | `true`                                          |
| securityContext.runAsUser                     | The user ID to run the container in the Pod as                                                    | `1000`                                          |
| securityContext.runAsGroup                    | The group ID to run the container in the Pod as                                                   | `1000`                                          |
| service.type                                  | The service type to create                                                                        | `ClusterIP`                                     |
| service.port                                  | The port to bind the app on and for the service to be set to                                      | `8080`                                          |
| ingress.enabled                               | Create an ingress manifests                                                                       | `false`                                         |
| ingress.realIPHeader                          | A header to forward, which contains the real client IP address                                    | `""`                                            |
| ingress.annotations                           | Set annotations for the ingress manifest                                                          | `{}`                                            |
| ingress.hosts                                 | The hosts which the ingress endpoint should be accessed from                                      |                                                 |
| ingress.tls                                   | References to TLS secrets                                                                         | `[]`                                            |
| resources                                     | Limits and requests for the Pods                                                                  | `{}`                                            |
| autoscaling.enabled                           | Enable autoscaling for the deployment                                                             | `false`                                         |
| autoscaling.minReplicas                       | The minimum amount of Pods to run                                                                 | `1`                                             |
| autoscaling.maxReplicas                       | The maximum amount of Pods to run                                                                 | `1`                                             |
| autoscaling.targetCPUUtilizationPercentage    | The individual Pod CPU amount until autoscaling occurs                                            | `80`                                            |
| autoscaling.targetMemoryUtilizationPercentage | The individual Pod Memory amount until autoscaling occurs                                         |                                                 |
| nodeSelector                                  | Declare the node labels for Pod scheduling                                                        | `{}`                                            |
| tolerations                                   | Declare the toleration labels for Pod scheduling                                                  | `[]`                                            |
| affinity                                      | Declare the affinity settings for the Pod scheduling                                              | `{}`                                            |
