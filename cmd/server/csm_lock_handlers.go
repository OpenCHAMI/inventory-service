// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/openchami/fabrica/pkg/events"
	"github.com/openchami/fabrica/pkg/resource"
	v1 "github.com/openchami/inventory-service/apis/inventory-service.openchami.org/v1"
	"github.com/openchami/inventory-service/cmd/plugins"
	"github.com/openchami/inventory-service/internal/storage"
)

// lockFilterSet is the internal, normalized set of component filters shared by
// the /hsm/v2/locks endpoints. Partition filtering is intentionally not part of
// this set: it is rejected at the handler boundary because inventory-service has
// no partition model (see the package README / SMD-compat notes).
type lockFilterSet struct {
	Type    []string
	State   []string
	Role    []string
	SubRole []string
	Class   []string
	Group   []string
}

// ---- generic helpers -------------------------------------------------------

func containsFold(list []string, val string) bool {
	for _, s := range list {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}

func anyFold(a, b []string) bool {
	for _, x := range a {
		if containsFold(b, x) {
			return true
		}
	}
	return false
}

// isRigid reports whether the processing model is the default all-or-nothing
// "rigid" model (anything other than an explicit "flexible" is rigid).
func isRigid(pm string) bool {
	return !strings.EqualFold(pm, "flexible")
}

// reservationDuration clamps a requested service-reservation duration (minutes)
// to the SMD-supported 1-15 minute range.
func reservationDuration(mins int) time.Duration {
	if mins <= 0 {
		mins = 1
	}
	if mins > 15 {
		mins = 15
	}
	return time.Duration(mins) * time.Minute
}

// newLockKey builds an SMD-style reservation/deputy key of the form
// "<xname>:<uuid>".
func newLockKey(id string) string {
	return id + ":" + uuid.NewString()
}

func newXnameResponse(success []string, failure []FailedXname) XnameResponse {
	if success == nil {
		success = []string{}
	}
	if failure == nil {
		failure = []FailedXname{}
	}
	return XnameResponse{
		Counts: &LockCounts{
			Total:   len(success) + len(failure),
			Success: len(success),
			Failure: len(failure),
		},
		Success: &XnameSuccess{ComponentIDs: success},
		Failure: failure,
	}
}

// ---- component / group resolution -----------------------------------------

// lockGroupMap builds a component-xname -> group-labels index used for the
// Group filter.
func lockGroupMap(ctx context.Context) (map[string][]string, error) {
	groups, err := plugins.Store.LoadAllGroups(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	for _, g := range groups {
		for _, id := range g.Spec.Members.IDs {
			m[id] = append(m[id], g.Spec.Label)
		}
	}
	return m, nil
}

func componentMatchesLockFilters(c *v1.Component, groups []string, f lockFilterSet) bool {
	if len(f.Type) > 0 && !containsFold(f.Type, c.Spec.Type) {
		return false
	}
	if len(f.State) > 0 && !containsFold(f.State, c.Spec.State) {
		return false
	}
	if len(f.Role) > 0 && !containsFold(f.Role, c.Spec.Role) {
		return false
	}
	if len(f.SubRole) > 0 && !containsFold(f.SubRole, c.Spec.SubRole) {
		return false
	}
	if len(f.Class) > 0 && !containsFold(f.Class, c.Spec.Class) {
		return false
	}
	if len(f.Group) > 0 && !anyFold(f.Group, groups) {
		return false
	}
	return true
}

// resolveLockTargets returns the components an operation should act on. When
// ids is non-empty, unknown ids are returned in notFound; otherwise every
// component matching the filters is returned.
func resolveLockTargets(r *http.Request, ids []string, f lockFilterSet) (targets []*v1.Component, notFound []string, err error) {
	ctx := r.Context()
	groupMap, err := lockGroupMap(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(ids) > 0 {
		for _, id := range ids {
			comp, err := plugins.Store.LoadComponentByID(ctx, id)
			if err == storage.ErrNotFound || comp == nil {
				notFound = append(notFound, id)
				continue
			}
			if err != nil {
				return nil, nil, err
			}
			if componentMatchesLockFilters(comp, groupMap[id], f) {
				targets = append(targets, comp)
			}
		}
		return targets, notFound, nil
	}

	comps, err := plugins.Store.LoadAllComponents(ctx)
	if err != nil {
		return nil, nil, err
	}
	for _, comp := range comps {
		if componentMatchesLockFilters(comp, groupMap[comp.Spec.ID], f) {
			targets = append(targets, comp)
		}
	}
	return targets, notFound, nil
}

// ---- lock/reservation helpers ---------------------------------------------

func lockExpired(lock *v1.Lock, now time.Time) bool {
	if lock == nil || lock.Spec.ExpirationTime == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, lock.Spec.ExpirationTime)
	if err != nil {
		return false
	}
	return now.After(t)
}

// activeLock returns the current, non-expired reservation for a component, or
// nil if there is none. Expired reservations are lazily deleted.
func activeLock(r *http.Request, id string, now time.Time) (*v1.Lock, error) {
	lock, err := plugins.Store.LoadLockByID(r.Context(), id)
	if err == storage.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if lockExpired(lock, now) {
		_ = plugins.Store.DeleteLockByID(r.Context(), id)
		return nil, nil
	}
	return lock, nil
}

// saveComponentFlag persists a component after an admin lock flag change and
// publishes an update event.
func saveComponentFlag(r *http.Request, comp *v1.Component) error {
	comp.Metadata.UpdatedAt = time.Now()
	if err := plugins.Store.SaveComponent(r.Context(), comp); err != nil {
		return err
	}
	if err := events.PublishResourceUpdated(r.Context(), "Component", comp.Metadata.UID, comp.Metadata.Name, comp, map[string]interface{}{"updatedAt": comp.Metadata.UpdatedAt}); err != nil {
		fmt.Printf("Warning: Failed to publish resource updated event for Component %s: %v\n", comp.Metadata.UID, err)
	}
	return nil
}

// ---- admin lock flag endpoints (lock/unlock/disable/repair) ----------------

// handleAdminLockFlag implements the shared body of the admin lock/unlock/
// disable/repair endpoints. setFlag mutates the component; when
// removeReservation is true any active reservation is also deleted (used by
// /disable).
func handleAdminLockFlag(w http.ResponseWriter, r *http.Request, setFlag func(*v1.Component), removeReservation bool) {
	var req LockFilter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if len(req.Partition) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Partition filter is not supported by inventory-service"))
		return
	}

	filters := lockFilterSet{Type: req.Type, State: req.State, Role: req.Role, SubRole: req.SubRole, Class: req.Class, Group: req.Group}
	targets, notFound, err := resolveLockTargets(r, req.ComponentIDs, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to resolve components: %w", err))
		return
	}

	var failure []FailedXname
	for _, id := range notFound {
		failure = append(failure, FailedXname{ID: id, Reason: lockReasonNotFound})
	}

	// Rigid processing: any failure aborts the whole request with no changes.
	if isRigid(req.ProcessingModel) && len(failure) > 0 {
		respondJSON(w, http.StatusOK, newXnameResponse(nil, failure))
		return
	}

	var success []string
	for _, comp := range targets {
		setFlag(comp)
		if err := saveComponentFlag(r, comp); err != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		if removeReservation {
			_ = plugins.Store.DeleteLockByID(r.Context(), comp.Spec.ID)
		}
		success = append(success, comp.Spec.ID)
	}
	respondJSON(w, http.StatusOK, newXnameResponse(success, failure))
}

// LockComponentsCsm handles POST /hsm/v2/locks/lock.
func LockComponentsCsm(w http.ResponseWriter, r *http.Request) {
	handleAdminLockFlag(w, r, func(c *v1.Component) { c.Spec.Locked = true }, false)
}

// UnlockComponentsCsm handles POST /hsm/v2/locks/unlock.
func UnlockComponentsCsm(w http.ResponseWriter, r *http.Request) {
	handleAdminLockFlag(w, r, func(c *v1.Component) { c.Spec.Locked = false }, false)
}

// DisableComponentsCsm handles POST /hsm/v2/locks/disable. Disabling
// reservations also removes any existing reservation.
func DisableComponentsCsm(w http.ResponseWriter, r *http.Request) {
	handleAdminLockFlag(w, r, func(c *v1.Component) { c.Spec.ReservationDisabled = true }, true)
}

// RepairComponentsCsm handles POST /hsm/v2/locks/repair.
func RepairComponentsCsm(w http.ResponseWriter, r *http.Request) {
	handleAdminLockFlag(w, r, func(c *v1.Component) { c.Spec.ReservationDisabled = false }, false)
}

// ---- reservation create (admin + service) ----------------------------------

// handleCreateReservations implements POST /hsm/v2/locks/reservations (admin,
// service=false, no expiration) and POST /hsm/v2/locks/service/reservations
// (service=true, timed expiration). A component may be reserved only when it is
// Unlocked, not ReservationDisabled and not already reserved.
func handleCreateReservations(w http.ResponseWriter, r *http.Request, service bool) {
	var req ReservationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if len(req.Partition) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Partition filter is not supported by inventory-service"))
		return
	}

	filters := lockFilterSet{Type: req.Type, State: req.State, Role: req.Role, SubRole: req.SubRole, Class: req.Class, Group: req.Group}
	targets, notFound, err := resolveLockTargets(r, req.ComponentIDs, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to resolve components: %w", err))
		return
	}

	now := time.Now()
	var failure []FailedXname
	for _, id := range notFound {
		failure = append(failure, FailedXname{ID: id, Reason: lockReasonNotFound})
	}

	var eligible []*v1.Component
	for _, comp := range targets {
		if comp.Spec.ReservationDisabled {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonDisabled})
			continue
		}
		if comp.Spec.Locked {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonLocked})
			continue
		}
		lock, err := activeLock(r, comp.Spec.ID, now)
		if err != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		if lock != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonReserved})
			continue
		}
		eligible = append(eligible, comp)
	}

	if isRigid(req.ProcessingModel) && len(failure) > 0 {
		respondJSON(w, http.StatusOK, ReservationCreateResponse{Success: []ReservationCreateSuccess{}, Failure: failure})
		return
	}

	var success []ReservationCreateSuccess
	for _, comp := range eligible {
		expiration := ""
		if service {
			expiration = now.Add(reservationDuration(req.ReservationDuration)).Format(time.RFC3339)
		}
		rk := newLockKey(comp.Spec.ID)
		dk := newLockKey(comp.Spec.ID)

		uid, err := resource.GenerateUIDForResource("Lock")
		if err != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		lock := &v1.Lock{
			APIVersion: comp.APIVersion,
			Kind:       "Lock",
			Spec: v1.LockSpec{
				ID:             comp.Spec.ID,
				ReservationKey: rk,
				DeputyKey:      dk,
				CreatedTime:    now.Format(time.RFC3339),
				ExpirationTime: expiration,
			},
		}
		lock.Metadata.Initialize(comp.Spec.ID, uid)

		if err := plugins.Store.SaveLock(r.Context(), lock); err != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		success = append(success, ReservationCreateSuccess{
			ID:             comp.Spec.ID,
			DeputyKey:      dk,
			ReservationKey: rk,
			ExpirationTime: expiration,
		})
	}

	if success == nil {
		success = []ReservationCreateSuccess{}
	}
	if failure == nil {
		failure = []FailedXname{}
	}
	respondJSON(w, http.StatusOK, ReservationCreateResponse{Success: success, Failure: failure})
}

