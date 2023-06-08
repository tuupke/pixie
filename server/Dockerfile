FROM golang:1.19-alpine AS build

ADD . /server
WORKDIR /server

USER root

# Run the installer
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -ldflags="-w -s" -o /pixie .

FROM node:current-alpine AS frontend

ADD ./fe /frontend

WORKDIR /frontend
RUN npm install
RUN npm run build

# Build a new image, copy everything over and set the entrypoint
FROM alpine:3

ENV DASHBOARD_LOCATION=/frontend
ENV DB_DSN=/pixie.sqlite?cache=shared&mode=rwc
ENV DB_DRIVER=sqlite

EXPOSE 4000
EXPOSE 631
EXPOSE 5353/udp

RUN apk add --no-cache sqlite

COPY --from=frontend /frontend/dist /frontend
COPY --from=build /pixie /pixie

CMD ["/pixie"]
