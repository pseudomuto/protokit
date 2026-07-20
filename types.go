package protokit

import (
	"fmt"
	"maps"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type (
	common struct {
		file     *FileDescriptor
		path     string
		LongName string
		FullName string

		OptionExtensions map[string]any
	}

	// An ImportedDescriptor describes a type that was imported by a FileDescriptor.
	ImportedDescriptor struct {
		common
	}

	// A FileDescriptor describes a single proto file with all of its messages, enums, services, etc.
	FileDescriptor struct {
		comments Comments
		*descriptorpb.FileDescriptorProto

		PackageComments *Comment
		SyntaxComments  *Comment
		EditionComments *Comment

		Enums      []*EnumDescriptor
		Extensions []*ExtensionDescriptor
		Imports    []*ImportedDescriptor
		Messages   []*Descriptor
		Services   []*ServiceDescriptor

		OptionExtensions map[string]any

		// optResolver resolves extensions defined in the parsed descriptor set so option
		// extensions can be decoded without the consumer registering them globally.
		optResolver extResolver
	}

	// An EnumDescriptor describe an enum type
	EnumDescriptor struct {
		common
		*descriptorpb.EnumDescriptorProto
		Parent   *Descriptor
		Values   []*EnumValueDescriptor
		Comments *Comment
	}

	// An EnumValueDescriptor describes an enum value
	EnumValueDescriptor struct {
		common
		*descriptorpb.EnumValueDescriptorProto
		Enum     *EnumDescriptor
		Comments *Comment
	}

	// An ExtensionDescriptor describes a protobuf extension. If it's a top-level extension it's parent will be `nil`
	ExtensionDescriptor struct {
		common
		*descriptorpb.FieldDescriptorProto
		Parent   *Descriptor
		Comments *Comment
	}

	// A Descriptor describes a message
	Descriptor struct {
		common
		*descriptorpb.DescriptorProto
		Parent     *Descriptor
		Comments   *Comment
		Enums      []*EnumDescriptor
		Extensions []*ExtensionDescriptor
		Fields     []*FieldDescriptor
		Messages   []*Descriptor
	}

	// A FieldDescriptor describes a message field
	FieldDescriptor struct {
		common
		*descriptorpb.FieldDescriptorProto
		Comments *Comment
		Message  *Descriptor
	}

	// A ServiceDescriptor describes a service
	ServiceDescriptor struct {
		common
		*descriptorpb.ServiceDescriptorProto
		Comments *Comment
		Methods  []*MethodDescriptor
	}

	// A MethodDescriptor describes a method in a service
	MethodDescriptor struct {
		common
		*descriptorpb.MethodDescriptorProto
		Comments *Comment
		Service  *ServiceDescriptor
	}
)

// GetFile returns the FileDescriptor that contains this object
func (c *common) GetFile() *FileDescriptor { return c.file }

// GetPackage returns the package this object is in
func (c *common) GetPackage() string { return c.file.GetPackage() }

// GetLongName returns the name prefixed with the dot-separated parent descriptor's name (if any)
func (c *common) GetLongName() string { return c.LongName }

// GetFullName returns the `LongName` prefixed with the package this object is in
func (c *common) GetFullName() string { return c.FullName }

// IsProto3 returns whether or not this is a proto3 object or uses proto3-like semantics
func (c *common) IsProto3() bool { return c.file.IsProto3() }

// GetEdition returns the edition of the file this object belongs to
func (c *common) GetEdition() descriptorpb.Edition { return c.file.GetEdition() }

// IsEditions returns whether or not this object belongs to a file using editions syntax
func (c *common) IsEditions() bool { return c.file.IsEditions() }

// extResolver resolves extension and message types when decoding option extensions. Both
// *protoregistry.Types (e.g. protoregistry.GlobalTypes) and *dynamicpb.Types satisfy it, and it
// matches the resolver expected by proto.UnmarshalOptions.
type extResolver interface {
	protoregistry.MessageTypeResolver
	protoregistry.ExtensionTypeResolver
}

// combinedResolver tries a primary resolver first, then a fallback. This lets us prefer concrete
// Go types registered in the global registry (so consumers that import generated extension packages
// still get concrete values) while falling back to types derived from the parsed descriptor set (so
// extensions that were never registered globally still decode, as dynamic values).
type combinedResolver struct {
	primary  extResolver
	fallback extResolver
}

func (r *combinedResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	if r.primary != nil {
		if et, err := r.primary.FindExtensionByName(field); err == nil {
			return et, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.FindExtensionByName(field)
	}
	return nil, protoregistry.NotFound
}

func (r *combinedResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	if r.primary != nil {
		if et, err := r.primary.FindExtensionByNumber(message, field); err == nil {
			return et, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.FindExtensionByNumber(message, field)
	}
	return nil, protoregistry.NotFound
}

func (r *combinedResolver) FindMessageByName(name protoreflect.FullName) (protoreflect.MessageType, error) {
	if r.primary != nil {
		if mt, err := r.primary.FindMessageByName(name); err == nil {
			return mt, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.FindMessageByName(name)
	}
	return nil, protoregistry.NotFound
}

func (r *combinedResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	if r.primary != nil {
		if mt, err := r.primary.FindMessageByURL(url); err == nil {
			return mt, nil
		}
	}
	if r.fallback != nil {
		return r.fallback.FindMessageByURL(url)
	}
	return nil, protoregistry.NotFound
}

// getOptions decodes the extension fields set on an options message, keyed by the extension's full
// name. It re-decodes the options bytes using the supplied resolver so that extensions defined in
// the parsed descriptor set are recognized even when they were not registered in the global type
// registry. Values are returned as protoreflect.Value.Interface() results: proto2 optional scalars
// come back by value (e.g. bool, not *bool), and message-typed extensions come back as a proto
// message (a concrete registered type when available, otherwise a *dynamicpb.Message).
func getOptions(res extResolver, options proto.Message) (m map[string]any) {
	// Re-decode against the resolver. The options message we receive was unmarshaled with the
	// global registry, so globally-registered extensions are already present as known fields and
	// only unregistered extensions sit in the unknown fields. Round-tripping the raw bytes through
	// the resolver promotes those to proper extension fields; skip it entirely when there are no
	// unknown fields, since there is then nothing to promote.
	if res != nil && len(options.ProtoReflect().GetUnknown()) > 0 {
		if raw, err := proto.Marshal(options); err == nil {
			decoded := options.ProtoReflect().New().Interface()
			if err := (proto.UnmarshalOptions{Resolver: res}).Unmarshal(raw, decoded); err == nil {
				options = decoded
			}
		}
	}

	options.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.IsExtension() {
			if m == nil {
				m = make(map[string]any)
			}
			m[string(fd.FullName())] = v.Interface()
		}
		return true
	})

	return m
}

func (c *common) setOptions(options proto.Message) {
	if opts := getOptions(c.file.optResolver, options); len(opts) > 0 {
		if c.OptionExtensions == nil {
			c.OptionExtensions = opts
			return
		}

		maps.Copy(c.OptionExtensions, opts)
	}
}

// FileDescriptor methods

// IsProto3 returns whether or not this file is a proto3 file or uses proto3-like semantics
func (f *FileDescriptor) IsProto3() bool {
	// Original proto3 syntax
	if f.GetSyntax() == "proto3" {
		return true
	}
	// Editions with proto3-like behavior (IMPLICIT field presence) match proto3 semantics
	if f.IsEditions() {
		if options := f.GetOptions(); options != nil {
			if features := options.GetFeatures(); features != nil {
				return features.GetFieldPresence() == descriptorpb.FeatureSet_IMPLICIT
			}
		}
	}
	return false
}

// GetEdition returns the edition of this file
func (f *FileDescriptor) GetEdition() descriptorpb.Edition { return f.FileDescriptorProto.GetEdition() }

// IsEditions returns whether or not this file uses the editions syntax
func (f *FileDescriptor) IsEditions() bool { return f.GetSyntax() == "editions" }

// GetEditionName returns the edition name as a string (e.g., "2023", "2024")
func (f *FileDescriptor) GetEditionName() string {
	if !f.IsEditions() {
		return ""
	}
	switch f.GetEdition() {
	case descriptorpb.Edition_EDITION_2023:
		return "2023"
	case descriptorpb.Edition_EDITION_2024:
		return "2024"
	case descriptorpb.Edition_EDITION_PROTO2:
		return "proto2"
	case descriptorpb.Edition_EDITION_PROTO3:
		return "proto3"
	case descriptorpb.Edition_EDITION_UNKNOWN, descriptorpb.Edition_EDITION_LEGACY:
		return "unknown"
	case descriptorpb.Edition_EDITION_1_TEST_ONLY:
		return "1_test_only"
	case descriptorpb.Edition_EDITION_2_TEST_ONLY:
		return "2_test_only"
	case descriptorpb.Edition_EDITION_99997_TEST_ONLY:
		return "99997_test_only"
	case descriptorpb.Edition_EDITION_99998_TEST_ONLY:
		return "99998_test_only"
	case descriptorpb.Edition_EDITION_99999_TEST_ONLY:
		return "99999_test_only"
	case descriptorpb.Edition_EDITION_MAX:
		return "max"
	default:
		return f.GetEdition().String()
	}
}

// GetPackageComments returns the file's package comments
func (f *FileDescriptor) GetPackageComments() *Comment { return f.PackageComments }

// GetSyntaxComments returns the file's syntax comments
func (f *FileDescriptor) GetSyntaxComments() *Comment { return f.SyntaxComments }

// GetEditionComments returns the file's edition comments
func (f *FileDescriptor) GetEditionComments() *Comment { return f.EditionComments }

// HasExplicitFieldPresence returns whether this file defaults to explicit field presence
// In editions 2023+, field presence is explicit by default (like proto2)
// In proto3, field presence is implicit by default
func (f *FileDescriptor) HasExplicitFieldPresence() bool {
	if f.IsEditions() {
		// Check custom field presence setting in editions
		if options := f.GetOptions(); options != nil {
			if features := options.GetFeatures(); features != nil {
				switch features.GetFieldPresence() {
				case descriptorpb.FeatureSet_IMPLICIT:
					return false
				case descriptorpb.FeatureSet_EXPLICIT, descriptorpb.FeatureSet_LEGACY_REQUIRED:
					return true
				case descriptorpb.FeatureSet_FIELD_PRESENCE_UNKNOWN:
					// Fall through to default behavior
				}
			}
		}
		// Editions 2023+ default to explicit field presence
		return true
	}
	// proto2 has explicit field presence, proto3 has implicit
	return f.GetSyntax() == "proto2"
}

// GetSyntaxType returns a more detailed syntax classification
func (f *FileDescriptor) GetSyntaxType() string {
	if f.IsEditions() {
		return "editions"
	}
	return f.GetSyntax()
}

// GetEnums returns the top-level enumerations defined in this file
func (f *FileDescriptor) GetEnums() []*EnumDescriptor { return f.Enums }

// GetExtensions returns the top-level (file) extensions defined in this file
func (f *FileDescriptor) GetExtensions() []*ExtensionDescriptor { return f.Extensions }

// GetImports returns the proto files imported by this file
func (f *FileDescriptor) GetImports() []*ImportedDescriptor { return f.Imports }

// GetMessages returns the top-level messages defined in this file
func (f *FileDescriptor) GetMessages() []*Descriptor { return f.Messages }

// GetServices returns the services defined in this file
func (f *FileDescriptor) GetServices() []*ServiceDescriptor { return f.Services }

// GetEnum returns the enumeration with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetEnum(name string) *EnumDescriptor {
	for _, e := range f.GetEnums() {
		if e.GetName() == name || e.GetLongName() == name {
			return e
		}
	}

	return nil
}

// GetMessage returns the message with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetMessage(name string) *Descriptor {
	for _, m := range f.GetMessages() {
		if m.GetName() == name || m.GetLongName() == name {
			return m
		}
	}

	return nil
}

