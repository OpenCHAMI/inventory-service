// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT
package main

// SMD-compatible request/response types for the /hsm/v2/locks endpoints.
// These mirror the shapes used by OpenCHAMI/smd so existing clients continue to
// work unmodified.

// Reason codes reported for a component that could not be processed. These match
// the values used by SMD.
const (
	lockReasonNotFound    = "NotFound"
	lockReasonLocked      = "Locked"
	lockReasonDisabled    = "Disabled"
	lockReasonReserved    = "Reserved"
	lockReasonServerError = "ServerError"
)

// processingModelRigid is the default all-or-nothing processing model; the
// alternative is "flexible" (best-effort, per-component).
const processingModelRigid = "rigid"

// LockFilter is the common request body for the admin lock/unlock/disable/repair
// and reservation-remove operations. Targets are the union of ComponentIDs and
// any components matching the supplied filters.
type LockFilter struct {
	ComponentIDs    []string `json:"ComponentIDs,omitempty"`
	Partition       []string `json:"Partition,omitempty"`
	Group           []string `json:"Group,omitempty"`
	Type            []string `json:"Type,omitempty"`
	State           []string `json:"State,omitempty"`
	Role            []string `json:"Role,omitempty"`
	SubRole         []string `json:"SubRole,omitempty"`
	Class           []string `json:"Class,omitempty"`
	ProcessingModel string   `json:"ProcessingModel,omitempty"`
}

// LockCounts summarizes the outcome of an admin lock operation.
type LockCounts struct {
	Total   int `json:"Total"`
	Success int `json:"Success"`
	Failure int `json:"Failure"`
}

// FailedXname is a single component that could not be processed.
type FailedXname struct {
	ID     string `json:"ID"`
	Reason string `json:"Reason"`
}

// XnameSuccess lists the components that were successfully processed.
type XnameSuccess struct {
	ComponentIDs []string `json:"ComponentIDs"`
}

// XnameResponse is the standard response for admin lock operations that do not
// return reservation keys.
type XnameResponse struct {
	Counts  *LockCounts   `json:"Counts,omitempty"`
	Success *XnameSuccess `json:"Success,omitempty"`
	Failure []FailedXname `json:"Failure,omitempty"`
}

// ReservationCreateRequest is the request body for creating admin and service
// reservations. ReservationDuration (minutes) only applies to service
// reservations; admin reservations never expire.
type ReservationCreateRequest struct {
	ComponentIDs        []string `json:"ComponentIDs,omitempty"`
	Partition           []string `json:"Partition,omitempty"`
	Group               []string `json:"Group,omitempty"`
	Type                []string `json:"Type,omitempty"`
	State               []string `json:"State,omitempty"`
	Role                []string `json:"Role,omitempty"`
	SubRole             []string `json:"SubRole,omitempty"`
	Class               []string `json:"Class,omitempty"`
	ProcessingModel     string   `json:"ProcessingModel,omitempty"`
	ReservationDuration int      `json:"ReservationDuration,omitempty"`
}

// ReservationCreateSuccess is a single successfully created reservation.
type ReservationCreateSuccess struct {
	ID             string `json:"ID"`
	DeputyKey      string `json:"DeputyKey"`
	ReservationKey string `json:"ReservationKey"`
	ExpirationTime string `json:"ExpirationTime,omitempty"`
}

// ReservationCreateResponse is returned by the reservation-create endpoints.
type ReservationCreateResponse struct {
	Success []ReservationCreateSuccess `json:"Success"`
	Failure []FailedXname              `json:"Failure"`
}

// XnameWithKey pairs a component xname with a reservation or deputy key.
type XnameWithKey struct {
	ID  string `json:"ID"`
	Key string `json:"Key"`
}

// ReservedKeys is the request body for releasing reservations by key.
type ReservedKeys struct {
	ReservationKeys []XnameWithKey `json:"ReservationKeys"`
	ProcessingModel string         `json:"ProcessingModel,omitempty"`
}

// ReservationRenewal is the request body for renewing service reservations.
type ReservationRenewal struct {
	ReservationKeys     []XnameWithKey `json:"ReservationKeys"`
	ProcessingModel     string         `json:"ProcessingModel,omitempty"`
	ReservationDuration int            `json:"ReservationDuration,omitempty"`
}

// DeputyKeys is the request body for validating service reservations.
type DeputyKeys struct {
	DeputyKeys []XnameWithKey `json:"DeputyKeys"`
}

// ReservationCheckSuccess is a single validated reservation.
type ReservationCheckSuccess struct {
	ID             string `json:"ID"`
	DeputyKey      string `json:"DeputyKey"`
	ReservationKey string `json:"ReservationKey"`
	ExpirationTime string `json:"ExpirationTime,omitempty"`
}

// ReservationCheckResponse is returned by the service reservation check endpoint.
type ReservationCheckResponse struct {
	Success []ReservationCheckSuccess `json:"Success"`
	Failure []FailedXname             `json:"Failure"`
}

// AdminStatusCheckRequest is the POST body for /locks/status.
type AdminStatusCheckRequest struct {
	ComponentIDs []string `json:"ComponentIDs,omitempty"`
	Partition    []string `json:"Partition,omitempty"`
	Group        []string `json:"Group,omitempty"`
	Type         []string `json:"Type,omitempty"`
	State        []string `json:"State,omitempty"`
	Role         []string `json:"Role,omitempty"`
	SubRole      []string `json:"SubRole,omitempty"`
	Class        []string `json:"Class,omitempty"`
}

// LockComponentStatus is the lock/reservation status for a single component.
type LockComponentStatus struct {
	ID                  string `json:"ID"`
	Locked              bool   `json:"Locked"`
	Reserved            bool   `json:"Reserved"`
	ReservationDisabled bool   `json:"ReservationDisabled"`
	CreatedTime         string `json:"CreatedTime,omitempty"`
	ExpirationTime      string `json:"ExpirationTime,omitempty"`
}

// AdminStatusCheckResponse is returned by /locks/status.
type AdminStatusCheckResponse struct {
	Components []LockComponentStatus `json:"Components"`
	NotFound   []string              `json:"NotFound,omitempty"`
}
