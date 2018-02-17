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

// ParseFile parses a `FileDescriptorProto` into a `File` struct.
func ParseFile(fd *descriptor.FileDescriptorProto) *FileDescriptor {
	comments := ParseComments(fd)
	ctx := ContextWithComments(context.Background(), comments)
	ctx = ContextWithPackage(ctx, fd.GetPackage())

	file := &FileDescriptor{
		FileDescriptorProto: fd,
		Description:         comments[fmt.Sprintf("%d", packageCommentPath)],
	}

	file.Enums = parseEnums(ctx, fd.GetEnumType())
	file.Messages = parseMessages(ctx, fd.GetMessageType())
	file.Services = parseServices(ctx, fd.GetService())

	return file
}

func parseEnums(ctx context.Context, protos []*descriptor.EnumDescriptorProto) []*EnumDescriptor {
	enums := make([]*EnumDescriptor, len(protos))
	comments, _ := CommentsFromContext(ctx)
	pkg, _ := PackageFromContext(ctx)
	commentPrefix, hasPrefix := LocationPrefixFromContext(ctx)
	message, hasMessage := MessageFromContext(ctx)

	for i, ed := range protos {
		commentPath := fmt.Sprintf("%d.%d", enumCommentPath, i)

		if hasPrefix {
			commentPath = fmt.Sprintf("%s.%d", commentPrefix, i)
		}

		subCtx := ContextWithLocationPrefix(ctx, commentPath)

		enums[i] = &EnumDescriptor{
			EnumDescriptorProto: ed,
			Values:              parseEnumValues(subCtx, ed.GetValue()),
			Description:         comments[commentPath],
			Package:             pkg,
		}

		if hasMessage && !strings.Contains(enums[i].GetName(), ".") {
			enums[i].Name = proto.String(fmt.Sprintf("%s.%s", message, ed.GetName()))
		}
	}

	return enums
}

func parseEnumValues(ctx context.Context, protos []*descriptor.EnumValueDescriptorProto) []*EnumValueDescriptor {
	values := make([]*EnumValueDescriptor, len(protos))
	comments, _ := CommentsFromContext(ctx)
	commentPrefix, _ := LocationPrefixFromContext(ctx)

	for i, vd := range protos {
		values[i] = &EnumValueDescriptor{
			EnumValueDescriptorProto: vd,
			Description:              comments[fmt.Sprintf("%s.%d.%d", commentPrefix, enumValueCommentPath, i)],
		}
	}

	return values
}

func parseMessages(ctx context.Context, protos []*descriptor.DescriptorProto) []*Descriptor {
	msgs := make([]*Descriptor, len(protos))
	comments, _ := CommentsFromContext(ctx)
	commentPrefix, hasPrefix := LocationPrefixFromContext(ctx)
	message, hasMessage := MessageFromContext(ctx)
	pkg, _ := PackageFromContext(ctx)

	for i, md := range protos {
		commentPath := fmt.Sprintf("%d.%d", messageCommentPath, i)
		if hasPrefix {
			commentPath = fmt.Sprintf("%s.%d.%d", commentPrefix, messageMessageCommentPath, i)
		}

		enumPath := fmt.Sprintf("%s.%d", commentPath, messageEnumCommentPath)
		enumCtx := ContextWithMessage(ContextWithLocationPrefix(ctx, enumPath), md.GetName())
		msgCtx := ContextWithMessage(ContextWithLocationPrefix(ctx, commentPath), md.GetName())

		msgs[i] = &Descriptor{
			DescriptorProto: md,
			Description:     comments[commentPath],
			Enums:           parseEnums(enumCtx, md.GetEnumType()),
			Fields:          parseMessageFields(msgCtx, md.GetField()),
			Messages:        parseMessages(msgCtx, md.GetNestedType()),
			Package:         pkg,
		}

		if hasMessage && !strings.Contains(msgs[i].GetName(), ".") {
			msgs[i].Name = proto.String(fmt.Sprintf("%s.%s", message, md.GetName()))
		}
	}

	return msgs
}

func parseMessageFields(ctx context.Context, protos []*descriptor.FieldDescriptorProto) []*FieldDescriptor {
	fields := make([]*FieldDescriptor, len(protos))
	comments, _ := CommentsFromContext(ctx)
	commentPrefix, _ := LocationPrefixFromContext(ctx)

	for i, fd := range protos {
		fields[i] = &FieldDescriptor{
			FieldDescriptorProto: fd,
			Description:          comments[fmt.Sprintf("%s.%d.%d", commentPrefix, messageFieldCommentPath, i)],
		}
	}

	return fields
}

func parseServices(ctx context.Context, protos []*descriptor.ServiceDescriptorProto) []*ServiceDescriptor {
	svcs := make([]*ServiceDescriptor, len(protos))
	comments, _ := CommentsFromContext(ctx)
	pkg, _ := PackageFromContext(ctx)

	for i, sd := range protos {
		commentPath := fmt.Sprintf("%d.%d", serviceCommentPath, i)
		subCtx := ContextWithLocationPrefix(ctx, commentPath)
		subCtx = ContextWithService(subCtx, sd.GetName())

		svcs[i] = &ServiceDescriptor{
			ServiceDescriptorProto: sd,
			Description:            comments[commentPath],
			Methods:                parseServiceMethods(subCtx, sd.GetMethod()),
			Package:                pkg,
		}
	}

	return svcs
}

func parseServiceMethods(ctx context.Context, protos []*descriptor.MethodDescriptorProto) []*MethodDescriptor {
	methods := make([]*MethodDescriptor, len(protos))

	pkg, _ := PackageFromContext(ctx)
	svc, _ := ServiceFromContext(ctx)
	comments, _ := CommentsFromContext(ctx)
	commentPrefix, _ := LocationPrefixFromContext(ctx)

	for i, md := range protos {
		methods[i] = &MethodDescriptor{
			MethodDescriptorProto: md,
			Description:           comments[fmt.Sprintf("%s.%d.%d", commentPrefix, serviceMethodCommentPath, i)],
			URL:                   fmt.Sprintf("/%s.%s/%s", pkg, svc, md.GetName()),
			InputRef:              typeRef(md.GetInputType()),
			OutputRef:             typeRef(md.GetOutputType()),
		}
	}

	return methods
}

func typeRef(typeStr string) *TypeReference {
	parts := strings.Split(typeStr, ".")
	pkg := ""

	if len(parts) > 2 {
		pkg = strings.Join(parts[1:len(parts)-1], ".")
	}

	return &TypeReference{
		Package:        pkg,
		TypeName:       parts[len(parts)-1],
		FullyQualified: parts[0] == "",
	}
}
