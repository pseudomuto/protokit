package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
	"github.com/pseudomuto/protokit"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func main() {
	if err := protokit.RunPlugin(new(plugin)); err != nil {
		log.Fatal(err)
	}
}

type plugin struct{}

func (p *plugin) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
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

	resp := new(pluginpb.CodeGeneratorResponse)
	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
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
		Description: fd.GetPackageComments().String(),
		Services:    svcs,
	}
}

type service struct {
	Name    string    `json:"name"`
	Methods []*method `json:"methods"`
}

func newService(sd *protokit.ServiceDescriptor) *service {
	methods := make([]*method, len(sd.GetMethods()))
	for i, md := range sd.GetMethods() {
		methods[i] = newMethod(md)
	}

	return &service{Name: sd.GetName(), Methods: methods}
}

type method struct {
	Name      string   `json:"name"`
	HTTPRules []string `json:"http_rules"`
}

func newMethod(md *protokit.MethodDescriptor) *method {
	httpRules := make([]string, 0)
	if httpRule, ok := md.OptionExtensions["google.api.http"].(*annotations.HttpRule); ok {
		switch httpRule.GetPattern().(type) {
		case *annotations.HttpRule_Get:
			httpRules = append(httpRules, fmt.Sprintf("GET %s", httpRule.GetGet()))
		case *annotations.HttpRule_Put:
			httpRules = append(httpRules, fmt.Sprintf("PUT %s", httpRule.GetPut()))
		case *annotations.HttpRule_Post:
			httpRules = append(httpRules, fmt.Sprintf("POST %s", httpRule.GetPost()))
		case *annotations.HttpRule_Delete:
			httpRules = append(httpRules, fmt.Sprintf("DELETE %s", httpRule.GetDelete()))
		case *annotations.HttpRule_Patch:
			httpRules = append(httpRules, fmt.Sprintf("PATCH %s", httpRule.GetPatch()))
		}
		// Append more for each rule in httpRule.AdditionalBindings...
	}

	return &method{Name: md.GetName(), HTTPRules: httpRules}
}
