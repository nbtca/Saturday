FROM golang as builder
COPY . /app
WORKDIR /app
RUN go env -w CGO_ENABLED=0 &&\
  go build -v -o saturday .

FROM alpine:latest as deploy
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata  &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/saturday /app
COPY --from=builder /app/migrations /app

ENV Port=80

EXPOSE 80
ENTRYPOINT [ "./saturday" ]