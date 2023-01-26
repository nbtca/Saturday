FROM alpine:latest

RUN mkdir /app
WORKDIR /app/
ARG APP_PORT=80

COPY ./saturday .

ENV GIN_MODE=release
ENTRYPOINT ["./saturday"]
EXPOSE ${APP_PORT}