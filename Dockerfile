# stage 1
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# stage 2
FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY .env.example .

RUN mv .env.example .env

COPY start.sh .
COPY wait-for.sh .
COPY /migration ./migration

EXPOSE 3031

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
