package protokit_test

import (
	"context"
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/stretchr/testify/require"
)

func TestContextWithFileDescriptor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	val, found := protokit.FileDescriptorFromContext(ctx)
	require.Nil(t, val)
	require.False(t, found)

	ctx = protokit.ContextWithFileDescriptor(ctx, new(protokit.FileDescriptor))
	val, found = protokit.FileDescriptorFromContext(ctx)
	require.NotNil(t, val)
	require.True(t, found)
}

func TestContextWithEnumDescriptor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	val, found := protokit.EnumDescriptorFromContext(ctx)
	require.Nil(t, val)
	require.False(t, found)

	ctx = protokit.ContextWithEnumDescriptor(ctx, new(protokit.EnumDescriptor))
	val, found = protokit.EnumDescriptorFromContext(ctx)
	require.NotNil(t, val)
	require.True(t, found)
}

func TestContextWithDescriptor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	val, found := protokit.DescriptorFromContext(ctx)
	require.Nil(t, val)
	require.False(t, found)

	ctx = protokit.ContextWithDescriptor(ctx, new(protokit.Descriptor))
	val, found = protokit.DescriptorFromContext(ctx)
	require.NotNil(t, val)
	require.True(t, found)
}

func TestContextWithServiceDescriptor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	val, found := protokit.ServiceDescriptorFromContext(ctx)
	require.Empty(t, val)
	require.False(t, found)

	ctx = protokit.ContextWithServiceDescriptor(ctx, new(protokit.ServiceDescriptor))
	val, found = protokit.ServiceDescriptorFromContext(ctx)
	require.NotNil(t, val)
	require.True(t, found)
}
