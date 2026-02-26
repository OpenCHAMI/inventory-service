// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	v1 "github.com/OpenCHAMI/smd2/apis/smd2.openchami.org/v1"
	"github.com/OpenCHAMI/smd2/internal/storage/ent"
	"github.com/OpenCHAMI/smd2/internal/storage/ent/predicate"
	entresource "github.com/OpenCHAMI/smd2/internal/storage/ent/resource"
)

// LoadComponentByID loads a single Component resource by its Spec.ID from Ent storage
func LoadComponentByID(ctx context.Context, id string) (*v1.Component, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	// Query by spec.ID and kind using a JSON predicate
	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("Component"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), id, sqljson.Path("ID")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load Component with ID %s: %w", id, err)
	}

	// Convert to Fabrica resource
	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.Component), nil
}

func LoadComponentEndpointByID(ctx context.Context, id string) (*v1.ComponentEndpoint, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("ComponentEndpoint"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), id, sqljson.Path("ID")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load ComponentEndpoint with ID %s: %w", id, err)
	}

	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.ComponentEndpoint), nil
}

func LoadRedfishEndpointByID(ctx context.Context, id string) (*v1.RedfishEndpoint, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("RedfishEndpoint"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), id, sqljson.Path("ID")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load RedfishEndpoint with ID %s: %w", id, err)
	}

	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.RedfishEndpoint), nil
}

func LoadEthernetInterfaceByID(ctx context.Context, id string) (*v1.EthernetInterface, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("EthernetInterface"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), id, sqljson.Path("ID")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load EthernetInterface with ID %s: %w", id, err)
	}

	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.EthernetInterface), nil
}

func LoadServiceEndpointByID(ctx context.Context, id string) (*v1.ServiceEndpoint, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("ServiceEndpoint"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), id, sqljson.Path("RedfishEndpointID")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load ServiceEndpoint with ID %s: %w", id, err)
	}

	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.ServiceEndpoint), nil
}

func LoadGroupByLabel(ctx context.Context, label string) (*v1.Group, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ("Group"),
			predicate.Resource(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(s.C(entresource.FieldSpec), label, sqljson.Path("label")))
			}),
		).
		WithLabels().
		WithAnnotations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load Group with label %s: %w", label, err)
	}

	fabricaResource, err := FromEntResource(ctx, entResource)
	if err != nil {
		return nil, err
	}

	return fabricaResource.(*v1.Group), nil
}
