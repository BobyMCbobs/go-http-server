# As a base image

```
FROM registry.gitlab.com/bobymcbobs/go-http-server:latest
ADD mysite /var/run/ko
```

# Helm

## Preliminary steps

Create a namespace:

```shell
kubectl create ns go-http-server
```

## Installation

Install with Helm:

```shell
helm install go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```

Note: to configure, please check out [the configuration docs](./configuration.md)

## Upgrading versions

Upgrade a release with Helm:

```shell
helm upgrade go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```

## Uninstalling

Uninstall with Helm:

```shell
helm uninstall go-http-server-dev \
  -n go-http-server-dev
```
