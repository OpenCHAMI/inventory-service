/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the /hsm/v2/locks routes registered in csm_routes.go.
 * Routes under test:
 *   GET/POST /hsm/v2/locks/status
 *   POST     /hsm/v2/locks/lock | /unlock | /disable | /repair
 *   POST     /hsm/v2/locks/reservations | .../release | .../remove
 *   POST     /hsm/v2/locks/service/reservations | .../renew | .../release | .../check
 *
 * Lock/unlock and disable/repair mutate the Component (Spec.Locked /
 * Spec.ReservationDisabled); reservations are stored as separate Lock records
 * keyed by component xname. These tests exercise the SMD-flat request/response
 * shapes and verify that lock state is observable through /locks/status.
 */

package resttests

import (
	"net/http"
	"testing"
)

const csmLocksBase = "/hsm/v2/locks"

// ─── CSM lock request / response shapes ───────────────────────────────────────

type lockFilterReq struct {
	ComponentIDs    []string `json:"ComponentIDs,omitempty"`
	ProcessingModel string   `json:"ProcessingModel,omitempty"`
}

type lockCounts struct {
	Total   int `json:"Total"`
	Success int `json:"Success"`
	Failure int `json:"Failure"`
}

type failedXname struct {
	ID     string `json:"ID"`
	Reason string `json:"Reason"`
}

type xnameSuccess struct {
	ComponentIDs []string `json:"ComponentIDs"`
}

type xnameResponse struct {
	Counts  *lockCounts   `json:"Counts"`
	Success *xnameSuccess `json:"Success"`
	Failure []failedXname `json:"Failure"`
}

type reservationCreateReq struct {
	ComponentIDs        []string `json:"ComponentIDs,omitempty"`
	ProcessingModel     string   `json:"ProcessingModel,omitempty"`
	ReservationDuration int      `json:"ReservationDuration,omitempty"`
}

type reservationSuccess struct {
	ID             string `json:"ID"`
	DeputyKey      string `json:"DeputyKey"`
	ReservationKey string `json:"ReservationKey"`
	ExpirationTime string `json:"ExpirationTime,omitempty"`
}

type reservationCreateResp struct {
	Success []reservationSuccess `json:"Success"`
	Failure []failedXname        `json:"Failure"`
}

type xnameWithKey struct {
	ID  string `json:"ID"`
	Key string `json:"Key"`
}

type reservedKeysReq struct {
	ReservationKeys []xnameWithKey `json:"ReservationKeys"`
	ProcessingModel string         `json:"ProcessingModel,omitempty"`
}

type reservationRenewalReq struct {
	ReservationKeys     []xnameWithKey `json:"ReservationKeys"`
	ProcessingModel     string         `json:"ProcessingModel,omitempty"`
	ReservationDuration int            `json:"ReservationDuration,omitempty"`
}

type deputyKeysReq struct {
	DeputyKeys []xnameWithKey `json:"DeputyKeys"`
}

type reservationCheckResp struct {
	Success []reservationSuccess `json:"Success"`
	Failure []failedXname        `json:"Failure"`
}

type statusReq struct {
	ComponentIDs []string `json:"ComponentIDs,omitempty"`
}

type lockComponentStatus struct {
	ID                  string `json:"ID"`
	Locked              bool   `json:"Locked"`
	Reserved            bool   `json:"Reserved"`
	ReservationDisabled bool   `json:"ReservationDisabled"`
	CreatedTime         string `json:"CreatedTime,omitempty"`
	ExpirationTime      string `json:"ExpirationTime,omitempty"`
}

