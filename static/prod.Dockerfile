FROM golang:1.12-alpine as builder
RUN apk add git
WORKDIR /home/app
ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -i -o appserver ./cmd

FROM alpine:latest as appserver
RUN apk add ca-certificates
COPY --from=builder /home/app/appserver /bin/appserver
RUN touch .env
ENTRYPOINT ["/bin/appserver"]
