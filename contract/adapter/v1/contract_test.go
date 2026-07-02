// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package contractv1

import (
	"context"
	"testing"

	adapterv1 "github.com/kernloom/kernloom-protocol/sdk/go/adapter/v1"
)

type emptyAdapter struct{}

func (emptyAdapter) Describe(context.Context) (*adapterv1.AdapterDescriptor, error) {
	return &adapterv1.AdapterDescriptor{
		AdapterId:       "empty.test",
		Name:            "Empty Contract Adapter",
		ProtocolVersion: adapterv1.ProtocolVersion,
		Facets: []string{
			adapterv1.FacetDescribe,
			adapterv1.FacetHealth,
		},
	}, nil
}

func (emptyAdapter) Health(context.Context) (*adapterv1.HealthResponse, error) {
	return &adapterv1.HealthResponse{Status: adapterv1.HealthServing}, nil
}

func TestEmptyAdapterPassesMinimalContract(t *testing.T) {
	RunMinimalContract(t, emptyAdapter{})
}

type emptyServiceAdapter struct {
	adapterv1.UnimplementedAdapterServiceServer
}

func (emptyServiceAdapter) Describe(context.Context, *adapterv1.DescribeRequest) (*adapterv1.DescribeResponse, error) {
	return &adapterv1.DescribeResponse{Adapter: &adapterv1.AdapterDescriptor{
		AdapterId:       "empty.service.test",
		Name:            "Empty Service Contract Adapter",
		ProtocolVersion: adapterv1.ProtocolVersion,
		Facets: []string{
			adapterv1.FacetDescribe,
			adapterv1.FacetHealth,
		},
		FacetDescriptors: []*adapterv1.FacetDescriptor{
			{Name: adapterv1.FacetDescribe, Status: adapterv1.FacetStatusImplemented},
			{Name: adapterv1.FacetHealth, Status: adapterv1.FacetStatusImplemented},
		},
	}}, nil
}

func (emptyServiceAdapter) Health(context.Context, *adapterv1.HealthRequest) (*adapterv1.HealthResponse, error) {
	return &adapterv1.HealthResponse{Status: adapterv1.HealthServing}, nil
}

func TestEmptyAdapterPassesServiceContract(t *testing.T) {
	RunServiceContract(t, emptyServiceAdapter{})
}
