package helloworld

// 生成 proto grpc
//go:generate protoc -I . -I ../../../proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. helloworld.proto

// 生成 proto http
//go:generate protoc -I . -I ../../../proto --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. helloworld.proto

// 生成 proto errors
// go:generate protoc -I . -I ../../../proto --go_out=paths=source_relative:. --go-errors_out=paths=source_relative:. helloworld.proto
