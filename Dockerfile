# syntax=docker/dockerfile:1

FROM golang:1.17-alpine
RUN mkdir /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY * ./
RUN go get -u github.com/gin-gonic/gin
RUN go build -o ./colour-color
EXPOSE 1212
CMD ["./colour-color"]