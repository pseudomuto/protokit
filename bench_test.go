package protokit_test

import (
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
)

func BenchmarkParseCodeRequest(b *testing.B) {
	fds, _ := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")
	files := utils.FilesToGenerate(req)

	for i := 0; i < b.N; i++ {
		for _, pf := range files {
			protokit.ParseFile(pf)
		}
	}
}