// CreateReservationsCsm handles POST /hsm/v2/locks/reservations.
func CreateReservationsCsm(w http.ResponseWriter, r *http.Request) {
	handleCreateReservations(w, r, false)
}

// CreateServiceReservationsCsm handles POST /hsm/v2/locks/service/reservations.
func CreateServiceReservationsCsm(w http.ResponseWriter, r *http.Request) {
	handleCreateReservations(w, r, true)
}

// ---- reservation release / remove / renew / check --------------------------

// handleReleaseReservations implements the key-matched release used by both
// /locks/reservations/release and /locks/service/reservations/release.
func handleReleaseReservations(w http.ResponseWriter, r *http.Request) {
	var req ReservedKeys
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	now := time.Now()
	var failure []FailedXname
	var toRelease []string
	for _, k := range req.ReservationKeys {
		lock, err := activeLock(r, k.ID, now)
		if err != nil {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonServerError})
			continue
		}
		if lock == nil || lock.Spec.ReservationKey != k.Key {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonNotFound})
			continue
		}
		toRelease = append(toRelease, k.ID)
	}

	if isRigid(req.ProcessingModel) && len(failure) > 0 {
		respondJSON(w, http.StatusOK, newXnameResponse(nil, failure))
		return
	}

	var success []string
	for _, id := range toRelease {
		if err := plugins.Store.DeleteLockByID(r.Context(), id); err != nil {
			failure = append(failure, FailedXname{ID: id, Reason: lockReasonServerError})
			continue
		}
		success = append(success, id)
	}
	respondJSON(w, http.StatusOK, newXnameResponse(success, failure))
}

