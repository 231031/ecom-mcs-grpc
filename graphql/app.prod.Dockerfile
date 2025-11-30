FROM golang:1.25.3-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/app
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]