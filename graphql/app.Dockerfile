FROM golang:1.25.3-alpine AS dev

# Set the folder name dynamically (change per microservice)
ARG FOLDER_NAME
ENV FOLDER_NAME=graphql

# Install necessary tools
RUN apk update && apk add --no-cache git curl

# Install air using Go
RUN go install github.com/air-verse/air@v1.63.0

WORKDIR /go/src/app

COPY go.mod go.sum .air.toml entrypoint.sh ./
RUN go mod download

COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql

# Copy and set up the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50001

# Use entrypoint script
CMD ["/entrypoint.sh"]

# CMD ["air", "-c", ".air.toml"]
# CMD ["go", "run", "./graphql/cmd/graphql/main.go"]


# FROM golang:1.25.3-alpine AS dev

# # Install necessary tools
# RUN apk --no-cache add gcc g++ make ca-certificates curl git

# WORKDIR /go/src/app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY account account
# COPY catalog catalog
# COPY order order
# COPY graphql graphql

# EXPOSE 50001

# CMD ["go", "run", "./graphql/cmd/graphql/main.go"]

