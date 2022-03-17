- [Helm](#org78ac655)
  - [Preliminary steps](#org28a66c3)
  - [Installation](#org9920be8)
  - [Upgrading versions](#org118a603)
  - [Uninstalling](#orgbb26cf8)



<a id="org78ac655"></a>

# Helm


<a id="org28a66c3"></a>

## Preliminary steps

Create a namespace:

```shell
kubectl create ns go-http-server
```


<a id="org9920be8"></a>

## Installation

Install with Helm:

```shell
helm install go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```

Note: to configure, please check out [the configuration docs](./configuration.md)


<a id="org118a603"></a>

## Upgrading versions

Upgrade a release with Helm:

```shell
helm upgrade go-http-server-dev \
  -n go-http-server-dev \
deployments/go-http-server
```


<a id="orgbb26cf8"></a>

## Uninstalling

Uninstall with Helm:

```shell
helm uninstall go-http-server-dev \
  -n go-http-server-dev
```
