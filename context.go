package protokit

import (
	"context"
)

type contextKey string

const (
	fileContextKey       = contextKey("file")
	descriptorContextKey = contextKey("descriptor")
	enumContextKey       = contextKey("enum")
	serviceContextKey    = contextKey("service")
)

// ContextWithFileDescriptor returns a new context with the attached `FileDescriptor`
func ContextWithFileDescriptor(ctx context.Context, fd *FileDescriptor) context.Context {
	return context.WithValue(ctx, fileContextKey, fd)
}

// FileDescriptorFromContext returns the `FileDescriptor` from the context and whether or not the key was found.
func FileDescriptorFromContext(ctx context.Context) (*FileDescriptor, bool) {
	val, ok := ctx.Value(fileContextKey).(*FileDescriptor)
	return val, ok
}

func ContextWithDescriptor(ctx context.Context, d *Descriptor) context.Context {
	return context.WithValue(ctx, descriptorContextKey, d)
}

func DescriptorFromContext(ctx context.Context) (*Descriptor, bool) {
	val, ok := ctx.Value(descriptorContextKey).(*Descriptor)
	return val, ok
}

func ContextWithEnumDescriptor(ctx context.Context, d *EnumDescriptor) context.Context {
	return context.WithValue(ctx, enumContextKey, d)
}

func EnumDescriptorFromContext(ctx context.Context) (*EnumDescriptor, bool) {
	val, ok := ctx.Value(enumContextKey).(*EnumDescriptor)
	return val, ok
}

// ContextWithServiceDescriptor returns a new context with `service`
func ContextWithServiceDescriptor(ctx context.Context, service *ServiceDescriptor) context.Context {
	return context.WithValue(ctx, serviceContextKey, service)
}

// ServiceDescriptorFromContext returns the `Service` from the context and whether or not the key was found.
func ServiceDescriptorFromContext(ctx context.Context) (*ServiceDescriptor, bool) {
	val, ok := ctx.Value(serviceContextKey).(*ServiceDescriptor)
	return val, ok
}
