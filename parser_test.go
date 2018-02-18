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
	assert.True(file.IsProto3())
	assert.Contains(file.GetDescription(), "The official documentation for the Todo API.\n\n")
}

func (assert *ParserTest) TestParseFileEnums() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetEnums(), 1)
	assert.Nil(file.GetEnum("swingandamiss"))

	enum := file.GetEnum("ListType")
	assert.Equal("ListType", enum.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.ListType", enum.GetFullName())
	assert.True(enum.IsProto3())
	assert.Nil(enum.GetParent())
	assert.NotNil(enum.GetFile())
	assert.Equal("An enumeration of list types", enum.GetDescription())
	assert.Equal("com.pseudomuto.protokit.v1", enum.GetPackage())
	assert.Len(enum.GetValues(), 2)

	assert.Equal("REMINDERS", enum.GetValues()[0].GetName())
	assert.Equal(enum, enum.GetValues()[0].GetEnum())
	assert.Equal("The reminders type.", enum.GetNamedValue("REMINDERS").GetDescription())

	assert.Nil(enum.GetNamedValue("whodis"))
}

func (assert *ParserTest) TestParseFileServices() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetServices(), 1)
	assert.Nil(file.GetService("swingandamiss"))

	svc := file.GetService("Todo")
	assert.Equal("Todo", svc.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Todo", svc.GetFullName())
	assert.NotNil(svc.GetFile())
	assert.True(svc.IsProto3())
	assert.Contains(svc.GetDescription(), "A service for managing \"todo\" items.\n\n")
	assert.Equal("com.pseudomuto.protokit.v1", svc.GetPackage())
	assert.Len(svc.GetMethods(), 2)

	m := svc.GetNamedMethod("CreateList")
	assert.Equal("CreateList", m.GetName())
	assert.Equal("Todo.CreateList", m.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Todo.CreateList", m.GetFullName())
	assert.NotNil(m.GetFile())
	assert.Equal(svc, m.GetService())
	assert.Equal("Create a new todo list", m.GetDescription())

	m = svc.GetNamedMethod("Todo.AddItem")
	assert.Equal("Add an item to your list\n\nAdds a new item to the specified list.", m.GetDescription())

	assert.Nil(svc.GetNamedMethod("wat"))
}

func (assert *ParserTest) TestParseFileMessages() {
	file := protokit.ParseFile(proto)
	assert.Len(file.GetMessages(), 6)
	assert.Nil(file.GetMessage("swingandamiss"))

	m := file.GetMessage("AddItemRequest")
	assert.Equal("AddItemRequest", m.GetName())
	assert.Equal("AddItemRequest", m.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.AddItemRequest", m.GetFullName())
	assert.NotNil(m.GetFile())
	assert.Nil(m.GetParent())
	assert.Equal("A request message for adding new items.", m.GetDescription())
	assert.Equal("com.pseudomuto.protokit.v1", m.GetPackage())
	assert.Len(m.GetMessageFields(), 3)
	assert.Nil(m.GetMessageField("swingandamiss"))

	f := m.GetMessageField("completed")
	assert.Equal("completed", f.GetName())
	assert.Equal("AddItemRequest.completed", f.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.AddItemRequest.completed", f.GetFullName())
	assert.NotNil(f.GetFile())
	assert.Equal(m, f.GetMessage())
	assert.Equal("Whether or not the item is completed.", f.GetDescription())
}

func (assert *ParserTest) TestParseFileMessageEnums() {
	m := protokit.ParseFile(proto).GetMessage("Item")
	assert.NotNil(m.GetFile())
	assert.Len(m.GetEnums(), 1)
	assert.Nil(m.GetEnum("whodis"))

	e := m.GetEnum("Status")
	assert.Equal("Status", e.GetName())
	assert.Equal("Item.Status", e.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Item.Status", e.GetFullName())
	assert.NotNil(e.GetFile())
	assert.Equal(m, e.GetParent())
	assert.Equal(e, m.GetEnum("Item.Status"))
	assert.Equal("An enumeration of possible statuses", e.GetDescription())
	assert.Len(e.GetValues(), 2)

	val := e.GetNamedValue("COMPLETED")
	assert.Equal("COMPLETED", val.GetName())
	assert.Equal("Item.Status.COMPLETED", val.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Item.Status.COMPLETED", val.GetFullName())
	assert.Equal("The completed status.", val.GetDescription())
	assert.NotNil(val.GetFile())
}

func (assert *ParserTest) TestParseFileNestedMessages() {
	m := protokit.ParseFile(proto).GetMessage("CreateListResponse")
	assert.Len(m.GetMessages(), 1)
	assert.Nil(m.GetMessage("whodis"))

	n := m.GetMessage("Status")
	assert.Equal(n, m.GetMessage("CreateListResponse.Status"))

	assert.Equal("Status", n.GetName())
	assert.Equal("CreateListResponse.Status", n.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.CreateListResponse.Status", n.GetFullName())
	assert.Equal("An internal status message", n.GetDescription())
	assert.NotNil(n.GetFile())
	assert.Equal(m, n.GetParent())

	f := n.GetMessageField("code")
	assert.Equal("CreateListResponse.Status.code", f.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.CreateListResponse.Status.code", f.GetFullName())
	assert.NotNil(f.GetFile())
	assert.Equal("The status code.", f.GetDescription())
}
