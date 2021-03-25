FROM golang:1.14.6-alpine3.11 AS build
WORKDIR /app
COPY main.go /app/
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY go.* /app/
ARG APP_BUILD_VERSION="0.0.0"
ARG APP_BUILD_HASH="???"
ARG APP_BUILD_DATE="???"
ARG APP_BUILD_MODE="development"
ARG GOARCH=""
RUN CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build \
  -a \
  -installsuffix cgo \
  -ldflags "-extldflags '-static' -s -w \
    -X gitlab.com/safesurfer/go-http-server/pkg/common.AppBuildVersion=$APP_BUILD_VERSION \
    -X gitlab.com/safesurfer/go-http-server/pkg/common.AppBuildHash=$APP_BUILD_HASH \
    -X gitlab.com/safesurfer/go-http-server/pkg/common.AppBuildDate=$APP_BUILD_DATE \
    -X gitlab.com/safesurfer/go-http-server/pkg/common.AppBuildMode=$APP_BUILD_MODE" \
  -o bin/go-http-server \
  main.go

FROM alpine:3.11 as extras
RUN apk add tzdata ca-certificates
RUN adduser -D user

FROM scratch
WORKDIR /app
ENV PATH=/app/bin
COPY --from=build /app/bin/go-http-server /app/bin/go-http-server
COPY --from=extras /etc/passwd /etc/passwd
COPY --from=extras /etc/group /etc/group
COPY --from=extras /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=extras /etc/ssl /etc/ssl
EXPOSE 8080
USER user
ENTRYPOINT ["/app/bin/go-http-server"]
