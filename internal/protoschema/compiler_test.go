// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package protoschema

import (
	"context"
	"testing"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestAdapterV1ProtoCompiles(t *testing.T) {
	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: []string{"../../proto"},
		},
	}

	files, err := compiler.Compile(context.Background(), "kernloom/adapter/v1/adapter.proto")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one compiled proto file, got %d", len(files))
	}

	service := files[0].Services().ByName(protoreflect.Name("AdapterService"))
	if service == nil {
		t.Fatal("compiled proto missing AdapterService")
	}
	for _, name := range []string{"Describe", "Health"} {
		if service.Methods().ByName(protoreflect.Name(name)) == nil {
			t.Fatalf("compiled proto missing %s method", name)
		}
	}
}
