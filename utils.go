package protokit

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	"errors"
	"io/ioutil"
	"path/filepath"
)

// LoadDescriptor loads file descriptor protos from a file on disk, and returns the named proto descriptor. This is
// useful mostly for testing purposes.
func LoadDescriptor(name string, pathSegments ...string) (*descriptor.FileDescriptorProto, error) {
	f, err := ioutil.ReadFile(filepath.Join(pathSegments...))
	if err != nil {
		return nil, err
	}

	set := new(descriptor.FileDescriptorSet)
	if err = proto.Unmarshal(f, set); err != nil {
		return nil, err
	}

	for _, pf := range set.GetFile() {
		if filepath.Base(pf.GetName()) == name {
			return pf, nil
		}
	}

	return nil, errors.New("FileDescriptor not found")
}
