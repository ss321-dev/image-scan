FROM golang:1.14.4-alpine3.11 AS build-env
RUN apk --no-cache add gcc=9.3.0-r0 musl-dev=1.1.24-r2 git=2.24.3-r0
COPY . /go/src/github.com/ss321-dev/dockle-cmd/
WORKDIR /go/src/github.com/ss321-dev/dockle-cmd/
RUN go build -a -installsuffix cgo -ldflags "-s -w" -o /app

FROM alpine:3.11
ENV GOPATH=/go
COPY --from=build-env /app /app
ENTRYPOINT ["/app"]
