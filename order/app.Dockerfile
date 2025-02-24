FROM golang:1.23.6-alpine3.21 AS dev

# Set the folder name dynamically (change per microservice)
ARG FOLDER_NAME
ENV FOLDER_NAME=order

# Install necessary tools
RUN apk update && apk add --no-cache git curl

# Install air using Go
RUN go install github.com/air-verse/air@latest

WORKDIR /go/src/app

COPY go.mod go.sum .air.toml entrypoint.sh ./
RUN go mod download

COPY account account
COPY catalog catalog
COPY order order

# Copy and set up the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50001

# Use entrypoint script
CMD ["/entrypoint.sh"]

# CMD ["air", "-c", ".air.toml"]




# FROM golang:1.23.6-alpine3.21 AS dev

# # Install necessary tools
# RUN apk --no-cache add gcc g++ make ca-certificates curl git

# WORKDIR /go/src/app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY account account
# COPY catalog catalog
# COPY order order

# EXPOSE 50001

# CMD ["go", "run", "./order/cmd/order/main.go"]


