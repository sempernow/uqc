# Stage 1 
#----------------------------------
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices
# https://docs.docker.com/engine/reference/builder/ 
# https://hub.docker.com/_/golang/ 
FROM golang:1.19.2-bullseye as builder

ARG PKG_NAME
ARG MODULE
ARG VENDOR
ARG SVN
ARG VER
ARG BUILT

ARG ARCH
# Settings for Golang compiler
ENV CGO_ENABLED 0
ENV GOARCH $ARCH
ENV GOHOSTARCH $ARCH

WORKDIR /work
COPY . .
#RUN go mod download -x

WORKDIR /work/app/${PKG_NAME}
RUN go build -a -ldflags="\
    -X '${MODULE}/app.Maker=${VENDOR}' \
    -X '${MODULE}/app.SVN=${SVN}' \
    -X '${MODULE}/app.Version=${VER}' \
    -X '${MODULE}/app.Built=${BUILT}'"

# Stage 2 
#----------------------------------
# https://hub.docker.com/_/alpine/ 
FROM alpine:3.16.3

ARG PKG_NAME
ARG ARCH
ARG HUB
ARG PRJ
ARG MODULE
ARG AUTHORS
ARG VENDOR
ARG SVN
ARG VER
ARG BUILT

#RUN apk --no-cache add curl
RUN apk --no-cache add jq

RUN mkdir -p /app/assets
RUN mkdir -p /app/cache

ENV PATH="/app:${PATH}"

COPY --from=builder /work/app/${PKG_NAME}/${PKG_NAME} /app/main

WORKDIR /app
CMD ["sleep", "1d"]
# CMD ["/app/main", "upsertpostschron"]

# https://github.com/opencontainers/image-spec/blob/master/annotations.md#pre-defined-annotation-keys 
LABEL image.authors="${AUTHORS}"
LABEL image.created="${BUILT}"
LABEL image.from="alpine"
LABEL image.hub="https://hub.docker.com/repository/docker/${HUB}/${PRJ}.${PKG_NAME}-${ARCH}"
LABEL image.revision="${SVN}"
LABEL image.source="https://${MODULE}/app/${PKG_NAME}"
LABEL image.title="${PKG_NAME}"
LABEL image.vendor="${VENDOR}"
LABEL image.version="${VER}"
