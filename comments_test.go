package protokit_test

import (
	"github.com/stretchr/testify/suite"

	"testing"

	"github.com/pseudomuto/protokit"
)

type CommentsTest struct {
	suite.Suite
	comments protokit.Comments
}

func TestComments(t *testing.T) {
	suite.Run(t, new(CommentsTest))
}

func (assert *CommentsTest) SetupSuite() {
	pf, err := protokit.LoadDescriptor("todo.proto", "fixtures", "fileset.pb")
	assert.NoError(err)

	assert.comments = protokit.ParseComments(pf)
}

func (assert *CommentsTest) TestComments() {
	tests := []struct {
		key   string
		value string
	}{
		{"6.0.2.1", "Add an item to your list\n\nAdds a new item to the specified list."}, // leading commend
		{"4.0.2.0", "The id of the list."},                                                // tailing comment
	}

	for _, test := range tests {
		assert.Equal(test.value, assert.comments[test.key])
	}
}
