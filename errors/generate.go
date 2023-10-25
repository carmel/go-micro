package errors

// 生成 pb
//go:generate protoc -I . -I ../proto --go_out=paths=source_relative:. --go-errors_out=paths=source_relative:. errors.proto
