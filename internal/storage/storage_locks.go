// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openchami/fabrica/pkg/fabrica"
	v1 "github.com/openchami/inventory-service/apis/inventory-service.openchami.org/v1"
	"github.com/openchami/inventory-service/internal/storage/ent"
	entresource "github.com/openchami/inventory-service/internal/storage/ent/resource"
)

// lockKind is the resource kind used to persist Lock records in the generic
// resource table. There is at most one Lock per component, enforced by the
// (resource_type, resource_id) unique index.
const lockKind = "Lock"

// lockFromEnt converts a stored Ent resource row into a *v1.Lock. Locks are
// hand-marshaled (rather than routed through the generated FromEntResource
// type switch) so the Lock resource stays entirely within hand-written files.
func lockFromEnt(entResource *ent.Resource) (*v1.Lock, error) {
	lock := &v1.Lock{
		APIVersion: entResource.APIVersion,
		Kind:       entResource.Kind,
		ID:         entResource.ResourceID,
		Metadata: fabrica.Metadata{
			Name:        entResource.Name,
			UID:         entResource.UID,
			CreatedAt:   entResource.CreatedAt,
			UpdatedAt:   entResource.UpdatedAt,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}
	if len(entResource.Spec) > 0 {
		if err := json.Unmarshal(entResource.Spec, &lock.Spec); err != nil {
			return nil, fmt.Errorf("failed to unmarshal spec for Lock: %w", err)
		}
	}
	return lock, nil
}

// LoadAllLocks returns every Lock record.
func LoadAllLocks(ctx context.Context) ([]*v1.Lock, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResources, err := entClient.Resource.Query().
		Where(entresource.KindEQ(lockKind)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load Lock resources: %w", err)
	}

	locks := make([]*v1.Lock, 0, len(entResources))
	for _, entResource := range entResources {
		lock, err := lockFromEnt(entResource)
		if err != nil {
			continue
		}
		locks = append(locks, lock)
	}
	return locks, nil
}

// LoadLock loads a single Lock by its UID.
func LoadLock(ctx context.Context, uid string) (*v1.Lock, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.UIDEQ(uid),
			entresource.KindEQ(lockKind),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load Lock %s: %w", uid, err)
	}
	return lockFromEnt(entResource)
}

// LoadLockByID loads a single Lock by its component xname (Spec.ID).
func LoadLockByID(ctx context.Context, id string) (*v1.Lock, error) {
	if entClient == nil {
		return nil, fmt.Errorf("ent client not initialized")
	}

	entResource, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ(lockKind),
			entresource.ResourceIDEQ(id),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to load Lock with ID %s: %w", id, err)
	}
	return lockFromEnt(entResource)
}

// SaveLock upserts a Lock keyed by its component xname (Spec.ID). Because there
// is exactly one Lock per component, existing rows are updated in place so the
// (resource_type, resource_id) unique index is never violated.
func SaveLock(ctx context.Context, lock *v1.Lock) error {
	if entClient == nil {
		return fmt.Errorf("ent client not initialized")
	}

	spec, err := json.Marshal(lock.Spec)
	if err != nil {
		return fmt.Errorf("failed to marshal Lock spec: %w", err)
	}

	existing, err := entClient.Resource.Query().
		Where(
			entresource.KindEQ(lockKind),
			entresource.ResourceIDEQ(lock.Spec.ID),
		).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("failed to check Lock existence: %w", err)
	}

	if ent.IsNotFound(err) {
		_, err = entClient.Resource.Create().
			SetUID(lock.Metadata.UID).
			SetName(lock.Metadata.Name).
			SetAPIVersion(lock.APIVersion).
			SetKind(lockKind).
			SetResourceType(lockKind).
			SetResourceID(lock.Spec.ID).
			SetSpec(spec).
			SetCreatedAt(lock.Metadata.CreatedAt).
			SetUpdatedAt(lock.Metadata.UpdatedAt).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create Lock: %w", err)
		}
		return nil
	}

	_, err = entClient.Resource.UpdateOne(existing).
		SetName(lock.Metadata.Name).
		SetAPIVersion(lock.APIVersion).
		SetSpec(spec).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update Lock: %w", err)
	}
	return nil
}

// DeleteLock deletes a Lock by its UID.
func DeleteLock(ctx context.Context, uid string) error {
	if entClient == nil {
		return fmt.Errorf("ent client not initialized")
	}

	deleted, err := entClient.Resource.Delete().
		Where(
			entresource.UIDEQ(uid),
			entresource.KindEQ(lockKind),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Lock %s: %w", uid, err)
	}
	if deleted == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteLockByID deletes a Lock by its component xname (Spec.ID). It is a no-op
// (returns nil) when no lock exists for the component.
func DeleteLockByID(ctx context.Context, id string) error {
	if entClient == nil {
		return fmt.Errorf("ent client not initialized")
	}

	_, err := entClient.Resource.Delete().
		Where(
			entresource.KindEQ(lockKind),
			entresource.ResourceIDEQ(id),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Lock with ID %s: %w", id, err)
	}
	return nil
}
