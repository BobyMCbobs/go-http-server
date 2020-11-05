- [Container build](#sec-1)

> Packaging your site with safesurfer/go-http-server

# Container build<a id="sec-1"></a>

```dockerfile
FROM registry.gitlab.com/safesurfer/go-http-server:1.0.0
COPY dist /app/site
```
