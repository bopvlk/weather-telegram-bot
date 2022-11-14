FROM golang:1.19.1-alpine

WORKDIR /bot

COPY . .

RUN go mod download

RUN go build -o /runbot ./cmd/telegrambot/

CMD [ "/runbot" ]