FROM golang:1.25.3-alpine
WORKDIR /app
COPY . .
RUN go build -o parser handlers.go
EXPOSE 8080
CMD ["/app/parser"]
