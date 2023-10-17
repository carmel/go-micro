# Gateway
[![Build Status](https://github.com/carmel/go-micro/gateway/workflows/Test/badge.svg?branch=main)](https://github.com/carmel/go-micro/gateway/actions?query=branch%3Amain)
[![codecov](https://codecov.io/gh/carmel/go-micro/gateway/branch/main/graph/badge.svg)](https://codecov.io/gh/carmel/go-micro/gateway)

HTTP -> Proxy -> Router -> Midware -> Client -> Selector -> Node

## Protocol
* HTTP -> HTTP  
* HTTP -> gRPC  
* gRPC -> gRPC  

## Encoding
* Protobuf Schemas

## Endpoint
* prefix: /api/echo/*
* path: /api/echo/hello
* regex: /api/echo/[a-z]+
* restful: /api/echo/{name}

## Midware
* cors
* auth
* color
* logging
* tracing
* metrics
* ratelimit
* datacenter
