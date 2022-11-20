# FROM golang:1.16-alpine
# WORKDIR /app
# COPY go.mod go.sum ./
# # fetch dependancies
# RUN go mod download && \
#     go mod verify

# COPY . .

# RUN go build -o ./gatewaygo ./cmd/server/main.go

# CMD [ "./gatewaygo" ]
# EXPOSE 8080

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /gtwaygo ./cmd/server/main.go
EXPOSE 8080

CMD [ "/gtwaygo" ]

