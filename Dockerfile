FROM golang AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies with retry and alternative proxy
RUN go env -w CGO_ENABLED=0 && \
    go env -w GOPROXY=https://proxy.golang.org,direct && \
    go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -v -o saturday .

FROM alpine:latest AS deploy
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata  &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/saturday /app
COPY --from=builder /app/migrations /app/migrations

ENV Port=80

EXPOSE 80
ENTRYPOINT [ "./saturday" ]