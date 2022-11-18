FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod go.sum ./
# fetch dependancies
RUN go mod download && \
    go mod verify

COPY . .

RUN go build -o ./servergo ./cmd/server/main.go

CMD [ "./servergo" ]
EXPOSE 9090
