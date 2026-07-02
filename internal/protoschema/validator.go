// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package protoschema

import (
	"fmt"
	"os"
	"strings"
)

func ValidateAdapterV1(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	source := string(data)

	required := []string{
		`syntax = "proto3";`,
		`package kernloom.adapter.v1;`,
		`option go_package = "github.com/kernloom/kernloom-protocol/sdk/go/adapter/v1;adapterv1";`,
		"service AdapterService",
		"rpc Describe(DescribeRequest) returns (DescribeResponse);",
		"rpc Health(HealthRequest) returns (HealthResponse);",
		"message AdapterDescriptor",
		"message FacetDescriptor",
		"message CapabilityDescriptor",
		"message ContextRequirementDescriptor",
		"message PrivilegeDescriptor",
	}
	for _, want := range required {
		if !strings.Contains(source, want) {
			return fmt.Errorf("adapter proto missing required declaration %q", want)
		}
	}
	return nil
}
