# Multi-step image builder to make smallest containers for deployment on poor man's cloud server

####################################################################################
# Base builder image.
FROM golang:1.19 as base
# Install OS and Go dependencies for build and code generation
RUN apt-get update
# Set up working directory.
RUN mkdir -p /work/bin
WORKDIR /work

ARG SERVICE
####################################################################################
# Builder for msg-scheduler
FROM base AS builder

# Copy over whole source tree.
COPY . app

WORKDIR /work/app

RUN pwd && ls -l

RUN go mod download
# Do all code generation (migrations and mocks)
RUN go generate ./...

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.4 #latest version is currently broken (august 30th)

RUN swag init -g api/main.go --parseDependency --parseInternal --parseDepth 3

WORKDIR /work/app/$SERVICE

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /work/bin/app/$SERVICE

#####################################################################################
# For root CA certificates.
FROM alpine:3.9 as ca
RUN apk add -U --no-cache ca-certificates

#####################################################################################
# Build minimal deployment image.
#FROM scratch
FROM alpine:3.4 
COPY --from=ca /etc/passwd /etc/passwd
COPY --from=ca /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /work/bin/app/$SERVICE .
ENTRYPOINT ["./app/${SERVICE}"]
