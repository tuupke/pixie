FROM golang:1.21-alpine AS build

ADD . /src
WORKDIR /src

USER root

# Run the installer
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -ldflags="-w -s" -o /cuproxy .

# Build a new image, copy everything over and set the entrypoint
FROM alpine:3

EXPOSE 631

COPY --from=build /cuproxy /cuproxy

CMD ["/cuproxy"]
