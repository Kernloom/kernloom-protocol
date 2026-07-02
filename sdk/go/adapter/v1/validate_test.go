// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package adapterv1

import (
	"strings"
	"testing"
)

func TestValidateDescriptorAcceptsCapabilityWithContextAndPrivilege(t *testing.T) {
	desc := baseDescriptor()
	desc.Capabilities = []*CapabilityDescriptor{{
		Id:             "runtime.block_source",
		Kind:           "runtime_action",
		RuntimeActions: []string{"deny_temporarily"},
	}}
	desc.ContextRequirements = []*ContextRequirementDescriptor{{
		Fact:       "source identity",
		Freshness:  "30s",
		Confidence: "high",
	}}
	desc.Privileges = []*PrivilegeDescriptor{{
		Id:     "klshield.source_block.write",
		Reason: "execute bounded runtime mitigation",
		Scope:  "application",
		Access: "write",
	}}

	if err := ValidateDescriptor(desc); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDescriptorRejectsUnderspecifiedCapability(t *testing.T) {
	desc := baseDescriptor()
	desc.Capabilities = []*CapabilityDescriptor{{Id: "config.observe"}}
	desc.ContextRequirements = []*ContextRequirementDescriptor{{Fact: "observed state"}}

	err := ValidateDescriptor(desc)
	if err == nil {
		t.Fatal("expected capability without kind to be rejected")
	}
	if !strings.Contains(err.Error(), "requires kind") {
		t.Fatalf("expected kind error, got %v", err)
	}
}

func TestValidateDescriptorRejectsRuntimeActionWithoutPrivilege(t *testing.T) {
	desc := baseDescriptor()
	desc.Capabilities = []*CapabilityDescriptor{{
		Id:             "runtime.block_source",
		Kind:           "runtime_action",
		RuntimeActions: []string{"deny_temporarily"},
	}}
	desc.ContextRequirements = []*ContextRequirementDescriptor{{Fact: "source identity"}}

	err := ValidateDescriptor(desc)
	if err == nil {
		t.Fatal("expected runtime action without privileges to be rejected")
	}
	if !strings.Contains(err.Error(), "without privilege descriptors") {
		t.Fatalf("expected privilege error, got %v", err)
	}
}

func baseDescriptor() *AdapterDescriptor {
	return &AdapterDescriptor{
		AdapterId:       "test.adapter",
		Name:            "Test Adapter",
		ProtocolVersion: ProtocolVersion,
		Facets:          []string{FacetDescribe, FacetHealth},
	}
}
