FROM golang:latest as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o webserver .

FROM scratch
COPY --from=builder /build/webserver /app/
WORKDIR /app
EXPOSE 8000
CMD ["./webserver"]
