FROM golang:1.19-alpine AS build

ADD . /server
WORKDIR /server

USER root

# Run the installer
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /server

FROM node:current-alpine AS frontend

ADD ./fe /frontend

WORKDIR /frontend
RUN npm install
RUN npm run build

# Build a new image, copy everything over and set the entrypoint
FROM alpine:3

COPY --from=frontend /frontend/dist /dist

ENV DASHBOARD_LOCATION=/dist
ENV DB_DSN=/pixie.sqlite?cache=shared&mode=rwc
ENV DB_DRIVER=sqlite


# The adding of the binary is the last step, so thatprevious steps can be cached
COPY --from=build /server /server

CMD ["/server"]
