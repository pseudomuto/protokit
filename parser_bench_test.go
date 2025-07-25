package protokit_test

import (
	"testing"

	"github.com/Djarvur/protokit"
	"github.com/Djarvur/protokit/utils"
)

func BenchmarkParseCodeGenRequest(b *testing.B) {
	fds, _ := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	req := utils.CreateGenRequest(fds, "booking.proto", "todo.proto")

	for i := 0; i < b.N; i++ {
		protokit.ParseCodeGenRequest(req)
	}
}
