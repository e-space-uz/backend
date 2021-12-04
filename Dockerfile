FROM golang:1.17 as builder

RUN mkdir -p $GOPATH/src/github.com/e-space-uz/backend
WORKDIR $GOPATH/src/github.com/e-space-uz/backend

COPY . ./

RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/backend /
FROM alpine
COPY --from=builder backend .
ENTRYPOINT [ "/backend" ]