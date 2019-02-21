# protoc-gen-pehredaar

Generates code for rights from proto messages for rights implememntation .

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./  ./src/rights/rights.proto

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./  ./src/pehredaar/pehredaar.proto

protoc -I /usr/local/include -I ./ -I ./src --go_out=plugins=grpc:./  ./src/example/example.proto

packr && go build && protoc -I /usr/local/include -I  ./ -I ./src --plugin=protoc-gen-pehredaar=protoc-gen-pehredaar  --pehredaar_out=:./  ./src/example/example.proto && goimports -w ./src