// ReleaseReservationsCsm handles POST /hsm/v2/locks/reservations/release.
func ReleaseReservationsCsm(w http.ResponseWriter, r *http.Request) {
	handleReleaseReservations(w, r)
}

// ReleaseServiceReservationsCsm handles POST /hsm/v2/locks/service/reservations/release.
func ReleaseServiceReservationsCsm(w http.ResponseWriter, r *http.Request) {
	handleReleaseReservations(w, r)
}

// RemoveReservationsCsm handles POST /hsm/v2/locks/reservations/remove. It
// forcibly deletes reservations by xname, ignoring keys.
func RemoveReservationsCsm(w http.ResponseWriter, r *http.Request) {
	var req LockFilter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if len(req.Partition) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Partition filter is not supported by inventory-service"))
		return
	}

	filters := lockFilterSet{Type: req.Type, State: req.State, Role: req.Role, SubRole: req.SubRole, Class: req.Class, Group: req.Group}
	targets, notFound, err := resolveLockTargets(r, req.ComponentIDs, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to resolve components: %w", err))
		return
	}

	now := time.Now()
	var failure []FailedXname
	for _, id := range notFound {
		failure = append(failure, FailedXname{ID: id, Reason: lockReasonNotFound})
	}

	// Pre-check which targets actually have a reservation to remove.
	var toRemove []string
	for _, comp := range targets {
		lock, err := activeLock(r, comp.Spec.ID, now)
		if err != nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		if lock == nil {
			failure = append(failure, FailedXname{ID: comp.Spec.ID, Reason: lockReasonNotFound})
			continue
		}
		toRemove = append(toRemove, comp.Spec.ID)
	}

	if isRigid(req.ProcessingModel) && len(failure) > 0 {
		respondJSON(w, http.StatusOK, newXnameResponse(nil, failure))
		return
	}

	var success []string
	for _, id := range toRemove {
		if err := plugins.Store.DeleteLockByID(r.Context(), id); err != nil {
			failure = append(failure, FailedXname{ID: id, Reason: lockReasonServerError})
			continue
		}
		success = append(success, id)
	}
	respondJSON(w, http.StatusOK, newXnameResponse(success, failure))
}

