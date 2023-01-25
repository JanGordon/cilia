# syntax=docker/dockerfile:1
FROM node:18-alpine
WORKDIR /app
COPY . .
RUN go build ./pkg/main.go
RUN ./main
EXPOSE 8080