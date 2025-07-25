# protokit

[![CI][github-svg]][github-ci]
[![codecov][codecov-svg]][codecov-url]
[![GoDoc][godoc-svg]][godoc-url]
[![Go Report Card][goreport-svg]][goreport-url]

Note: Due to a prolonged period of inactivity in the mainstream, this project is no longer a fork: we have appropriated it. The module name has been changed to `github.com/Djarvur/protokit`, so all imports must be changed. Original project: https://github.com/pseudomuto/protokit.

A starter kit for building protoc-plugins. Rather than write your own, you can just use an existing one.

See the [examples](examples/) directory for uh...examples.

## Getting Started

```golang
package main

import (
    "google.golang.org/protobuf/proto"
    plugin_go "google.golang.org/protobuf/types/pluginpb"
    "github.com/Djarvur/protokit"
    _ "google.golang.org/genproto/googleapis/api/annotations" // Support (google.api.http) option (from google/api/annotations.proto).

    "log"
)

func main() {
    // all the heavy lifting done for you!
    if err := protokit.RunPlugin(new(plugin)); err != nil {
        log.Fatal(err)
    }
}

// plugin is an implementation of protokit.Plugin
type plugin struct{}

func (p *plugin) Generate(in *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
    descriptors := protokit.ParseCodeGenRequest(req)

    resp := new(plugin_go.CodeGeneratorResponse)

    for _, d := range descriptors {
        // TODO: YOUR WORK HERE
        fileName := // generate a file name based on d.GetName()
        content := // generate content for the output file

        resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
            Name:    proto.String(fileName),
            Content: proto.String(content),
        })
    }

    return resp, nil
}
```

Then invoke your plugin via `protoc`. For example (assuming your app is called `thingy`):

`protoc --plugin=protoc-gen-thingy=./thingy -I. --thingy_out=. rpc/*.proto`

[github-svg]: https://github.com/Djarvur/protokit/actions/workflows/ci.yaml/badge.svg?branch=master
[github-ci]: https://github.com/Djarvur/protokit/actions/workflows/ci.yaml
[codecov-svg]: https://codecov.io/gh/pseudomuto/protokit/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/pseudomuto/protokit
[godoc-svg]: https://godoc.org/github.com/Djarvur/protokit?status.svg
[godoc-url]: https://godoc.org/github.com/Djarvur/protokit
[goreport-svg]: https://goreportcard.com/badge/github.com/Djarvur/protokit
[goreport-url]: https://goreportcard.com/report/github.com/Djarvur/protokit
