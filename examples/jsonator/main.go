package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"

	"encoding/json"
	"fmt"
	"log"
)

func main() {
	if err := protokit.RunPlugin(new(plugin)); err != nil {
		log.Fatal(err)
	}
}

type plugin struct{}

func (p *plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(req)
	files := make([]*file, len(descriptors))

	for i, d := range descriptors {
		files[i] = newFile(d)
	}

	data, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	resp := new(plugin_go.CodeGeneratorResponse)
	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String("output.json"),
		Content: proto.String(string(data)),
	})

	return resp, nil
}

type file struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newFile(fd *protokit.FileDescriptor) *file {
	return &file{
		Name:        fmt.Sprintf("%s.%s", fd.GetPackage(), fd.GetName()),
		Description: fd.GetDescription(),
	}
}