// GetService returns the service with the specified name (returns `nil` if not found)
func (f *FileDescriptor) GetService(name string) *ServiceDescriptor {
	for _, s := range f.GetServices() {
		if s.GetName() == name || s.GetLongName() == name {
			return s
		}
	}

	return nil
}

func (f *FileDescriptor) setOptions(options proto.Message) {
	if opts := getOptions(f.optResolver, options); len(opts) > 0 {
		if f.OptionExtensions == nil {
			f.OptionExtensions = opts
			return
		}

		maps.Copy(f.OptionExtensions, opts)
	}
}

// EnumDescriptor methods

// GetComments returns a description of this enum
func (e *EnumDescriptor) GetComments() *Comment { return e.Comments }

// GetParent returns the parent message (if any) that contains this enum
func (e *EnumDescriptor) GetParent() *Descriptor { return e.Parent }

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

// EnumValueDescriptor methods

// GetComments returns a description of the value
func (v *EnumValueDescriptor) GetComments() *Comment { return v.Comments }

// GetEnum returns the parent enumeration that contains this value
func (v *EnumValueDescriptor) GetEnum() *EnumDescriptor { return v.Enum }

// ExtensionDescriptor methods

// GetComments returns a description of the extension
func (e *ExtensionDescriptor) GetComments() *Comment { return e.Comments }

