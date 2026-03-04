// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT
package main

import "github.com/go-chi/chi/v5"

func RegisterCsmRoutes(r chi.Router) {

	// Component routes
	r.Route("/hsm/v2/State/Components", func(r chi.Router) {
		r.Get("/", GetComponentsCsm)
		r.Post("/", CreateComponentCsm)
		// r.Delete("/", DeleteAllComponentCsm) // todo (smd has it but maybe not needed)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetComponentCsm)
			r.Put("/", UpdateComponentCsm)
			r.Delete("/", DeleteComponentCsm)
		})
		// possible todo
		// Get /State/Components/ByNID/{nid}
		// Get /State/Components/Query/{xname}
		// Post /State/Components/Query
		// Post /State/Components/ByNID/Query
		// Patch /State/Components/BulkStateData
		// Patch /State/Components/{xname}/StateData
		// Patch /State/Components/{xname}/FlagOnly
		// Patch /State/Components/{xname}/Enabled
		// Patch /State/Components/{xname}/SoftwareStatus
		// Patch /State/Components/{xname}/Role
		// Patch /State/Components/{xname}/NID
		// Patch /State/Components/BulkStateData
		// Patch /State/Components/BulkFlagOnly
		// Patch /State/Components/BulkEnabled
		// Patch /State/Components/BulkSoftwareStatus
		// Patch /State/Components/BulkRole
		// Patch /State/Components/BulkNID
	})
	// ComponentEndpoint routes
	r.Route("/hsm/v2/Inventory/ComponentEndpoints", func(r chi.Router) {
		r.Get("/", GetComponentEndpointsCsm)
		r.Post("/", CreateComponentEndpointCsm)
		// todo (optional)
		// DELETE /Inventory/ComponentEndpoints
		// r.Delete("/", DeleteAllComponentCsm) // todo (smd has it but it is probably not needed)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetComponentEndpointCsm)
			r.Put("/", UpdateComponentEndpointCsm)
			r.Delete("/", DeleteComponentEndpointCsm)
		})
	})
	// ServiceEndpoint routes
	r.Route("/hsm/v2/Inventory/ServiceEndpoints", func(r chi.Router) {
		r.Get("/", GetServiceEndpointsCsm)
		r.Post("/", CreateServiceEndpointCsm)
		// todo (optional)
		// DELETE /Inventory/ServiceEndpoints
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetServiceEndpointCsm)
			r.Put("/", UpdateServiceEndpointCsm)
			r.Delete("/", DeleteServiceEndpointCsm)
			// GET /ServiceEndpoints/{id}/RedfishEndpoints/{redfish_endpoint_id}
			// DELETE /ServiceEndpoints/{id}/RedfishEndpoints/{redfish_endpoint_id}
		})
	})
	// RedfishEndpoint routes
	r.Route("/hsm/v2/Inventory/RedfishEndpoints", func(r chi.Router) {
		r.Get("/", GetRedfishEndpointsCsm)
		r.Post("/", CreateRedfishEndpointCsm)
		// todo (optional)
		// DELETE /Inventory/RedfishEndpoints
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetRedfishEndpointCsm)
			r.Put("/", UpdateRedfishEndpointV2)
			r.Delete("/", DeleteRedfishEndpointV2)
		})
	})
	// EthernetInterface routes
	r.Route("/hsm/v2/Inventory/EthernetInterfaces", func(r chi.Router) {
		r.Get("/", GetEthernetInterfacesCsm)
		r.Post("/", CreateEthernetInterfaceCsm)
		// todo (optional)
		// DELETE /Inventory/EthernetInterfaces
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetEthernetInterfaceCsm)
			r.Put("/", UpdateEthernetInterfaceCsm)
			r.Delete("/", DeleteEthernetInterfaceCsm)
		})
	})
	// Group routes
	r.Route("/hsm/v2/groups", func(r chi.Router) {
		r.Get("/", GetGroupsCsm)
		r.Post("/", CreateGroupCsm)
		r.Route("/{group_label}", func(r chi.Router) {
			r.Get("/", GetGroupCsm)
			r.Put("/", UpdateGroupCsm)
			r.Delete("/", DeleteGroupCsm)
			// GET /groups/{group_label}/members
			// POST /groups/{group_label}/members
			// PUT /groups/{group_label}/members
			// DELETE /groups/{group_label}/members/{member_id}
		})
	})
}
