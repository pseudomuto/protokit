package protokit

import (
	"context"
)

type contextKey string

const (
	fileContextKey           = contextKey("file")
	commentsContextKey       = contextKey("comments")
	locationPrefixContextKey = contextKey("locationPrefix")
	serviceContextKey        = contextKey("service")
	messageContextKey        = contextKey("message")
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

// ContextWithComments returns a new context with `comments`
func ContextWithComments(ctx context.Context, comments Comments) context.Context {
	return context.WithValue(ctx, commentsContextKey, comments)
}

// CommentsFromContext returns the `Comments` from the context and whether or not the key was found.
func CommentsFromContext(ctx context.Context) (Comments, bool) {
	val, ok := ctx.Value(commentsContextKey).(Comments)
	return val, ok
}

// ContextWithLocationPrefix returns a new context with `locationPrefix`
func ContextWithLocationPrefix(ctx context.Context, locationPrefix string) context.Context {
	return context.WithValue(ctx, locationPrefixContextKey, locationPrefix)
}

// LocationPrefixFromContext returns the `LocationPrefix` from the context and whether or not the key was found.
func LocationPrefixFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(locationPrefixContextKey).(string)
	return val, ok
}

// ContextWithService returns a new context with `service`
func ContextWithService(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceContextKey, service)
}

// ServiceFromContext returns the `Service` from the context and whether or not the key was found.
func ServiceFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(serviceContextKey).(string)
	return val, ok
}

// ContextWithMessage returns a new context with `message`
func ContextWithMessage(ctx context.Context, message string) context.Context {
	return context.WithValue(ctx, messageContextKey, message)
}

// MessageFromContext returns the `Message` from the context and whether or not the key was found.
func MessageFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(messageContextKey).(string)
	return val, ok
}
