package protokit_test

import (
	"github.com/stretchr/testify/suite"

	"context"
	"testing"

	"github.com/pseudomuto/protokit"
)

type ContextTest struct {
	suite.Suite
}

func TestContext(t *testing.T) {
	suite.Run(t, new(ContextTest))
}

func (assert *ContextTest) TestContextWithComments() {
	ctx := context.Background()

	val, found := protokit.CommentsFromContext(ctx)
	assert.Nil(val)
	assert.False(found)

	ctx = protokit.ContextWithComments(ctx, make(protokit.Comments))
	val, found = protokit.CommentsFromContext(ctx)
	assert.NotNil(val)
	assert.True(found)
}

func (assert *ContextTest) TestContextWithLocationPrefix() {
	ctx := context.Background()

	val, found := protokit.LocationPrefixFromContext(ctx)
	assert.Empty(val)
	assert.False(found)

	ctx = protokit.ContextWithLocationPrefix(ctx, "prefix")
	val, found = protokit.LocationPrefixFromContext(ctx)
	assert.Equal("prefix", val)
	assert.True(found)
}

func (assert *ContextTest) TestContextWithPackage() {
	ctx := context.Background()

	val, found := protokit.PackageFromContext(ctx)
	assert.Empty(val)
	assert.False(found)

	ctx = protokit.ContextWithPackage(ctx, "package")
	val, found = protokit.PackageFromContext(ctx)
	assert.Equal("package", val)
	assert.True(found)
}

func (assert *ContextTest) TestContextWithService() {
	ctx := context.Background()

	val, found := protokit.ServiceFromContext(ctx)
	assert.Empty(val)
	assert.False(found)

	ctx = protokit.ContextWithService(ctx, "MyService")
	val, found = protokit.ServiceFromContext(ctx)
	assert.Equal("MyService", val)
	assert.True(found)
}

func (assert *ContextTest) TestContextWithMessage() {
	ctx := context.Background()

	val, found := protokit.MessageFromContext(ctx)
	assert.Empty(val)
	assert.False(found)

	ctx = protokit.ContextWithMessage(ctx, "MyMessage")
	val, found = protokit.MessageFromContext(ctx)
	assert.Equal("MyMessage", val)
	assert.True(found)
}
