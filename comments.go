package protokit

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"strconv"
	"strings"
)

// Comments is a map of source location paths to values.
//
// The values are a concatenation of the leader and trailing comments (null delimited). If either are missing, they are
// simply removed so that the null terminator isn't included.
type Comments map[string]string

// ParseComments parses all comments within a proto file. The locations are encoded into the map by joining the paths
// with a "." character. E.g. `4.2.3.0`.
//
// Leading spaces are trimmed for each distinct value (leading, trailing) before joinging with `\x00`.
func ParseComments(fd *descriptor.FileDescriptorProto) Comments {
	comments := make(Comments)

	for _, loc := range fd.GetSourceCodeInfo().GetLocation() {
		leading := loc.GetLeadingComments()
		trailing := loc.GetTrailingComments()

		if leading == "" && trailing == "" {
			continue
		}

		path := loc.GetPath()
		key := make([]string, len(path))
		for idx, p := range path {
			key[idx] = strconv.Itoa(int(p))
		}

		parts := make([]string, 0, 2)
		if leading != "" {
			parts = append(parts, scrub(leading))
		}

		if trailing != "" {
			parts = append(parts, scrub(trailing))
		}

		comments[strings.Join(key, ".")] = strings.Join(parts, "\x00")
	}

	return comments
}

func scrub(str string) string {
	return strings.TrimSpace(strings.Replace(str, "\n ", "\n", -1))
}
