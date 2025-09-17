# ---------- Builder ----------
FROM golang:1.25.1-alpine AS builder
WORKDIR /src

# cache dependencies
COPY go.mod go.sum ./
# RUN go mod tidy
RUN go mod download

COPY . .

# build
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /app/main ./main.go

# ---------- Runtime ----------
FROM alpine:3.20
WORKDIR /app

# set tz location
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main /app/main

# create .env
COPY .env.example /app/.env

# default ENV
ENV PORT=8080 \
    DB_URL="mongodb://mongo:27017/" \
    DB_NAME="cctv_db" \
    SECRETKEY="ABC123" 

EXPOSE 8080
CMD ["/app/main"]
