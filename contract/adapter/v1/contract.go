// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package contractv1

import (
	"context"
	"testing"
	"time"

	adapterv1 "github.com/kernloom/kernloom-protocol/sdk/go/adapter/v1"
	"google.golang.org/grpc/metadata"
)

func RunMinimalContract(t *testing.T, impl adapterv1.MinimalAdapter) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	desc, err := impl.Describe(ctx)
	if err != nil {
		t.Fatalf("Describe returned error: %v", err)
	}
	if err := adapterv1.ValidateDescriptor(desc); err != nil {
		t.Fatalf("Describe returned invalid descriptor: %v", err)
	}

	health, err := impl.Health(ctx)
	if err != nil {
		t.Fatalf("Health returned error: %v", err)
	}
	switch health.Status {
	case adapterv1.HealthServing, adapterv1.HealthDegraded:
	default:
		t.Fatalf("Health returned non-serving status %q", health.Status)
	}
}

func RunServiceContract(t *testing.T, impl adapterv1.AdapterServiceServer) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	describe, err := impl.Describe(ctx, &adapterv1.DescribeRequest{})
	if err != nil {
		t.Fatalf("Describe returned error: %v", err)
	}
	desc := describe.GetAdapter()
	if err := adapterv1.ValidateDescriptor(desc); err != nil {
		t.Fatalf("Describe returned invalid descriptor: %v", err)
	}
	for _, facet := range adapterv1.ImplementedFacets(desc) {
		callImplementedFacet(t, ctx, impl, facet)
	}
}

func callImplementedFacet(t *testing.T, ctx context.Context, impl adapterv1.AdapterServiceServer, facet string) {
	t.Helper()
	switch facet {
	case adapterv1.FacetDescribe:
		if _, err := impl.Describe(ctx, &adapterv1.DescribeRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetHealth:
		health, err := impl.Health(ctx, &adapterv1.HealthRequest{})
		if err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
		switch health.GetStatus() {
		case adapterv1.HealthServing, adapterv1.HealthDegraded:
		default:
			t.Fatalf("%s returned non-serving status %q", facet, health.GetStatus())
		}
	case adapterv1.FacetPlanConfig:
		if _, err := impl.PlanConfig(ctx, &adapterv1.PlanConfigRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetValidateConfig:
		if _, err := impl.ValidateConfig(ctx, &adapterv1.ValidateConfigRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetReadObservedState:
		if _, err := impl.ReadObservedState(ctx, &adapterv1.ReadObservedStateRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetNormalizeState:
		if _, err := impl.NormalizeState(ctx, &adapterv1.NormalizeStateRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetProvideConformanceEvidence:
		if _, err := impl.ProvideConformanceEvidence(ctx, &adapterv1.ProvideConformanceEvidenceRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetReadSignals:
		if _, err := impl.ReadSignals(ctx, &adapterv1.ReadSignalsRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetStreamSignals:
		if err := impl.StreamSignals(&adapterv1.StreamSignalsRequest{}, &streamRecorder{ctx: ctx}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetProvideRelationships:
		if _, err := impl.ProvideRelationships(ctx, &adapterv1.ProvideRelationshipsRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetExecuteRuntimeAction:
		if _, err := impl.ExecuteRuntimeAction(ctx, &adapterv1.ExecuteRuntimeActionRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetGetRuntimeActionState:
		if _, err := impl.GetRuntimeActionState(ctx, &adapterv1.GetRuntimeActionStateRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetRevokeRuntimeAction:
		if _, err := impl.RevokeRuntimeAction(ctx, &adapterv1.RevokeRuntimeActionRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	case adapterv1.FacetProvideAttestationEvidence:
		if _, err := impl.ProvideAttestationEvidence(ctx, &adapterv1.ProvideAttestationEvidenceRequest{}); err != nil {
			t.Fatalf("%s returned error: %v", facet, err)
		}
	default:
		t.Fatalf("unknown implemented facet %q", facet)
	}
}

type streamRecorder struct {
	ctx  context.Context
	sent []*adapterv1.StreamSignalsResponse
}

func (s *streamRecorder) Send(resp *adapterv1.StreamSignalsResponse) error {
	s.sent = append(s.sent, resp)
	return nil
}

func (s *streamRecorder) SetHeader(metadata.MD) error {
	return nil
}

func (s *streamRecorder) SendHeader(metadata.MD) error {
	return nil
}

func (s *streamRecorder) SetTrailer(metadata.MD) {}

func (s *streamRecorder) Context() context.Context {
	return s.ctx
}

func (s *streamRecorder) SendMsg(any) error {
	return nil
}

func (s *streamRecorder) RecvMsg(any) error {
	return nil
}
