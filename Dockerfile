FROM golang:latest
WORKDIR /go/src/github.com/postgres-go
EXPOSE 8080
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main"]