type lockStatusResp struct {
	Components []lockComponentStatus `json:"Components"`
	NotFound   []string              `json:"NotFound"`
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// lockPost sends a POST to a /hsm/v2/locks sub-path, decodes the response into
// dst (may be nil), and asserts HTTP 200.
func lockPost(t *testing.T, path string, body, dst interface{}) {
	t.Helper()
	resp := doRequest(t, http.MethodPost, csmLocksBase+path, body)
	requireStatus(t, resp, http.StatusOK)
	if dst != nil {
		decodeJSON(t, resp, dst)
	} else {
		resp.Body.Close()
	}
}

// lockStatusFor returns the /locks/status entry for a single xname, or nil if
// the component was reported NotFound.
func lockStatusFor(t *testing.T, xname string) *lockComponentStatus {
	t.Helper()
	var out lockStatusResp
	lockPost(t, "/status", statusReq{ComponentIDs: []string{xname}}, &out)
	for i := range out.Components {
		if out.Components[i].ID == xname {
			return &out.Components[i]
		}
	}
	return nil
}

// ─── Tests ────────────────────────────────────────────────────────────────────

// TestAdminLockUnlockCsm verifies that /locks/lock and /locks/unlock flip the
// component's Locked flag, observable through /locks/status.
func TestAdminLockUnlockCsm(t *testing.T) {
	xname := "x9000c0s0b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})
	defer csmDelete(t, xname)

	if st := lockStatusFor(t, xname); st == nil || st.Locked {
		t.Fatalf("expected component to start unlocked, got %+v", st)
	}

	var lockResp xnameResponse
	lockPost(t, "/lock", lockFilterReq{ComponentIDs: []string{xname}}, &lockResp)
	if lockResp.Counts == nil || lockResp.Counts.Success != 1 || lockResp.Counts.Failure != 0 {
		t.Fatalf("expected 1 success/0 failure locking, got %+v", lockResp.Counts)
	}
	if st := lockStatusFor(t, xname); st == nil || !st.Locked {
		t.Fatalf("expected component to be Locked after /lock, got %+v", st)
	}

	var unlockResp xnameResponse
	lockPost(t, "/unlock", lockFilterReq{ComponentIDs: []string{xname}}, &unlockResp)
	if unlockResp.Counts == nil || unlockResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success unlocking, got %+v", unlockResp.Counts)
	}
	if st := lockStatusFor(t, xname); st == nil || st.Locked {
		t.Fatalf("expected component to be Unlocked after /unlock, got %+v", st)
	}
}

// TestAdminReservationLifecycleCsm creates, observes and releases an admin
// (non-expiring) reservation.
func TestAdminReservationLifecycleCsm(t *testing.T) {
	xname := "x9000c0s1b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})
	defer csmDelete(t, xname)

	var createResp reservationCreateResp
	lockPost(t, "/reservations", reservationCreateReq{ComponentIDs: []string{xname}}, &createResp)
	if len(createResp.Success) != 1 {
		t.Fatalf("expected 1 reservation, got success=%+v failure=%+v", createResp.Success, createResp.Failure)
	}
	rsv := createResp.Success[0]
	if rsv.ReservationKey == "" || rsv.DeputyKey == "" {
		t.Fatalf("expected reservation and deputy keys, got %+v", rsv)
	}
	if rsv.ExpirationTime != "" {
		t.Fatalf("expected admin reservation to have no expiration, got %q", rsv.ExpirationTime)
	}

	if st := lockStatusFor(t, xname); st == nil || !st.Reserved {
		t.Fatalf("expected component to be Reserved, got %+v", st)
	}

	// A second reservation attempt must fail with Reserved.
	var second reservationCreateResp
	lockPost(t, "/reservations", reservationCreateReq{ComponentIDs: []string{xname}}, &second)
	if len(second.Failure) != 1 || second.Failure[0].Reason != "Reserved" {
		t.Fatalf("expected Reserved failure on double reservation, got %+v", second)
	}

	// Release with the reservation key.
	var releaseResp xnameResponse
	lockPost(t, "/reservations/release", reservedKeysReq{
		ReservationKeys: []xnameWithKey{{ID: xname, Key: rsv.ReservationKey}},
	}, &releaseResp)
	if releaseResp.Counts == nil || releaseResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success releasing, got %+v", releaseResp)
	}
	if st := lockStatusFor(t, xname); st == nil || st.Reserved {
		t.Fatalf("expected component to be un-reserved after release, got %+v", st)
	}
}

