# build stage
ARG BASE=golang:1.22.6-alpine3.20
FROM ${BASE} AS builder

ARG ALPINE_PKG_BASE="make git"
ARG ALPINE_PKG_EXTRA=""
ARG ADD_BUILD_TAGS=""

RUN apk add --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}
WORKDIR /app

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
ARG MAKE="make build"
RUN $MAKE

# final stage
FROM alpine:3.18
LABEL license='MIT license' 
LABEL Name=crumbDB Version=${VERSION}

RUN apk --no-cache add ca-certificates dumb-init
# Ensure using latest versions of all installed packages to avoid any recent CVEs
RUN apk --no-cache upgrade

COPY --from=builder /app/LICENSE /LICENSE
COPY --from=builder /app/res/ /res/
COPY --from=builder /app/follower /follower

EXPOSE 50055

CMD ["/follower"]
