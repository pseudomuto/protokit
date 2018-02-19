package main

//go:generate go build -o ./jsonator main.go
//go:generate protoc --plugin=protoc-gen-jsonator=./jsonator -I. --jsonator_out=. ./sample.proto ./sample2.proto
//go:generate rm ./jsonator
