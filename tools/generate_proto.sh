cd ..
protoc --go_out=internal/grpc \
       --go-grpc_out=internal/grpc \
        internal/grpc/proto/*.proto