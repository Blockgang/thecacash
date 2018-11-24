FROM golang:latest as builder
MAINTAINER theca

RUN apt update && apt install -y curl
RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o syncer .

#FROM scratch
#COPY --from=builder /build/syncer /app/
#WORKDIR /app
#CMD ["./syncer"]

FROM ubuntu:18.04
# a lib (not sure which?) needed for api request
RUN apt update && apt install -y curl mysql-client iputils-ping iproute2 
#COPY blockchain_sync /tmp/
#CMD ["/tmp/blockchain_sync"]
COPY --from=builder /build/syncer /app/
WORKDIR /app
CMD ["./syncer"]
