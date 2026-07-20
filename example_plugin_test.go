package protokit_test

import (
	"log"

	pluginpb "google.golang.org/protobuf/types/pluginpb"

	"github.com/pseudomuto/protokit"
)

type plugin struct{}

func (p *plugin) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(r)
	resp := new(pluginpb.CodeGeneratorResponse)

	for _, desc := range descriptors {
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    new(desc.GetName() + ".out"),
			Content: new("Some relevant output"),
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
