FROM golang:1.21.4-alpine
WORKDIR /app/go-todo
COPY . .
RUN go build -o ./build/main ./cmd/go_todo/
EXPOSE 9000
CMD ["./build/main"]