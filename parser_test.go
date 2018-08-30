package protokit_test

import (
	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
)

var (
	proto2 *protokit.FileDescriptor
	proto3 *protokit.FileDescriptor

	noPackage *protokit.FileDescriptor
)

type ParserTest struct {
	suite.Suite
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTest))
}

func (assert *ParserTest) SetupSuite() {
	// File set with package set
	{
		set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
		assert.NoError(err)

		req := utils.CreateGenRequest(set, "booking.proto", "todo.proto")
		files := protokit.ParseCodeGenRequest(req)
		proto2 = files[0]
		proto3 = files[1]
	}

	// File set without package set
	{
		set, err := utils.LoadDescriptorSet("fixtures", "fileset_nopackage.pb")
		assert.NoError(err)

		req := utils.CreateGenRequest(set, "todo_nopackage.proto")
		files := protokit.ParseCodeGenRequest(req)
		noPackage = files[0]
	}
}

func (assert *ParserTest) TestFileParsing() {
	assert.True(proto3.IsProto3())
	assert.Equal("Top-level comments are attached to the syntax directive.", proto3.GetSyntaxComments().String())
	assert.Contains(proto3.GetPackageComments().String(), "The official documentation for the Todo API.\n\n")
	assert.Len(proto3.GetExtensions(), 0) // no extensions in proto3

	assert.False(proto2.IsProto3())
	assert.Len(proto2.GetExtensions(), 1)
}

