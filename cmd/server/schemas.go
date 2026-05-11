package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OpenCHAMI/inventory-service/schemas"
)

// schemaOverrides maps each schema filename to a pointer to the package-level
// []byte variable it should overwrite when a custom file is found.
var schemaOverrides = map[string]*[]byte{
	"components_schema.json":        &schemas.ComponentsSchema,
	"ethernet_interface_schema.json": &schemas.EthernetInterfaceSchema,
	"group_schema.json":             &schemas.GroupSchema,
	"hardware_schema.json":          &schemas.HardwareSchema,
	"redfish_endpoint_schema.json":  &schemas.RedfishEndpointSchema,
	"service_endpoint_schema.json":  &schemas.ServiceEndpointSchema,
}

// LoadCustomSchemas reads any schema JSON files found in dir and replaces the
// corresponding embedded defaults in the schemas package.  Files whose names
// are not recognised are silently ignored.
func LoadCustomSchemas(dir string) error {
	if dir == "" {
		return nil
	}

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("schema directory %q not accessible: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("schema path %q is not a directory", dir)
	}

	for filename, target := range schemaOverrides {
		path := filepath.Join(dir, filename)
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to read custom schema %q: %w", path, err)
		}
		*target = data
	}

	return nil
}
