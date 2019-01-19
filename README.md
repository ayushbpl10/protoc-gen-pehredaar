# protoc-gen-pehredaar

Generates code for rights from proto messages for rights implememntation .

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./example  ./example/proto/rights/rights.proto

protoc -I /usr/local/include -I ./  --go_out=plugins=grpc:./example  ./example/proto/rightsval/validator.proto

protoc -I /usr/local/include -I ./ -I /Users/appointy/Desktop/Desk/protoc-gen-rights/example/proto --go_out=plugins=grpc:./example  ./example/proto/example/example.proto

packr && go build && protoc -I /usr/local/include -I  ./ -I /Users/appointy/Desktop/Desk/protoc-gen-rights/example/proto --plugin=protoc-gen-rights=protoc-gen-rights  --rights_out=:./example  ./example/proto/example/example.proto && goimports -w ./example/pb
