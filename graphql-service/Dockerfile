FROM golang:1.19
RUN mkdir /app
WORKDIR /app
COPY . /app

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o main .
CMD ["/app/main"]
