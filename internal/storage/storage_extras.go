// Copyright © 2025-2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package storage

import (
	"context"
	"fmt"

	v1 "github.com/OpenCHAMI/smd2/apis/smd2.openchami.org/v1"
)

func LoadComponentByID(ctx context.Context, id string) (*v1.Component, error) {
	// todo Change to not have to read all components
	components, err := LoadAllComponents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load components: %w", err)
	}
	var component *v1.Component
	for _, c := range components {
		if c.Spec.ID == id {
			component = c
			break
		}
	}
	return component, nil
}
