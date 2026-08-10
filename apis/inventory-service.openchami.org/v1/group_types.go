package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

type Group struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       GroupSpec        `json:"spec" validate:"required"`
	Status     GroupStatus      `json:"status,omitempty"`
}

type GroupSpec struct {
	Description    string `json:"description,omitempty" validate:"max=200"`
	Label          string `json:"label" yaml:"label"`
	ExclusiveGroup string `json:"exclusiveGroup,omitempty" yaml:"exclusiveGroup,omitempty"`

	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Members Members  `json:"members" yaml:"members"`
}

type GroupStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

func (r *Group) Validate(ctx context.Context) error {

	return nil
}

func (r *Group) GetKind() string {
	return "Group"
}

func (r *Group) GetName() string {
	return r.Metadata.Name
}

func (r *Group) GetUID() string {
	return r.Metadata.UID
}

func (r *Group) IsHub() {}

type Members struct {
	IDs []string `json:"ids" yaml:"ids"`
}
type Membership struct {
	ID string `json:"id" yaml:"id"`

	GroupLabels []string `json:"groupLabels" yaml:"groupLabels"`

	PartitionName string `json:"partitionName" yaml:"partitionName"`
}
