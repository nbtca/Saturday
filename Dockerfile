FROM alpine:latest

RUN mkdir /app
WORKDIR /app/
ARG APP_PORT=${PORT}

COPY ./saturday .

ENV GIN_MODE=release
#ENTRYPOINT ["./saturday"]
EXPOSE ${APP_PORT}