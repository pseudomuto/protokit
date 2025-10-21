package utils_test

import (
	"slices"
	"testing"

	"github.com/pseudomuto/protokit/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateGenRequest(t *testing.T) {
	t.Parallel()

	fds, err := utils.LoadDescriptorSet("..", "fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")
	require.Equal(t, []string{"booking.proto", "todo.proto"}, req.GetFileToGenerate())

	expectedProtos := []string{
		"booking.proto",
		"google/protobuf/any.proto",
		"google/protobuf/descriptor.proto",
		"google/protobuf/timestamp.proto",
		"google/protobuf/duration.proto",
		"extend.proto",
		"todo.proto",
		"todo_import.proto",
		"edition2023.proto",
		"edition2024.proto",
		"edition2023_implicit.proto",
	}

	for _, pf := range req.GetProtoFile() {
		require.True(t, slices.Contains(expectedProtos, pf.GetName()), "Unexpected proto file: %s", pf.GetName())
	}
}

func TestFilesToGenerate(t *testing.T) {
	t.Parallel()

	fds, err := utils.LoadDescriptorSet("..", "fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(fds, "booking.proto")
	protos := utils.FilesToGenerate(req)
	require.Len(t, protos, 1)
	require.Equal(t, "booking.proto", protos[0].GetName())
}

func TestLoadDescriptorSet(t *testing.T) {
	t.Parallel()

	set, err := utils.LoadDescriptorSet("..", "fixtures", "fileset.pb")
	require.NoError(t, err)
	require.Len(t, set.GetFile(), 11)

	require.NotNil(t, utils.FindDescriptor(set, "todo.proto"))
	require.Nil(t, utils.FindDescriptor(set, "whodis.proto"))
}

func TestLoadDescriptorSetFileNotFound(t *testing.T) {
	t.Parallel()

	set, err := utils.LoadDescriptorSet("..", "fixtures", "notgonnadoit.pb")
	require.Nil(t, set)
	require.EqualError(t, err, "open ../fixtures/notgonnadoit.pb: no such file or directory")
}

func TestLoadDescriptorSetMarshalError(t *testing.T) {
	t.Parallel()

	set, err := utils.LoadDescriptorSet("..", "fixtures", "todo.proto")
	require.Nil(t, set)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proto:")
}

func TestLoadDescriptor(t *testing.T) {
	t.Parallel()

	proto, err := utils.LoadDescriptor("todo.proto", "..", "fixtures", "fileset.pb")
	require.NotNil(t, proto)
	require.NoError(t, err)
}

func TestLoadDescriptorFileNotFound(t *testing.T) {
	t.Parallel()

	proto, err := utils.LoadDescriptor("todo.proto", "..", "fixtures", "notgonnadoit.pb")
	require.Nil(t, proto)
	require.EqualError(t, err, "open ../fixtures/notgonnadoit.pb: no such file or directory")
}

func TestLoadDescriptorMarshalError(t *testing.T) {
	t.Parallel()

	proto, err := utils.LoadDescriptor("todo.proto", "..", "fixtures", "todo.proto")
	require.Nil(t, proto)
	require.Error(t, err)
	require.Contains(t, err.Error(), "proto:")
}

func TestLoadDescriptorDescriptorNotFound(t *testing.T) {
	t.Parallel()

	proto, err := utils.LoadDescriptor("nothere.proto", "..", "fixtures", "fileset.pb")
	require.Nil(t, proto)
	require.EqualError(t, err, "FileDescriptor not found")
}
