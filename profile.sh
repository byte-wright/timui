#!/bin/bash

go test -v -bench=. -cpuprofile=cpu.out -o bench.test ./cmd/gridtest/
go tool pprof -http=:8080 bench.test cpu.out