// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package plugins

import (
	"context"

	v1 "github.com/openchami/inventory-service/apis/inventory-service.openchami.org/v1"
)

// StorageExtras defines the interface for extra storage operations
// that are not part of the standard CRUD interface.
type StorageExtras interface {
	LoadComponentByID(ctx context.Context, id string) (*v1.Component, error)
	LoadComponentEndpointByID(ctx context.Context, id string) (*v1.ComponentEndpoint, error)
	LoadRedfishEndpointByID(ctx context.Context, id string) (*v1.RedfishEndpoint, error)
	LoadEthernetInterfaceByID(ctx context.Context, id string) (*v1.EthernetInterface, error)
	LoadServiceEndpointByID(ctx context.Context, id string) (*v1.ServiceEndpoint, error)
	LoadGroupByLabel(ctx context.Context, label string) (*v1.Group, error)
	LoadHardwareByID(ctx context.Context, id string) (*v1.Hardware, error)
	LoadServiceEndpointsByRedfishType(ctx context.Context, redfishType string) ([]*v1.ServiceEndpoint, error)
	LoadServiceEndpointsByRedfishTypeAndID(ctx context.Context, redfishType string, redfishID string) ([]*v1.ServiceEndpoint, error)

	// Lock records back the SMD-compatible /hsm/v2/locks endpoints. Locks are
	// keyed by component xname (Spec.ID) and are not exposed as a native REST
	// resource.
	LoadAllLocks(ctx context.Context) ([]*v1.Lock, error)
	LoadLock(ctx context.Context, uid string) (*v1.Lock, error)
	LoadLockByID(ctx context.Context, id string) (*v1.Lock, error)
	SaveLock(ctx context.Context, resource *v1.Lock) error
	DeleteLock(ctx context.Context, uid string) error
	DeleteLockByID(ctx context.Context, id string) error
}
