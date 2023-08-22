FROM golang:1.18 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deployment

FROM alpine:latest

RUN apk add --no-cache tzdata

ENV TZ=America/Sao_Paulo

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/deployment /deployment

CMD ["tail", "-f", "/dev/null"]
