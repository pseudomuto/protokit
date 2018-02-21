package protokit

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"

	"context"
	"fmt"
	"strings"

	"github.com/pseudomuto/protokit/utils"
)

const (
	// tag numbers in FileDescriptorProto
	packageCommentPath   = 2
	messageCommentPath   = 4
	enumCommentPath      = 5
	serviceCommentPath   = 6
	extensionCommentPath = 7

	// tag numbers in DescriptorProto
	messageFieldCommentPath     = 2 // field
	messageMessageCommentPath   = 3 // nested_type
	messageEnumCommentPath      = 4 // enum_type
	messageExtensionCommentPath = 6 // extension

	// tag numbers in EnumDescriptorProto
	enumValueCommentPath = 2 // value

	// tag numbers in ServiceDescriptorProto
	serviceMethodCommentPath = 2
)

// ParseCodeGenRequest parses the given request into `FileDescriptor` objects. Only the `req.FilesToGenerate` will be parsed
// here.
//
// For example, given the following invocation, only booking.proto will be parsed even if it imports other protos:
//
//     protoc --plugin=protoc-gen-test=./test -I. protos/booking.proto
func ParseCodeGenRequest(req *plugin_go.CodeGeneratorRequest) []*FileDescriptor {
	files := make([]*FileDescriptor, len(req.GetFileToGenerate()))

	for i, pf := range utils.FilesToGenerate(req) {
		files[i] = ParseFile(pf)
	}

	return files
}

// ParseFile parses a `FileDescriptorProto` into a `FileDescriptor` struct.
func ParseFile(fd *descriptor.FileDescriptorProto) *FileDescriptor {
	comments := ParseComments(fd)

	file := &FileDescriptor{
		comments:            comments,
		FileDescriptorProto: fd,
		Comments:            comments[fmt.Sprintf("%d", packageCommentPath)],
	}

	ctx := ContextWithFileDescriptor(context.Background(), file)
	file.Enums = parseEnums(ctx, fd.GetEnumType())
	file.Extensions = parseExtensions(ctx, fd.GetExtension())
	file.Messages = parseMessages(ctx, fd.GetMessageType())
	file.Services = parseServices(ctx, fd.GetService())

	return file
}

func parseEnums(ctx context.Context, protos []*descriptor.EnumDescriptorProto) []*EnumDescriptor {
	enums := make([]*EnumDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	parent, hasParent := DescriptorFromContext(ctx)

	for i, ed := range protos {
		longName := ed.GetName()
		commentPath := fmt.Sprintf("%d.%d", enumCommentPath, i)

		if hasParent {
			longName = fmt.Sprintf("%s.%s", parent.GetLongName(), longName)
			commentPath = fmt.Sprintf("%s.%d.%d", parent.path, messageEnumCommentPath, i)
		}

		enums[i] = &EnumDescriptor{
			common:              newCommon(file, commentPath, longName),
			EnumDescriptorProto: ed,
			Comments:            file.comments[commentPath],
			Parent:              parent,
		}

		subCtx := ContextWithEnumDescriptor(ctx, enums[i])
		enums[i].Values = parseEnumValues(subCtx, ed.GetValue())
	}

	return enums
}

func parseEnumValues(ctx context.Context, protos []*descriptor.EnumValueDescriptorProto) []*EnumValueDescriptor {
	values := make([]*EnumValueDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	enum, _ := EnumDescriptorFromContext(ctx)

	for i, vd := range protos {
		longName := fmt.Sprintf("%s.%s", enum.GetLongName(), vd.GetName())

		values[i] = &EnumValueDescriptor{
			common: newCommon(file, "", longName),
			EnumValueDescriptorProto: vd,
			Enum:     enum,
			Comments: file.comments[fmt.Sprintf("%s.%d.%d", enum.path, enumValueCommentPath, i)],
		}
	}

	return values
}

func parseExtensions(ctx context.Context, protos []*descriptor.FieldDescriptorProto) []*ExtensionDescriptor {
	exts := make([]*ExtensionDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	parent, hasParent := DescriptorFromContext(ctx)

	for i, ext := range protos {
		commentPath := fmt.Sprintf("%d.%d", extensionCommentPath, i)
		longName := fmt.Sprintf("%s.%s", ext.GetExtendee(), ext.GetName())

		if strings.Contains(longName, file.GetPackage()) {
			parts := strings.Split(ext.GetExtendee(), ".")
			longName = fmt.Sprintf("%s.%s", parts[len(parts)-1], ext.GetName())
		}

		if hasParent {
			commentPath = fmt.Sprintf("%s.%d.%d", parent.path, messageExtensionCommentPath, i)
		}

		exts[i] = &ExtensionDescriptor{
			common:               newCommon(file, commentPath, longName),
			FieldDescriptorProto: ext,
			Comments:             file.comments[commentPath],
			Parent:               parent,
		}
	}

	return exts
}

func parseMessages(ctx context.Context, protos []*descriptor.DescriptorProto) []*Descriptor {
	msgs := make([]*Descriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	parent, hasParent := DescriptorFromContext(ctx)

	for i, md := range protos {
		longName := md.GetName()
		commentPath := fmt.Sprintf("%d.%d", messageCommentPath, i)

		if hasParent {
			longName = fmt.Sprintf("%s.%s", parent.GetLongName(), longName)
			commentPath = fmt.Sprintf("%s.%d.%d", parent.path, messageMessageCommentPath, i)
		}

		msgs[i] = &Descriptor{
			common:          newCommon(file, commentPath, longName),
			DescriptorProto: md,
			Comments:        file.comments[commentPath],
			Parent:          parent,
		}

		msgCtx := ContextWithDescriptor(ctx, msgs[i])
		msgs[i].Enums = parseEnums(msgCtx, md.GetEnumType())
		msgs[i].Extensions = parseExtensions(msgCtx, md.GetExtension())
		msgs[i].Fields = parseMessageFields(msgCtx, md.GetField())
		msgs[i].Messages = parseMessages(msgCtx, md.GetNestedType())
	}

	return msgs
}

func parseMessageFields(ctx context.Context, protos []*descriptor.FieldDescriptorProto) []*FieldDescriptor {
	fields := make([]*FieldDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)
	message, _ := DescriptorFromContext(ctx)

	for i, fd := range protos {
		longName := fmt.Sprintf("%s.%s", message.GetLongName(), fd.GetName())

		fields[i] = &FieldDescriptor{
			common:               newCommon(file, "", longName),
			FieldDescriptorProto: fd,
			Comments:             file.comments[fmt.Sprintf("%s.%d.%d", message.path, messageFieldCommentPath, i)],
			Message:              message,
		}
	}

	return fields
}

func parseServices(ctx context.Context, protos []*descriptor.ServiceDescriptorProto) []*ServiceDescriptor {
	svcs := make([]*ServiceDescriptor, len(protos))
	file, _ := FileDescriptorFromContext(ctx)

	for i, sd := range protos {
		longName := sd.GetName()
		commentPath := fmt.Sprintf("%d.%d", serviceCommentPath, i)

		svcs[i] = &ServiceDescriptor{
			common:                 newCommon(file, commentPath, longName),
			ServiceDescriptorProto: sd,
			Comments:               file.comments[commentPath],
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
		longName := fmt.Sprintf("%s.%s", svc.GetLongName(), md.GetName())

		methods[i] = &MethodDescriptor{
			common:                newCommon(file, "", longName),
			MethodDescriptorProto: md,
			Service:               svc,
			Comments:              file.comments[fmt.Sprintf("%s.%d.%d", svc.path, serviceMethodCommentPath, i)],
		}
	}

	return methods
}