// GetParent returns the descriptor that defined this extension (if any)
func (e *ExtensionDescriptor) GetParent() *Descriptor { return e.Parent }

// Descriptor methods

// GetComments returns a description of the message
func (m *Descriptor) GetComments() *Comment { return m.Comments }

// GetParent returns the parent descriptor (if any) that defines this descriptor
func (m *Descriptor) GetParent() *Descriptor { return m.Parent }

// GetEnums returns the nested enumerations within the message
func (m *Descriptor) GetEnums() []*EnumDescriptor { return m.Enums }

// GetExtensions returns the message-level extensions defined by this message
func (m *Descriptor) GetExtensions() []*ExtensionDescriptor { return m.Extensions }

// GetMessages returns the nested messages within the message
func (m *Descriptor) GetMessages() []*Descriptor { return m.Messages }

// GetMessageFields returns the message fields
func (m *Descriptor) GetMessageFields() []*FieldDescriptor { return m.Fields }

// GetEnum returns the enum with the specified name. The name can be either simple, or fully qualified (returns `nil` if
// not found)
func (m *Descriptor) GetEnum(name string) *EnumDescriptor {
	for _, e := range m.GetEnums() {
		// can lookup by name or message prefixed name (qualified)
		if e.GetName() == name || e.GetLongName() == name {
			return e
		}
	}

	return nil
}

