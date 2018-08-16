package protokit_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func registerTestExtensions() {
	var E_ExtendService = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.ServiceOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_service",
		Tag:           "varint,20000,opt,name=extend_service,json=extendService",
		Filename:      "booking.proto",
	}

	var E_ExtendMethod = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_method",
		Tag:           "varint,20000,opt,name=extend_method,json=extendMethod",
		Filename:      "booking.proto",
	}

	var E_ExtendEnum = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.EnumOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_enum",
		Tag:           "varint,20000,opt,name=extend_enum,json=extendEnum",
		Filename:      "booking.proto",
	}

	var E_ExtendEnumValue = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.EnumValueOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_enum_value",
		Tag:           "varint,20000,opt,name=extend_enum_value,json=extendEnumValue",
		Filename:      "booking.proto",
	}

	var E_ExtendMessage = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_message",
		Tag:           "varint,20000,opt,name=extend_message,json=extendMessage",
		Filename:      "booking.proto",
	}

	var E_ExtendField = &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         20000,
		Name:          "com.pseudomuto.protokit.v1.extend_field",
		Tag:           "varint,20000,opt,name=extend_field,json=extendField",
		Filename:      "booking.proto",
	}

	proto.RegisterExtension(E_ExtendService)
	proto.RegisterExtension(E_ExtendMethod)
	proto.RegisterExtension(E_ExtendEnum)
	proto.RegisterExtension(E_ExtendEnumValue)
	proto.RegisterExtension(E_ExtendMessage)
	proto.RegisterExtension(E_ExtendField)
}
