FROM golang:latest as builder

COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/

RUN go get -d -v

RUN go build -o /go/src/webserver

FROM scratch
COPY --from=builder /go/src/webserver /go/src/webserver

EXPOSE 8000
ENTRYPOINT ["/go/src/webserver"]
