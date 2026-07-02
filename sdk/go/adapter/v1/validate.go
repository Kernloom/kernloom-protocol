// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package adapterv1

import (
	"context"
	"errors"
	"fmt"
)

const ProtocolVersion = "adapter/v1"

const (
	FacetDescribe                   = "Describe"
	FacetHealth                     = "Health"
	FacetPlanConfig                 = "PlanConfig"
	FacetValidateConfig             = "ValidateConfig"
	FacetReadObservedState          = "ReadObservedState"
	FacetNormalizeState             = "NormalizeState"
	FacetProvideConformanceEvidence = "ProvideConformanceEvidence"
	FacetReadSignals                = "ReadSignals"
	FacetStreamSignals              = "StreamSignals"
	FacetProvideRelationships       = "ProvideRelationships"
	FacetExecuteRuntimeAction       = "ExecuteRuntimeAction"
	FacetGetRuntimeActionState      = "GetRuntimeActionState"
	FacetRevokeRuntimeAction        = "RevokeRuntimeAction"
	FacetProvideAttestationEvidence = "ProvideAttestationEvidence"
)

const (
	HealthServing    = "serving"
	HealthDegraded   = "degraded"
	HealthNotServing = "not_serving"
	HealthUnknown    = "unknown"
)

type MinimalAdapter interface {
	Describe(context.Context) (*AdapterDescriptor, error)
	Health(context.Context) (*HealthResponse, error)
}

func ValidateDescriptor(desc *AdapterDescriptor) error {
	if desc == nil {
		return errors.New("adapter descriptor is nil")
	}
	if desc.AdapterId == "" {
		return errors.New("adapter descriptor requires adapter id")
	}
	if desc.Name == "" {
		return errors.New("adapter descriptor requires name")
	}
	if desc.ProtocolVersion != ProtocolVersion {
		return fmt.Errorf("adapter protocol version %q does not match %q", desc.ProtocolVersion, ProtocolVersion)
	}
	if !hasFacet(desc.Facets, FacetDescribe) {
		return fmt.Errorf("adapter %q must declare Describe facet", desc.AdapterId)
	}
	if !hasFacet(desc.Facets, FacetHealth) {
		return fmt.Errorf("adapter %q must declare Health facet", desc.AdapterId)
	}
	for _, capability := range desc.Capabilities {
		if capability.GetId() == "" {
			return fmt.Errorf("adapter %q contains capability without id", desc.AdapterId)
		}
	}
	for _, privilege := range desc.Privileges {
		if privilege.GetId() == "" {
			return fmt.Errorf("adapter %q contains privilege without id", desc.AdapterId)
		}
	}
	return nil
}

func hasFacet(facets []string, required string) bool {
	for _, facet := range facets {
		if facet == required {
			return true
		}
	}
	return false
}
