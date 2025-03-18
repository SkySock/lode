FROM golang:1.24-alpine AS builder
ARG SERVICE_NAME
ARG SERVICE_DIR=services/${SERVICE_NAME}

WORKDIR /app

COPY . .
RUN cd ${SERVICE_DIR} && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/${SERVICE_NAME} ./cmd/main.go