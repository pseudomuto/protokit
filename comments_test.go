package protokit_test

import (
	"fmt"
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	"github.com/stretchr/testify/require"
)

func TestComments(t *testing.T) {
	t.Parallel()

	pf, err := utils.LoadDescriptor("todo.proto", "fixtures", "fileset.pb")
	require.NoError(t, err)

	comments := protokit.ParseComments(pf)

	tests := []struct {
		key      string
		leading  string
		trailing string
	}{
		{"6.0.2.1", "Add an item to your list\n\nAdds a new item to the specified list.", ""}, // leading commend
		{"4.0.2.0", "", "The id of the list."},                                                // tailing comment
	}

	for _, test := range tests {
		require.Equal(t, test.leading, comments[test.key].GetLeading())
		require.Equal(t, test.trailing, comments[test.key].GetTrailing())
		require.Empty(t, comments[test.key].GetDetached())
	}

	require.NotNil(t, comments.Get("WONTBETHERE"))
	require.Empty(t, comments.Get("WONTBETHERE").String())
}

// Join the leading and trailing comments together
func ExampleComment_String() {
	c := &protokit.Comment{Leading: "Some leading comment", Trailing: "Some trailing comment"}
	fmt.Println(c.String())
	// Output: Some leading comment
	//
	// Some trailing comment
}
