/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the unprotected CSM service status routes registered in
 * csm_routes.go (RegisterUnprotectedCsmRoutes).
 * Routes under test:
 *   GET /hsm/v2/service/ready
 *   GET /hsm/v2/service/liveness
 */

package resttests

import (
	"net/http"
	"testing"
)

// TestGetReadinessCsm verifies GET /hsm/v2/service/ready returns HTTP 200
// with a ServiceReadiness body.
func TestGetReadinessCsm(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/hsm/v2/service/ready", nil)
	requireStatus(t, resp, http.StatusOK)

	var readiness struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	decodeJSON(t, resp, &readiness)
	if readiness.Code != 0 {
		t.Errorf("expected Code=0, got %d", readiness.Code)
	}
	if readiness.Message == "" {
		t.Error("expected non-empty Message")
	}
}

// TestGetLivenessCsm verifies GET /hsm/v2/service/liveness returns HTTP 204.
func TestGetLivenessCsm(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/hsm/v2/service/liveness", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected HTTP 204, got %d", resp.StatusCode)
	}
}
