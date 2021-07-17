FROM golang:1.16.6-alpine3.14 AS builder

COPY . /github.com/dobriychelpozitivniy/telegram-go-pocket-bot/
WORKDIR /github.com/dobriychelpozitivniy/telegram-go-pocket-bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/dobriychelpozitivniy/telegram-go-pocket-bot/bin/bot .
COPY --from=0 /github.com/dobriychelpozitivniy/telegram-go-pocket-bot/configs configs/
COPY --from=0 /github.com/dobriychelpozitivniy/telegram-go-pocket-bot/.env .

EXPOSE 80

CMD ["./bot"]