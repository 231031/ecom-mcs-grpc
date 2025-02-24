FROM golang:1.23.6-alpine3.21 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/app
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/account

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 50001
CMD ["app"]