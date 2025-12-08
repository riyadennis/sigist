FROM golang:1.24 as builder
RUN mkdir /app
WORKDIR /app
COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o main .

FROM debian:bullseye-slim
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
