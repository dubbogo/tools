// protoc-gen-go-triple is a plugin for the Google protocol buffer compiler to
// generate Go code. Install it by building this program and making it
// accessible within your PATH with the name:
//	protoc-gen-go-triple
//
// The 'go-grpc' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//	protoc --go-triple_out=. path/to/file.proto
//
// This generates Go service definitions for the protocol buffer defined by
// file.proto.  With that input, the output will be written to:
//	path/to/file_triple.pb.go
package main

import (
	"flag"
	"fmt"
)

import (
	"google.golang.org/protobuf/compiler/protogen"

	"google.golang.org/protobuf/types/pluginpb"
)

import (
	"github.com/dubbogo/tools/cmd/protoc-gen-go-triple/triple"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-triple %v\n", triple.Version)
		return
	}

	var flags flag.FlagSet
	triple.RequireUnimplemented = flags.Bool("require_unimplemented_servers", true, "set to false to match legacy behavior")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			triple.GenerateFile(gen, f)
		}
		return nil
	})
}
