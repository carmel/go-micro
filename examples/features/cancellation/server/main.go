//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Package main is the server main package for cancellation demo.
package main

import (
	trpc "go-micro"
	"go-micro/examples/features/common"
	"go-micro/log"
	pb "go-micro/testdata/trpc/helloworld"
)

func main() {
	// Init server.
	s := trpc.NewServer()

	// Register service.
	pb.RegisterGreeterService(s, &common.GreeterServerImpl{})

	// Serve and listen.
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