// TestServiceReservationCheckRenewCsm exercises the service reservation flow:
// create (timed), check by deputy key, renew, then release.
func TestServiceReservationCheckRenewCsm(t *testing.T) {
	xname := "x9000c0s2b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})
	defer csmDelete(t, xname)

	var createResp reservationCreateResp
	lockPost(t, "/service/reservations", reservationCreateReq{
		ComponentIDs:        []string{xname},
		ReservationDuration: 2,
	}, &createResp)
	if len(createResp.Success) != 1 {
		t.Fatalf("expected 1 service reservation, got %+v / %+v", createResp.Success, createResp.Failure)
	}
	rsv := createResp.Success[0]
	if rsv.ExpirationTime == "" {
		t.Fatalf("expected service reservation to have an expiration time")
	}

	// Validate the deputy key.
	var checkResp reservationCheckResp
	lockPost(t, "/service/reservations/check", deputyKeysReq{
		DeputyKeys: []xnameWithKey{{ID: xname, Key: rsv.DeputyKey}},
	}, &checkResp)
	if len(checkResp.Success) != 1 || checkResp.Success[0].ID != xname {
		t.Fatalf("expected deputy key check to succeed, got %+v", checkResp)
	}

	// Renew with the reservation key.
	var renewResp xnameResponse
	lockPost(t, "/service/reservations/renew", reservationRenewalReq{
		ReservationKeys:     []xnameWithKey{{ID: xname, Key: rsv.ReservationKey}},
		ReservationDuration: 5,
	}, &renewResp)
	if renewResp.Counts == nil || renewResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success renewing, got %+v", renewResp)
	}

	// Release.
	var releaseResp xnameResponse
	lockPost(t, "/service/reservations/release", reservedKeysReq{
		ReservationKeys: []xnameWithKey{{ID: xname, Key: rsv.ReservationKey}},
	}, &releaseResp)
	if releaseResp.Counts == nil || releaseResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success releasing service reservation, got %+v", releaseResp)
	}
}

// TestReservationDisabledRejectedCsm verifies /disable sets ReservationDisabled
// and that reservations are then rejected with the Disabled reason, and that
// /repair clears the flag.
func TestReservationDisabledRejectedCsm(t *testing.T) {
	xname := "x9000c0s3b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})
	defer csmDelete(t, xname)

	var disableResp xnameResponse
	lockPost(t, "/disable", lockFilterReq{ComponentIDs: []string{xname}}, &disableResp)
	if disableResp.Counts == nil || disableResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success disabling, got %+v", disableResp)
	}
	if st := lockStatusFor(t, xname); st == nil || !st.ReservationDisabled {
		t.Fatalf("expected ReservationDisabled after /disable, got %+v", st)
	}

	var createResp reservationCreateResp
	lockPost(t, "/reservations", reservationCreateReq{ComponentIDs: []string{xname}}, &createResp)
	if len(createResp.Failure) != 1 || createResp.Failure[0].Reason != "Disabled" {
		t.Fatalf("expected Disabled failure, got %+v", createResp)
	}

	var repairResp xnameResponse
	lockPost(t, "/repair", lockFilterReq{ComponentIDs: []string{xname}}, &repairResp)
	if repairResp.Counts == nil || repairResp.Counts.Success != 1 {
		t.Fatalf("expected 1 success repairing, got %+v", repairResp)
	}
	if st := lockStatusFor(t, xname); st == nil || st.ReservationDisabled {
		t.Fatalf("expected ReservationDisabled cleared after /repair, got %+v", st)
	}
}

// TestLockStatusNotFoundCsm verifies that an unknown xname is reported in the
// NotFound list of /locks/status.
func TestLockStatusNotFoundCsm(t *testing.T) {
	var out lockStatusResp
	lockPost(t, "/status", statusReq{ComponentIDs: []string{"xDOESNOTEXIST0"}}, &out)
	found := false
	for _, id := range out.NotFound {
		if id == "xDOESNOTEXIST0" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected unknown xname in NotFound, got %+v", out)
	}
}

// TestDeleteComponentRemovesLockCsm verifies that deleting a component also
// removes its reservation (no orphaned Lock record).
func TestDeleteComponentRemovesLockCsm(t *testing.T) {
	xname := "x9000c0s4b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})

	var createResp reservationCreateResp
	lockPost(t, "/reservations", reservationCreateReq{ComponentIDs: []string{xname}}, &createResp)
	if len(createResp.Success) != 1 {
		t.Fatalf("expected reservation before delete, got %+v", createResp)
	}

	csmDelete(t, xname)

	// Recreate the component; its lock must be gone (Reserved == false).
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "Ready"})
	defer csmDelete(t, xname)
	if st := lockStatusFor(t, xname); st == nil || st.Reserved {
		t.Fatalf("expected no reservation after component delete/recreate, got %+v", st)
	}
}
