FROM golang:1.16-alpine

WORKDIR /aopp

COPY * .
RUN go mod download

EXPOSE 6000
CMD go run main.go