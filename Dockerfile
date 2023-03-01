FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

RUN mkdir /app
WORKDIR /app/

COPY ./saturday .

ENV GIN_MODE=release