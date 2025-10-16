package protokit_test

import (
	"log"

	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
)

type plugin struct{}

func (p *plugin) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(r)
	resp := new(pluginpb.CodeGeneratorResponse)

	for _, desc := range descriptors {
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(desc.GetName() + ".out"),
			Content: proto.String("Some relevant output"),
		})
	}

	return resp, nil
}

// An example of running a custom plugin. This would be in your main.go file.
func ExampleRunPlugin() {
	// in func main() {}
	if err := protokit.RunPlugin(new(plugin)); err != nil {
		log.Fatal(err)
	}
}
