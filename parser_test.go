package protokit_test

import (
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	"github.com/stretchr/testify/require"
)

func setupParserTest(t *testing.T) (*protokit.FileDescriptor, *protokit.FileDescriptor) {
	t.Helper()

	set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(set, "booking.proto", "todo.proto")
	files := protokit.ParseCodeGenRequest(req)
	proto2 := files[0]
	proto3 := files[1]

	return proto2, proto3
}

func TestFileParsing(t *testing.T) {
	t.Parallel()

	proto2, proto3 := setupParserTest(t)

	require.True(t, proto3.IsProto3())
	require.Equal(t, "Top-level comments are attached to the syntax directive.", proto3.GetSyntaxComments().String())
	require.Contains(t, proto3.GetPackageComments().String(), "The official documentation for the Todo API.\n\n")
	require.Empty(t, proto3.GetExtensions()) // no extensions in proto3

	require.False(t, proto2.IsProto3())
	require.Len(t, proto2.GetExtensions(), 1)
}

func TestFileImports(t *testing.T) {
	t.Parallel()

	_, proto3 := setupParserTest(t)

	require.Len(t, proto3.GetImports(), 2)

	imp := proto3.GetImports()[0]
	require.NotNil(t, imp.GetFile())
	require.Equal(t, "ListItemDetails", imp.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.ListItemDetails", imp.GetFullName())
}

func TestFileEnums(t *testing.T) {
	t.Parallel()

	_, proto3 := setupParserTest(t)

	require.Len(t, proto3.GetEnums(), 1)
	require.Nil(t, proto3.GetEnum("swingandamiss"))

	enum := proto3.GetEnum("ListType")
	require.Equal(t, "ListType", enum.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.ListType", enum.GetFullName())
	require.True(t, enum.IsProto3())
	require.Nil(t, enum.GetParent())
	require.NotNil(t, enum.GetFile())
	require.Equal(t, "An enumeration of list types", enum.GetComments().String())
	require.Equal(t, "com.pseudomuto.protokit.v1", enum.GetPackage())
	require.Len(t, enum.GetValues(), 2)

	require.Equal(t, "REMINDERS", enum.GetValues()[0].GetName())
	require.Equal(t, enum, enum.GetValues()[0].GetEnum())
	require.Equal(t, "The reminders type.", enum.GetNamedValue("REMINDERS").GetComments().String())

	require.Nil(t, enum.GetNamedValue("whodis"))
}

func TestFileExtensions(t *testing.T) {
	t.Parallel()

	proto2, _ := setupParserTest(t)

	ext := proto2.GetExtensions()[0]
	require.Nil(t, ext.GetParent())
	require.Equal(t, "country", ext.GetName())
	require.Equal(t, "BookingStatus.country", ext.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.BookingStatus.country", ext.GetFullName())
	require.Equal(t, "The country the booking occurred in.", ext.GetComments().String())
}

func TestServices(t *testing.T) {
	t.Parallel()

	_, proto3 := setupParserTest(t)

	require.Len(t, proto3.GetServices(), 1)
	require.Nil(t, proto3.GetService("swingandamiss"))

	svc := proto3.GetService("Todo")
	require.Equal(t, "Todo", svc.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.Todo", svc.GetFullName())
	require.NotNil(t, svc.GetFile())
	require.True(t, svc.IsProto3())
	require.Contains(t, svc.GetComments().String(), "A service for managing \"todo\" items.\n\n")
	require.Equal(t, "com.pseudomuto.protokit.v1", svc.GetPackage())
	require.Len(t, svc.GetMethods(), 2)

	m := svc.GetNamedMethod("CreateList")
	require.Equal(t, "CreateList", m.GetName())
	require.Equal(t, "Todo.CreateList", m.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.Todo.CreateList", m.GetFullName())
	require.NotNil(t, m.GetFile())
	require.Equal(t, svc, m.GetService())
	require.Equal(t, "Create a new todo list", m.GetComments().String())

	m = svc.GetNamedMethod("Todo.AddItem")
	require.Equal(t, "Add an item to your list\n\nAdds a new item to the specified list.", m.GetComments().String())

	require.Nil(t, svc.GetNamedMethod("wat"))
}

func TestFileMessages(t *testing.T) {
	t.Parallel()

	proto2, proto3 := setupParserTest(t)

	require.Len(t, proto3.GetMessages(), 6)
	require.Nil(t, proto3.GetMessage("swingandamiss"))

	m := proto3.GetMessage("AddItemRequest")
	require.Equal(t, "AddItemRequest", m.GetName())
	require.Equal(t, "AddItemRequest", m.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.AddItemRequest", m.GetFullName())
	require.NotNil(t, m.GetFile())
	require.Nil(t, m.GetParent())
	require.Equal(t, "A request message for adding new items.", m.GetComments().String())
	require.Equal(t, "com.pseudomuto.protokit.v1", m.GetPackage())
	require.Len(t, m.GetMessageFields(), 3)
	require.Nil(t, m.GetMessageField("swingandamiss"))

	// no extensions in proto3
	require.Empty(t, m.GetExtensions())

	f := m.GetMessageField("completed")
	require.Equal(t, "completed", f.GetName())
	require.Equal(t, "AddItemRequest.completed", f.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.AddItemRequest.completed", f.GetFullName())
	require.NotNil(t, f.GetFile())
	require.Equal(t, m, f.GetMessage())
	require.Equal(t, "Whether or not the item is completed.", f.GetComments().String())

	// just making sure google.protobuf.Any fields aren't special
	m = proto3.GetMessage("List")
	f = m.GetMessageField("details")
	require.Equal(t, "details", f.GetName())
	require.Equal(t, "List.details", f.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.List.details", f.GetFullName())

	// oneof fields should just expand to fields
	m = proto2.GetMessage("Booking")
	require.NotNil(t, m.GetMessageField("reference_num"))
	require.NotNil(t, m.GetMessageField("reference_tag"))
	require.Equal(t, "the numeric reference number", m.GetMessageField("reference_num").GetComments().String())
}

func TestMessageEnums(t *testing.T) {
	t.Parallel()

	_, proto3 := setupParserTest(t)

	m := proto3.GetMessage("Item")
	require.NotNil(t, m.GetFile())
	require.Len(t, m.GetEnums(), 1)
	require.Nil(t, m.GetEnum("whodis"))

	e := m.GetEnum("Status")
	require.Equal(t, "Status", e.GetName())
	require.Equal(t, "Item.Status", e.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.Item.Status", e.GetFullName())
	require.NotNil(t, e.GetFile())
	require.Equal(t, m, e.GetParent())
	require.Equal(t, e, m.GetEnum("Item.Status"))
	require.Equal(t, "An enumeration of possible statuses", e.GetComments().String())
	require.Len(t, e.GetValues(), 2)

	val := e.GetNamedValue("COMPLETED")
	require.Equal(t, "COMPLETED", val.GetName())
	require.Equal(t, "Item.Status.COMPLETED", val.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.Item.Status.COMPLETED", val.GetFullName())
	require.Equal(t, "The completed status.", val.GetComments().String())
	require.NotNil(t, val.GetFile())
}

func TestMessageExtensions(t *testing.T) {
	t.Parallel()

	proto2, _ := setupParserTest(t)

	m := proto2.GetMessage("Booking")
	ext := m.GetExtensions()[0]
	require.Equal(t, m, ext.GetParent())
	require.Equal(t, int32(101), ext.GetNumber())
	require.Equal(t, "optional_field_1", ext.GetName())
	require.Equal(t, "BookingStatus.optional_field_1", ext.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.BookingStatus.optional_field_1", ext.GetFullName())
	require.Equal(t, "An optional field to be used however you please.", ext.GetComments().String())
}

func TestNestedMessages(t *testing.T) {
	t.Parallel()

	_, proto3 := setupParserTest(t)

	m := proto3.GetMessage("CreateListResponse")
	require.Len(t, m.GetMessages(), 1)
	require.Nil(t, m.GetMessage("whodis"))

	n := m.GetMessage("Status")
	require.Equal(t, n, m.GetMessage("CreateListResponse.Status"))

	// no extensions in proto3
	require.Empty(t, n.GetExtensions())

	require.Equal(t, "Status", n.GetName())
	require.Equal(t, "CreateListResponse.Status", n.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.CreateListResponse.Status", n.GetFullName())
	require.Equal(t, "An internal status message", n.GetComments().String())
	require.NotNil(t, n.GetFile())
	require.Equal(t, m, n.GetParent())

	f := n.GetMessageField("code")
	require.Equal(t, "CreateListResponse.Status.code", f.GetLongName())
	require.Equal(t, "com.pseudomuto.protokit.v1.CreateListResponse.Status.code", f.GetFullName())
	require.NotNil(t, f.GetFile())
	require.Equal(t, "The status code.", f.GetComments().String())
}

func TestExtendedOptions(t *testing.T) {
	t.Parallel()

	proto2, _ := setupParserTest(t)

	require.Contains(t, proto2.OptionExtensions, "com.pseudomuto.protokit.v1.extend_file")

	extendedValue, ok := proto2.OptionExtensions["com.pseudomuto.protokit.v1.extend_file"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	service := proto2.GetService("BookingService")
	require.Contains(t, service.OptionExtensions, "com.pseudomuto.protokit.v1.extend_service")

	extendedValue, ok = service.OptionExtensions["com.pseudomuto.protokit.v1.extend_service"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	method := service.GetNamedMethod("BookVehicle")
	require.Contains(t, method.OptionExtensions, "com.pseudomuto.protokit.v1.extend_method")

	extendedValue, ok = method.OptionExtensions["com.pseudomuto.protokit.v1.extend_method"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	message := proto2.GetMessage("Booking")
	require.Contains(t, message.OptionExtensions, "com.pseudomuto.protokit.v1.extend_message")

	extendedValue, ok = message.OptionExtensions["com.pseudomuto.protokit.v1.extend_message"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	field := message.GetMessageField("payment_received")
	require.Contains(t, field.OptionExtensions, "com.pseudomuto.protokit.v1.extend_field")

	extendedValue, ok = field.OptionExtensions["com.pseudomuto.protokit.v1.extend_field"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	enum := proto2.GetEnum("BookingType")
	require.Contains(t, enum.OptionExtensions, "com.pseudomuto.protokit.v1.extend_enum")

	extendedValue, ok = enum.OptionExtensions["com.pseudomuto.protokit.v1.extend_enum"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	enumValue := enum.GetNamedValue("FUTURE")
	require.Contains(t, enumValue.OptionExtensions, "com.pseudomuto.protokit.v1.extend_enum_value")

	extendedValue, ok = enumValue.OptionExtensions["com.pseudomuto.protokit.v1.extend_enum_value"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	_, proto3 := setupParserTest(t)
	require.Contains(t, proto3.OptionExtensions, "com.pseudomuto.protokit.v1.extend_file")

	extendedValue, ok = proto3.OptionExtensions["com.pseudomuto.protokit.v1.extend_file"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	service = proto3.GetService("Todo")
	require.Contains(t, service.OptionExtensions, "com.pseudomuto.protokit.v1.extend_service")

	extendedValue, ok = service.OptionExtensions["com.pseudomuto.protokit.v1.extend_service"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	method = service.GetNamedMethod("CreateList")
	require.Contains(t, method.OptionExtensions, "com.pseudomuto.protokit.v1.extend_method")

	extendedValue, ok = method.OptionExtensions["com.pseudomuto.protokit.v1.extend_method"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	message = proto3.GetMessage("List")
	require.Contains(t, message.OptionExtensions, "com.pseudomuto.protokit.v1.extend_message")

	extendedValue, ok = message.OptionExtensions["com.pseudomuto.protokit.v1.extend_message"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	field = message.GetMessageField("name")
	require.Contains(t, field.OptionExtensions, "com.pseudomuto.protokit.v1.extend_field")

	extendedValue, ok = field.OptionExtensions["com.pseudomuto.protokit.v1.extend_field"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	enum = proto3.GetEnum("ListType")
	require.Contains(t, enum.OptionExtensions, "com.pseudomuto.protokit.v1.extend_enum")

	extendedValue, ok = enum.OptionExtensions["com.pseudomuto.protokit.v1.extend_enum"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)

	enumValue = enum.GetNamedValue("CHECKLIST")
	require.Contains(t, enumValue.OptionExtensions, "com.pseudomuto.protokit.v1.extend_enum_value")

	extendedValue, ok = enumValue.OptionExtensions["com.pseudomuto.protokit.v1.extend_enum_value"].(*bool)
	require.True(t, ok)
	require.True(t, *extendedValue)
}
