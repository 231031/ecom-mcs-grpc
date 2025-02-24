FROM golang:1.23.6-alpine3.21 AS dev

# Set the folder name dynamically (change per microservice)
ARG FOLDER_NAME
ENV FOLDER_NAME=account

# Install necessary tools
RUN apk update && apk add --no-cache git curl

# Install air using Go
RUN go install github.com/air-verse/air@latest

WORKDIR /go/src/app

COPY go.mod go.sum .air.toml entrypoint.sh ./
RUN go mod download

COPY account account

# Copy and set up the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50001

# Use entrypoint script
CMD ["/entrypoint.sh"]

# CMD ["air", "-c", ".air.toml"]
# CMD ["go", "run", "./account/cmd/account/main.go"]

