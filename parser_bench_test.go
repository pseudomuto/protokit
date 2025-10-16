package protokit_test

import (
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
)

func BenchmarkParseCodeGenRequest(b *testing.B) {
	fds, _ := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")

	for b.Loop() {
		protokit.ParseCodeGenRequest(req)
	}
}