// RenewServiceReservationsCsm handles POST /hsm/v2/locks/service/reservations/renew.
func RenewServiceReservationsCsm(w http.ResponseWriter, r *http.Request) {
	var req ReservationRenewal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	now := time.Now()
	var failure []FailedXname
	var toRenew []*v1.Lock
	for _, k := range req.ReservationKeys {
		lock, err := activeLock(r, k.ID, now)
		if err != nil {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonServerError})
			continue
		}
		if lock == nil || lock.Spec.ReservationKey != k.Key {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonNotFound})
			continue
		}
		toRenew = append(toRenew, lock)
	}

	if isRigid(req.ProcessingModel) && len(failure) > 0 {
		respondJSON(w, http.StatusOK, newXnameResponse(nil, failure))
		return
	}

	var success []string
	for _, lock := range toRenew {
		lock.Spec.ExpirationTime = now.Add(reservationDuration(req.ReservationDuration)).Format(time.RFC3339)
		if err := plugins.Store.SaveLock(r.Context(), lock); err != nil {
			failure = append(failure, FailedXname{ID: lock.Spec.ID, Reason: lockReasonServerError})
			continue
		}
		success = append(success, lock.Spec.ID)
	}
	respondJSON(w, http.StatusOK, newXnameResponse(success, failure))
}

