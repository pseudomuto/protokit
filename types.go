package protokit

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"fmt"
)

type common struct {
	file *FileDescriptor
}

func (c *common) GetFile() *FileDescriptor { return c.file }
func (c *common) GetPackage() string       { return c.file.GetPackage() }
func (c *common) IsProto3() bool           { return c.file.GetSyntax() == "proto3" }

// A TypeReference represents a reference to a type. It includes the package, name, and whether or not the reference is
// fully qualified (name starts with a ".").
type TypeReference struct {
	Package        string // The package name (if available)
	TypeName       string // The name (without package)
	FullyQualified bool   // Whether or not the reference is full qualified
}

// GetPackage returns the package name (if available)
func (tr *TypeReference) GetPackage() string { return tr.Package }

// GetTypeName returns the name of the type (without the package)
func (tr *TypeReference) GetTypeName() string { return tr.TypeName }

// GetFullyQualified returns whether or not the type if fully-qualified
func (tr *TypeReference) GetFullyQualified() bool { return tr.FullyQualified }

// A FileDescriptor describes a single proto file with all of its messages, enums, services, etc.
type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	Description string
	Enums       []*EnumDescriptor
	Messages    []*Descriptor
	Services    []*ServiceDescriptor
}

// IsProto3 returns whether or not this file is a proto3 file
func (f *FileDescriptor) IsProto3() bool { return f.GetSyntax() == "proto3" }

// GetDescription returns the file's package comments
func (f *FileDescriptor) GetDescription() string { return f.Description }

// GetEnums returns the top-level enumerations defined in this file
func (f *FileDescriptor) GetEnums() []*EnumDescriptor { return f.Enums }

// GetMessages returns the top-level messages defined in this file
func (f *FileDescriptor) GetMessages() []*Descriptor { return f.Messages }

// GetServices returns the services defined in this file
func (f *FileDescriptor) GetServices() []*ServiceDescriptor { return f.Services }

// GetEnum returns the enumeration with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetEnum(name string) *EnumDescriptor {
	for _, e := range f.GetEnums() {
		if e.GetName() == name {
			return e
		}
	}

	return nil
}

// GetMessage returns the message with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetMessage(name string) *Descriptor {
	for _, m := range f.GetMessages() {
		if m.GetName() == name {
			return m
		}
	}

	return nil
}

// GetService returns the service with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetService(name string) *ServiceDescriptor {
	for _, s := range f.GetServices() {
		if s.GetName() == name {
			return s
		}
	}

	return nil
}

// An EnumDescriptor describe an enum type
type EnumDescriptor struct {
	common
	*descriptor.EnumDescriptorProto
	Values      []*EnumValueDescriptor
	Description string
}

// GetDescription returns a description of this enum
func (e *EnumDescriptor) GetDescription() string { return e.Description }

// GetValues returns the available values for this enum
func (e *EnumDescriptor) GetValues() []*EnumValueDescriptor { return e.Values }

// GetNamedValue returns the value with the specified name (returns `nil` if not found)
func (e *EnumDescriptor) GetNamedValue(name string) *EnumValueDescriptor {
	for _, v := range e.GetValues() {
		if v.GetName() == name {
			return v
		}
	}

	return nil
}

// An EnumValueDescriptor describes an enum value
type EnumValueDescriptor struct {
	common
	*descriptor.EnumValueDescriptorProto
	Description string
}

// GetDescription returns a description of the value
func (v *EnumValueDescriptor) GetDescription() string { return v.Description }

// A Descriptor describes a message
type Descriptor struct {
	common
	*descriptor.DescriptorProto
	Description string
	Enums       []*EnumDescriptor
	Fields      []*FieldDescriptor
	Messages    []*Descriptor
}

// GetDescription returns a description of the message
func (m *Descriptor) GetDescription() string { return m.Description }

// GetEnums returns the nested enumerations within the message
func (m *Descriptor) GetEnums() []*EnumDescriptor { return m.Enums }

// GetMessages returns the nested messages within the message
func (m *Descriptor) GetMessages() []*Descriptor { return m.Messages }

// GetMessageFields returns the message fields
func (m *Descriptor) GetMessageFields() []*FieldDescriptor { return m.Fields }

// GetEnum returns the enum with the specified name. The name can be either simple, or fully qualified (returns `nil` if
// not found)
func (m *Descriptor) GetEnum(name string) *EnumDescriptor {
	qn := fmt.Sprintf("%s.%s", m.GetName(), name)

	for _, e := range m.GetEnums() {
		// can lookup by name or message prefixed name (qualified)
		if e.GetName() == name || e.GetName() == qn {
			return e
		}
	}

	return nil
}

// GetMessage returns the nested message with the specified name. The name can be simple or fully qualified (returns
// `nil` if not found)
func (m *Descriptor) GetMessage(name string) *Descriptor {
	qn := fmt.Sprintf("%s.%s", m.GetName(), name)

	for _, msg := range m.GetMessages() {
		// can lookup by name or message prefixed name (qualified)
		if msg.GetName() == name || msg.GetName() == qn {
			return msg
		}
	}

	return nil
}

// GetMessageField returns the field with the specified name (returns `nil` if not found)
func (m *Descriptor) GetMessageField(name string) *FieldDescriptor {
	for _, f := range m.GetMessageFields() {
		if f.GetName() == name {
			return f
		}
	}

	return nil
}

// A FieldDescriptor describes a message field
type FieldDescriptor struct {
	common
	*descriptor.FieldDescriptorProto
	Description string
}

// GetDescription returns a description of the field
func (mf *FieldDescriptor) GetDescription() string { return mf.Description }

// A ServiceDescriptor describes a service
type ServiceDescriptor struct {
	common
	*descriptor.ServiceDescriptorProto
	Description string
	Methods     []*MethodDescriptor
}

// GetDescription returns a description of the service
func (s *ServiceDescriptor) GetDescription() string { return s.Description }

// GetMethods returns the methods for the service
func (s *ServiceDescriptor) GetMethods() []*MethodDescriptor { return s.Methods }

// A MethodDescriptor describes a method in a service
type MethodDescriptor struct {
	common
	*descriptor.MethodDescriptorProto
	InputRef    *TypeReference
	OutputRef   *TypeReference
	Description string
	URL         string
}

// GetDescription returns a description of the method
func (m *MethodDescriptor) GetDescription() string { return m.Description }

// GetURL returns the URL for the method
func (m *MethodDescriptor) GetURL() string { return m.URL }

// GetInputRef returns a reference to the input type
func (m *MethodDescriptor) GetInputRef() *TypeReference { return m.InputRef }

// GetOutputRef returns a reference to the output type
func (m *MethodDescriptor) GetOutputRef() *TypeReference { return m.OutputRef }
