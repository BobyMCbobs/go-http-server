# Deployment

> How to deploy go-http-server

## As a base image

With a Dockerfile

```
FROM registry.gitlab.com/bobymcbobs/go-http-server:latest
ADD mysite /var/run/ko
```

see [usage](./usage.md) for more examples.

With [crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane)

```
BASE_IMAGE=registry.gitlab.com/bobymcbobs/go-http-server:latest
mkdir -p output/var/run
mv some/path/to/dir output/var/run/ko
crane append --base="$BASE_IMAGE" --new_layer=<(cd output/ && tar -f - -c .) --new_tag=registry.example.com/image:tag
```

please note: the intended _output_ folder must be resulting in /var/run/ko inside, as that's the default serve location for a built container.
