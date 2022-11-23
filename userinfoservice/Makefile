gen:
	 protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative  \
        proto/*.proto


gen:
    protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/googleapis/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  -I$GOPATH/src/github.com/envoyproxy/protoc-gen-validate \
  --go-grpc_out require_unimplemented_servers=false:$GOPATH/src \
  --go_out $GOPATH/src \
  --openapiv2_out . \
  --openapiv2_opt logtostderr=true \
  --grpc-gateway_out $GOPATH/src \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt generate_unbound_methods=true \
  --validate_out="lang=go:$GOPATH/src" \
  proto/*.proto
