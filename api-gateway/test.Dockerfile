FROM golang:1.12 

WORKDIR /home/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .