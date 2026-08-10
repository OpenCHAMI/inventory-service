package v1

import (
	"context"
	"encoding/json"

	"github.com/openchami/fabrica/pkg/fabrica"
)

type ServiceEndpoint struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   fabrica.Metadata      `json:"metadata"`
	Spec       ServiceEndpointSpec   `json:"spec" validate:"required"`
	Status     ServiceEndpointStatus `json:"status,omitempty"`
}

type ServiceEndpointSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	ServiceDescription

	RfEndpointFQDN string `json:"RedfishEndpointFQDN" yaml:"RedfishEndpointFQDN"`
	URL            string `json:"RedfishURL" yaml:"RedfishURL"`

	ServiceInfo json.RawMessage `json:"ServiceInfo,omitempty" yaml:"ServiceInfo,omitempty"`
}

type ServiceEndpointStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

func (r *ServiceEndpoint) Validate(ctx context.Context) error {

	return nil
}

func (r *ServiceEndpoint) GetKind() string {
	return "ServiceEndpoint"
}

func (r *ServiceEndpoint) GetName() string {
	return r.Metadata.Name
}

func (r *ServiceEndpoint) GetUID() string {
	return r.Metadata.UID
}

func (r *ServiceEndpoint) IsHub() {}

type ServiceDescription struct {
	RfEndpointID string `json:"RedfishEndpointID" yaml:"RedfishEndpointID"`
	RedfishType  string `json:"RedfishType" yaml:"RedfishType"`

	RedfishSubtype string `json:"RedfishSubtype,omitempty" yaml:"RedfishSubtype,omitempty"`

	UUID    string `json:"UUID" yaml:"UUID"`
	OdataID string `json:"OdataID" yaml:"OdataID"`
}
