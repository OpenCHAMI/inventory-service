/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the service metadata routes registered directly in
 * cmd/server/main.go (public, unauthenticated endpoints).
 * Routes under test:
 *   GET /openapi.json
 *   GET /docs
 */

package resttests

import (
	"net/http"
	"strings"
	"testing"
)

// TestGetOpenAPISpec verifies GET /openapi.json returns HTTP 200 with a
// valid OpenAPI 3.0 JSON document.
func TestGetOpenAPISpec(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/openapi.json", nil)
	requireStatus(t, resp, http.StatusOK)

	var spec struct {
		OpenAPI string `json:"openapi"`
		Info    struct {
			Title string `json:"title"`
		} `json:"info"`
		Paths map[string]any `json:"paths"`
	}
	decodeJSON(t, resp, &spec)

	if !strings.HasPrefix(spec.OpenAPI, "3.") {
		t.Errorf("expected openapi version 3.x, got %q", spec.OpenAPI)
	}
	if spec.Info.Title == "" {
		t.Error("expected non-empty info.title")
	}
	if len(spec.Paths) == 0 {
		t.Error("expected at least one path in the OpenAPI spec")
	}
	if _, ok := spec.Paths["/components"]; !ok {
		t.Error("expected /components to be documented in the OpenAPI spec")
	}
}

// TestGetSwaggerUI verifies GET /docs returns HTTP 200 with an HTML page.
func TestGetSwaggerUI(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/docs", nil)
	requireStatus(t, resp, http.StatusOK)
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type to contain text/html, got %q", contentType)
	}
}
