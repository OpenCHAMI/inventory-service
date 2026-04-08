// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT
package v1

type Membership struct {
	ID            string   `json:"id"`
	GroupLabels   []string `json:"groupLabels"`
	PartitionName string   `json:"partitionName"`
}