// GetMessage returns the nested message with the specified name. The name can be simple or fully qualified (returns
// `nil` if not found)
func (m *Descriptor) GetMessage(name string) *Descriptor {
	for _, msg := range m.GetMessages() {
		// can lookup by name or message prefixed name (qualified)
		if msg.GetName() == name || msg.GetLongName() == name {
			return msg
		}
	}

	return nil
}

// GetMessageField returns the field with the specified name (returns `nil` if not found)
func (m *Descriptor) GetMessageField(name string) *FieldDescriptor {
	for _, f := range m.GetMessageFields() {
		if f.GetName() == name || f.GetLongName() == name {
			return f
		}
	}

	return nil
}

// FieldDescriptor methods

// GetComments returns a description of the field
func (mf *FieldDescriptor) GetComments() *Comment { return mf.Comments }

// GetMessage returns the descriptor that defines this field
func (mf *FieldDescriptor) GetMessage() *Descriptor { return mf.Message }

// ServiceDescriptor methods

// GetComments returns a description of the service
func (s *ServiceDescriptor) GetComments() *Comment { return s.Comments }

// GetMethods returns the methods for the service
func (s *ServiceDescriptor) GetMethods() []*MethodDescriptor { return s.Methods }

// GetNamedMethod returns the method with the specified name (if found)
func (s *ServiceDescriptor) GetNamedMethod(name string) *MethodDescriptor {
	for _, m := range s.GetMethods() {
		if m.GetName() == name || m.GetLongName() == name {
			return m
		}
	}

	return nil
}

// MethodDescriptor methods

// GetComments returns a description of the method
func (m *MethodDescriptor) GetComments() *Comment { return m.Comments }

// GetService returns the service descriptor that defines this method
func (m *MethodDescriptor) GetService() *ServiceDescriptor { return m.Service }

// newCommon creates a new common struct with the given parameters.
func newCommon(f *FileDescriptor, path, longName string) common {
	fn := longName
	if !strings.HasPrefix(fn, ".") {
		fn = fmt.Sprintf("%s.%s", f.GetPackage(), longName)
	}

	return common{
		file:     f,
		path:     path,
		LongName: longName,
		FullName: fn,
	}
}
