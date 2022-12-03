# Stage 1 
#----------------------------------
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices
# https://docs.docker.com/engine/reference/builder/ 
# https://hub.docker.com/_/golang/ 
FROM golang:1.15.8 as builder

ARG PKG_NAME
ARG MODULE
ARG VENDOR
ARG SVN
ARG VER
ARG BUILT

ARG ARCH
# Settings for Golang compiler
# How to enable cgo @ Alpine build  https://megamorf.gitlab.io/2019/09/08/alpine-go-builds-with-cgo-enabled/
# RUN apk update
# RUN apk upgrade
# RUN apk add --update go=1.8.3-r0 gcc=6.3.0-r4 g++=6.3.0-r4
# ENV CGO_ENABLED 1
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
#FROM alpine:20210212
#FROM alpine:3.13.5
#FROM alpine:3.14.3
# @ Makefile : sed -i "s/COMMON_IMAGE/${COMMON_IMAGE}/" ${PATH_REL_DOCKER_BUILD}/svc.alpine.this.dockerfile
FROM COMMON_IMAGE

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

#ENV APP_SERVICE_NAME ${PKG_NAME}
#... security vulnerability ???

#RUN apk --no-cache add curl
RUN mkdir -p /app/assets
COPY --from=builder /work/infra/docker/build/healthcheck-svc.sh /app/healthcheck.sh

COPY --from=builder /work/app/${PKG_NAME}/${PKG_NAME} /app/main

WORKDIR /app
CMD ["/app/main"]

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
