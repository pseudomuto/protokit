package main

//go:generate protoc --descriptor_set_out=fileset.pb --include_imports --include_source_info -I. ./booking.proto ./todo.proto
//go:generate protoc --descriptor_set_out=fileset_nopackage.pb --include_imports --include_source_info -I. ./todo_nopackage.proto
