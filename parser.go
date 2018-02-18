package protokit

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"context"
	"fmt"
	"strings"
)

const (
	// tag numbers in FileDescriptorProto
	packageCommentPath = 2
	messageCommentPath = 4
	enumCommentPath    = 5
	serviceCommentPath = 6

	// tag numbers in DescriptorProto
	messageFieldCommentPath   = 2 // field
	messageMessageCommentPath = 3 // nested_type
	messageEnumCommentPath    = 4 // enum_type

	// tag numbers in EnumDescriptorProto
	enumValueCommentPath = 2 // value

	// tag numbers in ServiceDescriptorProto
	serviceMethodCommentPath = 2
)

// ParseFile parses a `FileDescriptorProto` into a `FileDescriptor` struct.
func ParseFile(fd *descriptor.FileDescriptorProto) *FileDescriptor {
	comments := ParseComments(fd)

	file := &FileDescriptor{
		comments:            comments,
		FileDescriptorProto: fd,
		Description:         comments[fmt.Sprintf("%d", packageCommentPath)],
	}

	ctx := ContextWithFileDescriptor(context.Background(), file)
	file.Enums = parseEnums(ctx, fd.GetEnumType())
	file.Messages = parseMessages(ctx, fd.GetMessageType())
	file.Services = parseServices(ctx, fd.GetService())

	return file
}

func parseEnums(ctx context.Context, protos []*descriptor.EnumDescriptorProto) []*EnumDescriptor {
	enums := make([]*EnumDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	parent, hasParent := DescriptorFromContext(ctx)

	for i, ed := range protos {
		commentPath := fmt.Sprintf("%d.%d", enumCommentPath, i)

		if hasParent {
			commentPath = fmt.Sprintf("%s.%d.%d", parent.path, messageEnumCommentPath, i)
		}

		enums[i] = &EnumDescriptor{
			common:              common{file: file, index: i, path: commentPath},
			EnumDescriptorProto: ed,
			Description:         file.comments[commentPath],
			Parent:              parent,
		}

		subCtx := ContextWithEnumDescriptor(ctx, enums[i])
		enums[i].Values = parseEnumValues(subCtx, ed.GetValue())

		if hasParent && !strings.Contains(enums[i].GetName(), ".") {
			enums[i].Name = proto.String(fmt.Sprintf("%s.%s", parent.GetName(), ed.GetName()))
		}
	}

	return enums
}

func parseEnumValues(ctx context.Context, protos []*descriptor.EnumValueDescriptorProto) []*EnumValueDescriptor {
	values := make([]*EnumValueDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	enum, _ := EnumDescriptorFromContext(ctx)

	for i, vd := range protos {
		values[i] = &EnumValueDescriptor{
			common: common{file: file, index: i},
			EnumValueDescriptorProto: vd,
			Enum:        enum,
			Description: file.comments[fmt.Sprintf("%s.%d.%d", enum.path, enumValueCommentPath, i)],
		}
	}

	return values
}

func parseMessages(ctx context.Context, protos []*descriptor.DescriptorProto) []*Descriptor {
	msgs := make([]*Descriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	parent, hasParent := DescriptorFromContext(ctx)

	for i, md := range protos {
		commentPath := fmt.Sprintf("%d.%d", messageCommentPath, i)
		if hasParent {
			commentPath = fmt.Sprintf("%s.%d.%d", parent.path, messageMessageCommentPath, i)
		}

		msgs[i] = &Descriptor{
			common:          common{file: file, index: i, path: commentPath},
			DescriptorProto: md,
			Description:     file.comments[commentPath],
			Parent:          parent,
		}

		msgCtx := ContextWithDescriptor(ctx, msgs[i])
		msgs[i].Fields = parseMessageFields(msgCtx, md.GetField())
		msgs[i].Messages = parseMessages(msgCtx, md.GetNestedType())
		msgs[i].Enums = parseEnums(msgCtx, md.GetEnumType())

		if hasParent && !strings.Contains(msgs[i].GetName(), ".") {
			msgs[i].Name = proto.String(fmt.Sprintf("%s.%s", parent.GetName(), md.GetName()))
		}
	}

	return msgs
}

func parseMessageFields(ctx context.Context, protos []*descriptor.FieldDescriptorProto) []*FieldDescriptor {
	fields := make([]*FieldDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	message, _ := DescriptorFromContext(ctx)

	for i, fd := range protos {
		fields[i] = &FieldDescriptor{
			common:               common{file: file, index: i},
			FieldDescriptorProto: fd,
			Description:          file.comments[fmt.Sprintf("%s.%d.%d", message.path, messageFieldCommentPath, i)],
			Message:              message,
		}
	}

	return fields
}

func parseServices(ctx context.Context, protos []*descriptor.ServiceDescriptorProto) []*ServiceDescriptor {
	svcs := make([]*ServiceDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)

	for i, sd := range protos {
		commentPath := fmt.Sprintf("%d.%d", serviceCommentPath, i)

		svcs[i] = &ServiceDescriptor{
			common:                 common{file: file, index: i, path: commentPath},
			ServiceDescriptorProto: sd,
			Description:            file.comments[commentPath],
		}

		svcCtx := ContextWithServiceDescriptor(ctx, svcs[i])
		svcs[i].Methods = parseServiceMethods(svcCtx, sd.GetMethod())
	}

	return svcs
}

func parseServiceMethods(ctx context.Context, protos []*descriptor.MethodDescriptorProto) []*MethodDescriptor {
	methods := make([]*MethodDescriptor, len(protos))

	file, _ := FileDescriptorFromContext(ctx)
	svc, _ := ServiceDescriptorFromContext(ctx)

	for i, md := range protos {
		methods[i] = &MethodDescriptor{
			common:                common{file: file, index: i},
			MethodDescriptorProto: md,
			Service:               svc,
			Description:           file.comments[fmt.Sprintf("%s.%d.%d", svc.path, serviceMethodCommentPath, i)],
		}
	}

	return methods
}
