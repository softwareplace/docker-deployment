FROM golang:1.18 AS builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deployment

FROM alpine:latest

RUN apk add --no-cache tzdata

RUN mkdir /app
WORKDIR /app

ENV TZ=America/Sao_Paulo

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/deployment ./
COPY --from=builder /app/runner ./

CMD ["sh", "./runner"]

