- [Helm](#sec-1)
  - [Preliminary steps](#sec-1-1)
  - [Installation](#sec-1-2)
  - [Upgrading versions](#sec-1-3)
  - [Uninstalling](#sec-1-4)


# Helm<a id="sec-1"></a>

## Preliminary steps<a id="sec-1-1"></a>

Create a namespace:

```shell
kubectl create ns go-http-server
```

## Installation<a id="sec-1-2"></a>

Install with Helm:

```shell
helm install go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```

Note: to configure, please check out [the configuration docs](./configuration.md)

## Upgrading versions<a id="sec-1-3"></a>

Upgrade a release with Helm:

```shell
helm upgrade go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```

## Uninstalling<a id="sec-1-4"></a>

Uninstall with Helm:

```shell
helm uninstall go-http-server-dev \
  -n go-http-server-dev
```