// CheckServiceReservationsCsm handles POST /hsm/v2/locks/service/reservations/check.
// It validates deputy keys and returns the matching reservation details.
func CheckServiceReservationsCsm(w http.ResponseWriter, r *http.Request) {
	var req DeputyKeys
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	now := time.Now()
	var success []ReservationCheckSuccess
	var failure []FailedXname
	for _, k := range req.DeputyKeys {
		lock, err := activeLock(r, k.ID, now)
		if err != nil {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonServerError})
			continue
		}
		if lock == nil || lock.Spec.DeputyKey != k.Key {
			failure = append(failure, FailedXname{ID: k.ID, Reason: lockReasonNotFound})
			continue
		}
		success = append(success, ReservationCheckSuccess{
			ID:             k.ID,
			DeputyKey:      lock.Spec.DeputyKey,
			ReservationKey: lock.Spec.ReservationKey,
			ExpirationTime: lock.Spec.ExpirationTime,
		})
	}

	if success == nil {
		success = []ReservationCheckSuccess{}
	}
	if failure == nil {
		failure = []FailedXname{}
	}
	respondJSON(w, http.StatusOK, ReservationCheckResponse{Success: success, Failure: failure})
}

// ---- status ----------------------------------------------------------------

func writeLockStatus(w http.ResponseWriter, r *http.Request, ids []string, filters lockFilterSet) {
	targets, notFound, err := resolveLockTargets(r, ids, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to resolve components: %w", err))
		return
	}

	now := time.Now()
	comps := make([]LockComponentStatus, 0, len(targets))
	for _, comp := range targets {
		lock, err := activeLock(r, comp.Spec.ID, now)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("failed to load lock for %s: %w", comp.Spec.ID, err))
			return
		}
		st := LockComponentStatus{
			ID:                  comp.Spec.ID,
			Locked:              comp.Spec.Locked,
			ReservationDisabled: comp.Spec.ReservationDisabled,
			Reserved:            lock != nil,
		}
		if lock != nil {
			st.CreatedTime = lock.Spec.CreatedTime
			st.ExpirationTime = lock.Spec.ExpirationTime
		}
		comps = append(comps, st)
	}
	respondJSON(w, http.StatusOK, AdminStatusCheckResponse{Components: comps, NotFound: notFound})
}

func queryValues(q url.Values, names ...string) []string {
	var out []string
	for _, n := range names {
		out = append(out, q[n]...)
	}
	return out
}

// GetLockStatusCsm handles GET /hsm/v2/locks/status. Filters are supplied as
// (repeatable) query parameters; both lower-case and SMD-style capitalized keys
// are accepted.
func GetLockStatusCsm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(queryValues(q, "partition", "Partition")) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Partition filter is not supported by inventory-service"))
		return
	}
	filters := lockFilterSet{
		Type:    queryValues(q, "type", "Type"),
		State:   queryValues(q, "state", "State"),
		Role:    queryValues(q, "role", "Role"),
		SubRole: queryValues(q, "subrole", "SubRole"),
		Class:   queryValues(q, "class", "Class"),
		Group:   queryValues(q, "group", "Group"),
	}
	ids := queryValues(q, "componentid", "ComponentID", "id")
	writeLockStatus(w, r, ids, filters)
}

// PostLockStatusCsm handles POST /hsm/v2/locks/status.
func PostLockStatusCsm(w http.ResponseWriter, r *http.Request) {
	var req AdminStatusCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if len(req.Partition) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Partition filter is not supported by inventory-service"))
		return
	}
	filters := lockFilterSet{Type: req.Type, State: req.State, Role: req.Role, SubRole: req.SubRole, Class: req.Class, Group: req.Group}
	writeLockStatus(w, r, req.ComponentIDs, filters)
}
