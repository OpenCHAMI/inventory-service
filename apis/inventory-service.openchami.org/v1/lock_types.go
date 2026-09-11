// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

// Lock is the SMD-compatible lock / reservation record. There is at most one
// Lock per Component, keyed by the component xname (Spec.ID). It is persisted in
// the same generic resource table as every other resource, but is not exposed
// as a native REST resource; it is only reached through the /hsm/v2/locks
// SMD-compatibility endpoints.
//
// This type mirrors the layout that Fabrica generates for the other resources
// (see component_types.go) so that a future `fabrica generate` run can adopt it
// without churn. It is hand-maintained today because the pinned Fabrica release
// (0.4.9) is not available in this workspace.
type Lock struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	ID         string           `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       LockSpec         `json:"spec" yaml:"spec" validate:"required"`
	Status     LockStatus       `json:"status,omitempty" yaml:"status,omitempty"`
}

// LockSpec is the durable reservation state for a single component.
type LockSpec struct {
	ID             string `json:"ID" yaml:"ID"`
	ReservationKey string `json:"ReservationKey,omitempty" yaml:"ReservationKey,omitempty"`
	DeputyKey      string `json:"DeputyKey,omitempty" yaml:"DeputyKey,omitempty"`
	CreatedTime    string `json:"CreatedTime,omitempty" yaml:"CreatedTime,omitempty"`
	// ExpirationTime is an RFC3339 timestamp. An empty value denotes an
	// indefinite (admin) reservation.
	ExpirationTime string `json:"ExpirationTime,omitempty" yaml:"ExpirationTime,omitempty"`
}

type LockStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
}

func (r *Lock) Validate(ctx context.Context) error {
	return nil
}

func (r *Lock) GetKind() string {
	return "Lock"
}

func (r *Lock) GetName() string {
	return r.Metadata.Name
}

func (r *Lock) GetUID() string {
	return r.Metadata.UID
}

func (r *Lock) IsHub() {}
