FROM golang:1.18-alpine AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build .

FROM alpine:3.16.0
WORKDIR /app

COPY --from=build /app/reisen .
COPY --from=build /app/config.json .
COPY --from=build /app/templates ./templates

RUN ls
CMD ./reisen config.json
