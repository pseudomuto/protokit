package protokit_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
)

func TestRunPlugin(t *testing.T) {
	t.Parallel()

	fds, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")
	data, err := proto.Marshal(req)
	require.NoError(t, err)

	in := bytes.NewBuffer(data)
	out := new(bytes.Buffer)

	require.NoError(t, protokit.RunPluginWithIO(new(OkPlugin), in, out))
	require.NotEmpty(t, out)
}

func TestRunPluginInputError(t *testing.T) {
	t.Parallel()

	in := bytes.NewBufferString("Not a codegen request")
	out := new(bytes.Buffer)

	err := protokit.RunPluginWithIO(nil, in, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proto:")
	require.Empty(t, out)
}

func TestRunPluginNoFilesToGenerate(t *testing.T) {
	t.Parallel()

	fds, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(fds)
	data, err := proto.Marshal(req)
	require.NoError(t, err)

	in := bytes.NewBuffer(data)
	out := new(bytes.Buffer)

	err = protokit.RunPluginWithIO(new(ErrorPlugin), in, out)
	require.EqualError(t, err, "no files were supplied to the generator")
	require.Empty(t, out)
}

func TestRunPluginGeneratorError(t *testing.T) {
	t.Parallel()

	fds, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")
	data, err := proto.Marshal(req)
	require.NoError(t, err)

	in := bytes.NewBuffer(data)
	out := new(bytes.Buffer)

	err = protokit.RunPluginWithIO(new(ErrorPlugin), in, out)
	require.EqualError(t, err, "generator error")
	require.Empty(t, out)
}

type ErrorPlugin struct{}

func (ep *ErrorPlugin) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	return nil, errors.New("generator error")
}

type OkPlugin struct{}

func (op *OkPlugin) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	resp := new(pluginpb.CodeGeneratorResponse)
	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String("myfile.out"),
		Content: proto.String("someoutput"),
	})

	return resp, nil
}
