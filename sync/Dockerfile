FROM golang:latest as builder
MAINTAINER jengi77

RUN apt update && apt install -y curl
RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o syncer .

FROM scratch
COPY --from=builder /build/syncer /app/
WORKDIR /app
CMD ["./syncer"]
