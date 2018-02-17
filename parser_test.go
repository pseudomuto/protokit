package protokit_test

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/pseudomuto/protokit"
)

var proto *descriptor.FileDescriptorProto

type ParserTest struct {
	suite.Suite
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTest))
}

func (assert *ParserTest) SetupSuite() {
	var err error
	proto, err = protokit.LoadDescriptor("todo.proto", "fixtures", "fileset.pb")
	assert.NoError(err)
}

func (assert *ParserTest) TestParseFile() {
	file := protokit.ParseFile(proto)
	assert.Contains(file.GetDescription(), "The official documentation for the Todo API.\n\n")
}

func (assert *ParserTest) TestParseFileEnums() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetEnums(), 1)
	assert.Nil(file.GetEnum("swingandamiss"))

	enum := file.GetEnum("ListType")
	assert.Equal("An enumeration of list types", enum.GetDescription())
	assert.Equal("com.pseudomuto.protokit.v1", enum.GetPackage())
	assert.Len(enum.GetValues(), 2)

	assert.Equal("REMINDERS", enum.GetValues()[0].GetName())
	assert.Equal("The reminders type.", enum.GetNamedValue("REMINDERS").GetDescription())

	assert.Nil(enum.GetNamedValue("whodis"))
}

func (assert *ParserTest) TestParseFileServices() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetServices(), 1)
	assert.Nil(file.GetService("swingandamiss"))

	svc := file.GetService("Todo")
	assert.Contains(svc.GetDescription(), "A service for managing \"todo\" items.\n\n")
	assert.Equal("com.pseudomuto.protokit.v1", svc.GetPackage())
	assert.Len(svc.GetMethods(), 2)

	m := svc.GetMethods()[0]
	assert.Equal("Create a new todo list", m.GetDescription())

	assert.Equal("com.pseudomuto.protokit.v1", m.GetInputRef().GetPackage())
	assert.Equal("CreateListRequest", m.GetInputRef().GetTypeName())
	assert.True(m.GetInputRef().GetFullyQualified())

	assert.Equal("com.pseudomuto.protokit.v1", m.GetOutputRef().GetPackage())
	assert.Equal("CreateListResponse", m.GetOutputRef().GetTypeName())
	assert.True(m.GetInputRef().GetFullyQualified())

	m = svc.GetMethods()[1]
	assert.Equal("/com.pseudomuto.protokit.v1.Todo/AddItem", m.GetURL())
	assert.Equal("Add an item to your list\n\nAdds a new item to the specified list.", m.GetDescription())
}

func (assert *ParserTest) TestParseFileMessages() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetMessages(), 6)
	assert.Nil(file.GetMessage("swingandamiss"))

	m := file.GetMessage("AddItemRequest")
	assert.Equal("A request message for adding new items.", m.GetDescription())
	assert.Equal("com.pseudomuto.protokit.v1", m.GetPackage())
	assert.Len(m.GetMessageFields(), 3)
	assert.Nil(m.GetMessageField("swingandamiss"))

	f := m.GetMessageField("completed")
	assert.Equal("Whether or not the item is completed.", f.GetDescription())
}

func (assert *ParserTest) TestParseFileMessageEnums() {
	m := protokit.ParseFile(proto).GetMessage("Item")
	assert.Len(m.GetEnums(), 1)
	assert.Nil(m.GetEnum("whodis"))

	e := m.GetEnum("Status")
	assert.Equal("Item.Status", e.GetName())
	assert.Equal(e, m.GetEnum("Item.Status"))
	assert.Equal("An enumeration of possible statuses", e.GetDescription())
	assert.Len(e.GetValues(), 2)

	assert.Equal("COMPLETED", e.GetValues()[1].GetName())
	assert.Equal("The completed status.", e.GetNamedValue("COMPLETED").GetDescription())
}

func (assert *ParserTest) TestParseFileNestedMessages() {
	m := protokit.ParseFile(proto).GetMessage("CreateListResponse")
	assert.Len(m.GetMessages(), 1)
	assert.Nil(m.GetMessage("whodis"))

	n := m.GetMessage("Status")
	assert.Equal(n, m.GetMessage("CreateListResponse.Status"))
	assert.Equal("CreateListResponse.Status", n.GetName())
	assert.Equal("An internal status message", n.GetDescription())

	f := n.GetMessageField("code")
	assert.Equal("The status code.", f.GetDescription())
}
