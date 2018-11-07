package main

//go:generate go build -o ./jsonator main.go
//go:generate mkdir -p third_party/google/api
//go:generate curl -sSL -o third_party/google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
//go:generate curl -sSL -o third_party/google/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto
//go:generate protoc --plugin=protoc-gen-jsonator=./jsonator -I. -Ithird_party --jsonator_out=. ./sample.proto ./sample2.proto
//go:generate rm -rf third_party
//go:generate rm ./jsonator
