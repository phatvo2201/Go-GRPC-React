FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod go.sum ./
# fetch dependancies
RUN go mod download && \
    go mod verify

COPY . .

RUN go build -o ./infogo ./cmd/server/main.go

CMD [ "./infogo" ]
EXPOSE 9091