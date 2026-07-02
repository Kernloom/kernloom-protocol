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

const (
	FacetStatusImplemented = "implemented"
	FacetStatusPlanned     = "planned"
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
	facetStatus, err := validateFacets(desc)
	if err != nil {
		return err
	}
	if facetStatus[FacetDescribe] != FacetStatusImplemented {
		return fmt.Errorf("adapter %q must declare Describe facet as implemented", desc.AdapterId)
	}
	if facetStatus[FacetHealth] != FacetStatusImplemented {
		return fmt.Errorf("adapter %q must declare Health facet as implemented", desc.AdapterId)
	}
	for _, capability := range desc.Capabilities {
		if capability.GetId() == "" {
			return fmt.Errorf("adapter %q contains capability without id", desc.AdapterId)
		}
		if capability.GetKind() == "" {
			return fmt.Errorf("adapter %q capability %q requires kind", desc.AdapterId, capability.GetId())
		}
		if len(capability.GetActions()) == 0 && len(capability.GetRuntimeActions()) == 0 {
			return fmt.Errorf("adapter %q capability %q requires actions or runtime actions", desc.AdapterId, capability.GetId())
		}
	}
	for _, privilege := range desc.Privileges {
		if privilege.GetId() == "" {
			return fmt.Errorf("adapter %q contains privilege without id", desc.AdapterId)
		}
	}
	if len(desc.Capabilities) > 0 && len(desc.ContextRequirements) == 0 {
		return fmt.Errorf("adapter %q declares capabilities without context requirements", desc.AdapterId)
	}
	if declaresRuntimeAction(desc.Capabilities) && len(desc.Privileges) == 0 {
		return fmt.Errorf("adapter %q declares runtime actions without privilege descriptors", desc.AdapterId)
	}
	return nil
}

func ImplementedFacets(desc *AdapterDescriptor) []string {
	status, err := validateFacets(desc)
	if err != nil {
		return nil
	}
	var facets []string
	for name, value := range status {
		if value == FacetStatusImplemented {
			facets = append(facets, name)
		}
	}
	return facets
}

func validateFacets(desc *AdapterDescriptor) (map[string]string, error) {
	status := map[string]string{}
	facetSet := map[string]struct{}{}
	for _, facet := range desc.Facets {
		if facet == "" {
			return nil, fmt.Errorf("adapter %q contains empty facet name", desc.AdapterId)
		}
		if _, exists := facetSet[facet]; exists {
			return nil, fmt.Errorf("adapter %q declares duplicate facet %q", desc.AdapterId, facet)
		}
		facetSet[facet] = struct{}{}
	}
	if len(desc.FacetDescriptors) == 0 {
		for _, facet := range desc.Facets {
			status[facet] = FacetStatusImplemented
		}
		return status, nil
	}
	for _, facet := range desc.FacetDescriptors {
		if facet.GetName() == "" {
			return nil, fmt.Errorf("adapter %q contains facet without name", desc.AdapterId)
		}
		switch facet.GetStatus() {
		case FacetStatusImplemented, FacetStatusPlanned:
		default:
			return nil, fmt.Errorf("adapter %q facet %q has unsupported status %q", desc.AdapterId, facet.GetName(), facet.GetStatus())
		}
		if _, exists := status[facet.GetName()]; exists {
			return nil, fmt.Errorf("adapter %q declares duplicate facet %q", desc.AdapterId, facet.GetName())
		}
		status[facet.GetName()] = facet.GetStatus()
		if _, exists := facetSet[facet.GetName()]; !exists {
			return nil, fmt.Errorf("adapter %q facet descriptor %q is missing from facets list", desc.AdapterId, facet.GetName())
		}
	}
	for facet := range facetSet {
		if _, exists := status[facet]; !exists {
			return nil, fmt.Errorf("adapter %q facet %q is missing a facet descriptor", desc.AdapterId, facet)
		}
	}
	return status, nil
}

func declaresRuntimeAction(capabilities []*CapabilityDescriptor) bool {
	for _, capability := range capabilities {
		if len(capability.GetRuntimeActions()) > 0 {
			return true
		}
	}
	return false
}
