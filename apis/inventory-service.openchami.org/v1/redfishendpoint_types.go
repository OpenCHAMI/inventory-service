package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

type RedfishEndpoint struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   fabrica.Metadata      `json:"metadata"`
	Spec       RedfishEndpointSpec   `json:"spec" validate:"required"`
	Status     RedfishEndpointStatus `json:"status,omitempty"`
}

type RedfishEndpointSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID" yaml:"ID"`

	Type string `json:"Type" yaml:"Type"`
	Name string `json:"Name,omitempty" yaml:"Name,omitempty"`

	Hostname string `json:"Hostname" yaml:"Hostname"`
	Domain   string `json:"Domain" yaml:"Domain"`
	FQDN     string `json:"FQDN" yaml:"FQDN"`
	Enabled  bool   `json:"Enabled" yaml:"Enabled"`

	UUID     string `json:"UUID,omitempty" yaml:"UUID,omitempty"`
	User     string `json:"User" yaml:"User"`
	Password string `json:"Password" yaml:"Password"`

	UseSSDP     bool `json:"UseSSDP,omitempty" yaml:"UseSSDP,omitempty"`
	MACRequired bool `json:"MACRequired,omitempty" yaml:"MACRequired,omitempty"`

	MACAddr            string `json:"MACAddr,omitempty" yaml:"MACAddr,omitempty"`
	IPAddress          string `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
	RedsicoverOnUpdate bool   `json:"RediscoverOnUpdate" yaml:"RediscoverOnUpdate"`

	TemplateID string `json:"TemplateID,omitempty" yaml:"TemplateID,omitempty"`

	DiscoveryInfo DiscoveryInfo `json:"DiscoveryInfo" yaml:"DiscoveryInfo"`
}

type RedfishEndpointStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

func (r *RedfishEndpoint) Validate(ctx context.Context) error {

	return nil
}

func (r *RedfishEndpoint) GetKind() string {
	return "RedfishEndpoint"
}

func (r *RedfishEndpoint) GetName() string {
	return r.Metadata.Name
}

func (r *RedfishEndpoint) GetUID() string {
	return r.Metadata.UID
}

func (r *RedfishEndpoint) IsHub() {}

type DiscoveryInfo struct {
	LastAttempt string `json:"LastDiscoveryAttempt,omitempty" yaml:"LastDiscoveryAttempt,omitempty"`

	LastStatus     string `json:"LastDiscoveryStatus" yaml:"LastDiscoveryStatus"`
	RedfishVersion string `json:"RedfishVersion,omitempty" yaml:"RedfishVersion,omitempty"`
}