func (assert *ParserTest) TestFileImports() {
	assert.Require().Len(proto3.GetImports(), 2)

	imp := proto3.GetImports()[0]
	assert.NotNil(imp.GetFile())
	assert.Equal("ListItemDetails", imp.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.ListItemDetails", imp.GetFullName())
}

func (assert *ParserTest) TestFileEnums() {
	assert.Len(proto3.GetEnums(), 1)
	assert.Nil(proto3.GetEnum("swingandamiss"))

	enum := proto3.GetEnum("ListType")
	assert.Equal("ListType", enum.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.ListType", enum.GetFullName())
	assert.True(enum.IsProto3())
	assert.Nil(enum.GetParent())
	assert.NotNil(enum.GetFile())
	assert.Equal("An enumeration of list types", enum.GetComments().String())
	assert.Equal("com.pseudomuto.protokit.v1", enum.GetPackage())
	assert.Len(enum.GetValues(), 2)

	assert.Equal("REMINDERS", enum.GetValues()[0].GetName())
	assert.Equal(enum, enum.GetValues()[0].GetEnum())
	assert.Equal("The reminders type.", enum.GetNamedValue("REMINDERS").GetComments().String())

	assert.Nil(enum.GetNamedValue("whodis"))
}

func (assert *ParserTest) TestFileExtensions() {
	ext := proto2.GetExtensions()[0]
	assert.Nil(ext.GetParent())
	assert.Equal("country", ext.GetName())
	assert.Equal("BookingStatus.country", ext.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.BookingStatus.country", ext.GetFullName())
	assert.Equal("The country the booking occurred in.", ext.GetComments().String())
}

func (assert *ParserTest) TestServices() {
	assert.Len(proto3.GetServices(), 1)
	assert.Nil(proto3.GetService("swingandamiss"))

	svc := proto3.GetService("Todo")
	assert.Equal("Todo", svc.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Todo", svc.GetFullName())
	assert.NotNil(svc.GetFile())
	assert.True(svc.IsProto3())
	assert.Contains(svc.GetComments().String(), "A service for managing \"todo\" items.\n\n")
	assert.Equal("com.pseudomuto.protokit.v1", svc.GetPackage())
	assert.Len(svc.GetMethods(), 2)

	m := svc.GetNamedMethod("CreateList")
	assert.Equal("CreateList", m.GetName())
	assert.Equal("Todo.CreateList", m.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Todo.CreateList", m.GetFullName())
	assert.NotNil(m.GetFile())
	assert.Equal(svc, m.GetService())
	assert.Equal("Create a new todo list", m.GetComments().String())

	m = svc.GetNamedMethod("Todo.AddItem")
	assert.Equal("Add an item to your list\n\nAdds a new item to the specified list.", m.GetComments().String())

	assert.Nil(svc.GetNamedMethod("wat"))
}

func (assert *ParserTest) TestFileMessages() {
	assert.Len(proto3.GetMessages(), 6)
	assert.Nil(proto3.GetMessage("swingandamiss"))

	m := proto3.GetMessage("AddItemRequest")
	assert.Equal("AddItemRequest", m.GetName())
	assert.Equal("AddItemRequest", m.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.AddItemRequest", m.GetFullName())
	assert.NotNil(m.GetFile())
	assert.Nil(m.GetParent())
	assert.Equal("A request message for adding new items.", m.GetComments().String())
	assert.Equal("com.pseudomuto.protokit.v1", m.GetPackage())
	assert.Len(m.GetMessageFields(), 3)
	assert.Nil(m.GetMessageField("swingandamiss"))

	// no extensions in proto3
	assert.Len(m.GetExtensions(), 0)

	f := m.GetMessageField("completed")
	assert.Equal("completed", f.GetName())
	assert.Equal("AddItemRequest.completed", f.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.AddItemRequest.completed", f.GetFullName())
	assert.NotNil(f.GetFile())
	assert.Equal(m, f.GetMessage())
	assert.Equal("Whether or not the item is completed.", f.GetComments().String())

	// just making sure google.protobuf.Any fields aren't special
	m = proto3.GetMessage("List")
	f = m.GetMessageField("details")
	assert.Equal("details", f.GetName())
	assert.Equal("List.details", f.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.List.details", f.GetFullName())

	// oneof fields should just expand to fields
	m = proto2.GetMessage("Booking")
	assert.NotNil(m.GetMessageField("reference_num"))
	assert.NotNil(m.GetMessageField("reference_tag"))
	assert.Equal("the numeric reference number", m.GetMessageField("reference_num").GetComments().String())
}

func (assert *ParserTest) TestMessageEnums() {
	m := proto3.GetMessage("Item")
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
	assert.Equal("An enumeration of possible statuses", e.GetComments().String())
	assert.Len(e.GetValues(), 2)

	val := e.GetNamedValue("COMPLETED")
	assert.Equal("COMPLETED", val.GetName())
	assert.Equal("Item.Status.COMPLETED", val.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.Item.Status.COMPLETED", val.GetFullName())
	assert.Equal("The completed status.", val.GetComments().String())
	assert.NotNil(val.GetFile())
}

func (assert *ParserTest) TestMessageExtensions() {
	m := proto2.GetMessage("Booking")
	ext := m.GetExtensions()[0]
	assert.Equal(m, ext.GetParent())
	assert.Equal(int32(101), ext.GetNumber())
	assert.Equal("optional_field_1", ext.GetName())
	assert.Equal("BookingStatus.optional_field_1", ext.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.BookingStatus.optional_field_1", ext.GetFullName())
	assert.Equal("An optional field to be used however you please.", ext.GetComments().String())
}

func (assert *ParserTest) TestNestedMessages() {
	m := proto3.GetMessage("CreateListResponse")
	assert.Len(m.GetMessages(), 1)
	assert.Nil(m.GetMessage("whodis"))

	n := m.GetMessage("Status")
	assert.Equal(n, m.GetMessage("CreateListResponse.Status"))

	// no extensions in proto3
	assert.Len(n.GetExtensions(), 0)

	assert.Equal("Status", n.GetName())
	assert.Equal("CreateListResponse.Status", n.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.CreateListResponse.Status", n.GetFullName())
	assert.Equal("An internal status message", n.GetComments().String())
	assert.NotNil(n.GetFile())
	assert.Equal(m, n.GetParent())

	f := n.GetMessageField("code")
	assert.Equal("CreateListResponse.Status.code", f.GetLongName())
	assert.Equal("com.pseudomuto.protokit.v1.CreateListResponse.Status.code", f.GetFullName())
	assert.NotNil(f.GetFile())
	assert.Equal("The status code.", f.GetComments().String())
}

func (assert *ParserTest) TestNoPackage() {
	assert.Equal("", noPackage.GetPackage())

	m := noPackage.GetMessage("Outer")
	assert.Equal("Outer", m.GetFullName())
}
