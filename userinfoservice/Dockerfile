#
#FROM golang:1.16-alpine as build
#
#WORKDIR /app
#
#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download
#
#COPY . .
#
#RUN go build -o ./userinfo ./cmd/server/main.go
#
#
#
#FROM alpine as proxy
#
#COPY --from=build /app/userinfo ./
#COPY --from=build /app/app.env ./

#EXPOSE 9091
#
#ENTRYPOINT [ "./userinfo" ]


FROM golang:1.16-alpine as build
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

#WORKDIR $GOPATH/src/smallest-golang/app/
WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/server/main.go



FROM scratch as userinfo

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

COPY --from=build /main .
COPY --from=build /app/app.env .


USER small-user:small-user

CMD ["./main"]