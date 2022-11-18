FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod go.sum ./
# fetch dependancies
RUN go mod download && \
    go mod verify

COPY . .

RUN go build -o ./gatewaygo ./cmd/server/main.go

CMD [ "./gatewaygo" ]
EXPOSE 8080

