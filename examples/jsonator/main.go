package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"

	"bytes"
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

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")

	if err := enc.Encode(files); err != nil {
		return nil, err
	}

	resp := new(plugin_go.CodeGeneratorResponse)
	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String("output.json"),
		Content: proto.String(buf.String()),
	})

	return resp, nil
}

type file struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Services    []*service `json:"services"`
}

func newFile(fd *protokit.FileDescriptor) *file {
	svcs := make([]*service, len(fd.GetServices()))
	for i, sd := range fd.GetServices() {
		svcs[i] = newService(sd)
	}

	return &file{
		Name:        fmt.Sprintf("%s.%s", fd.GetPackage(), fd.GetName()),
		Description: fd.GetComments().String(),
		Services:    svcs,
	}
}

type service struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

func newService(sd *protokit.ServiceDescriptor) *service {
	methods := make([]string, len(sd.GetMethods()))
	for i, md := range sd.GetMethods() {
		methods[i] = md.GetName()
	}

	return &service{Name: sd.GetName(), Methods: methods}
}
