package v1

import "context"

var ObjValidator Validator

// Validator defines validation operations for all resource types.
type Validator interface {
	ValidateComponent(ctx context.Context, r *Component) error
	ValidateComponentEndpoint(ctx context.Context, r *ComponentEndpoint) error
	ValidateEthernetInterface(ctx context.Context, r *EthernetInterface) error
	ValidateGroup(ctx context.Context, r *Group) error
	ValidateRedfishEndpoint(ctx context.Context, r *RedfishEndpoint) error
	ValidateServiceEndpoint(ctx context.Context, r *ServiceEndpoint) error
}
