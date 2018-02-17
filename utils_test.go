package protokit_test

import (
	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/pseudomuto/protokit"
)

type UtilsTest struct {
	suite.Suite
}

func TestUtils(t *testing.T) {
	suite.Run(t, new(UtilsTest))
}

func (assert *UtilsTest) TestLoadDescriptor() {
	proto, err := protokit.LoadDescriptor("todo.proto", "fixtures", "fileset.pb")
	assert.NotNil(proto)
	assert.NoError(err)
}

func (assert *UtilsTest) TestLoadDescriptorFileNotFound() {
	proto, err := protokit.LoadDescriptor("todo.proto", "fixtures", "notgonnadoit.pb")
	assert.Nil(proto)
	assert.EqualError(err, "open fixtures/notgonnadoit.pb: no such file or directory")
}

func (assert *UtilsTest) TestLoadDescriptorMarshalError() {
	proto, err := protokit.LoadDescriptor("todo.proto", "fixtures", "todo.proto")
	assert.Nil(proto)
	assert.EqualError(err, "proto: can't skip unknown wire type 7 for descriptor.FileDescriptorSet")
}

func (assert *UtilsTest) TestLoadDescriptorDescriptorNotFound() {
	proto, err := protokit.LoadDescriptor("nothere.proto", "fixtures", "fileset.pb")
	assert.Nil(proto)
	assert.EqualError(err, "FileDescriptor not found")
}
