package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/inventory-service/schemas"
)

type Component struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	ID         string           `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       ComponentSpec    `json:"spec" yaml:"spec" validate:"required"`
	Status     ComponentStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

type ComponentSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID" yaml:"ID"`

	Type     string      `json:"Type" yaml:"Type"`
	State    string      `json:"State,omitempty" yaml:"State,omitempty"`
	Flag     string      `json:"Flag,omitempty" yaml:"Flag,omitempty"`
	Enabled  *bool       `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`
	SwStatus string      `json:"SoftwareStatus,omitempty" yaml:"SoftwareStatus,omitempty"`
	Role     string      `json:"Role,omitempty" yaml:"Role,omitempty"`
	SubRole  string      `json:"SubRole,omitempty" yaml:"SubRole,omitempty"`
	NID      json.Number `json:"NID,omitempty" yaml:"NID,omitempty"`
	Subtype  string      `json:"Subtype,omitempty" yaml:"Subtype,omitempty"`
	NetType  string      `json:"NetType,omitempty" yaml:"NetType,omitempty"`
	Arch     string      `json:"Arch,omitempty" yaml:"Arch,omitempty"`
	Class    string      `json:"Class,omitempty" yaml:"Class,omitempty"`

	ReservationDisabled bool `json:"ReservationDisabled,omitempty" yaml:"ReservationDisabled,omitempty"`
	Locked              bool `json:"Locked,omitempty" yaml:"Locked,omitempty"`
}

type ComponentStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
}

func (r *Component) Validate(ctx context.Context) error {
	var schema jsonschema.Schema
	if err := json.Unmarshal(schemas.ComponentsSchema, &schema); err != nil {
		return fmt.Errorf("loading component schema: %w", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolving component schema: %w", err)
	}

	resourceJSON, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshaling resource for validation: %w", err)
	}

	var instance any
	if err := json.Unmarshal(resourceJSON, &instance); err != nil {
		return fmt.Errorf("unmarshaling resource for validation: %w", err)
	}

	return resolved.Validate(instance)
}

func (r *Component) GetKind() string {
	return "Component"
}

func (r *Component) GetName() string {
	return r.Metadata.Name
}

func (r *Component) GetUID() string {
	return r.Metadata.UID
}

func (r *Component) IsHub() {}
