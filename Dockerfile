FROM golang:1.20rc1-alpine3.17

WORKDIR /usr/bot

COPY ./ ./

RUN go mod download

RUN go build -o runbot ./cmd/telegrambot/

CMD [ "./runbot" ]
