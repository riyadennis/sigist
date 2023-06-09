FROM golang:1.19
RUN mkdir /app
WORKDIR /app
COPY . /app

RUN go mod download
ENV LOG_LEVEL=debug
ENV PORT=:8080
ENV ENVIRONMENT=dev
ENV MIGRATIONS_PATH=/app/migrations
RUN CGO_ENABLED=1 GOOS=linux go build -o main .
CMD ["/app/main"]
