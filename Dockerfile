FROM golang:1.20rc1-alpine3.17

WORKDIR /usr/bot/
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN ls -l
RUN go build -o /usr/bot/runbot ./cmd/telegrambot/
RUN ls -l
CMD /usr/bot/runbot
