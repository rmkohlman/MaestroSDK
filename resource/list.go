package resource

import (
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/yaml.v3"
)

// ResourceList represents a kubectl-style List wrapper for multiple resources.
// This is a meta-resource that is not stored in the database — it is produced
// by "dvm get all -o yaml/json" and consumed by "dvm apply -f -" for round-trip
// export/import.
type ResourceList struct {
	APIVersion string         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string         `json:"kind"       yaml:"kind"`
	Metadata   map[string]any `json:"metadata"   yaml:"metadata"`
	Items      []any          `json:"items"      yaml:"items"`
}

// DependencyOrder defines the order in which resources should be applied
// when processing a List document. Dependencies come before dependents.
var DependencyOrder = []string{
	"Ecosystem",
	"Domain",
	"GitRepo",
	"App",
	"Registry",
	"Workspace",
	"Credential",
	"NvimPlugin",
	"NvimTheme",
	"NvimPackage",
	"TerminalPrompt",
	"TerminalPackage",
}

// NewResourceList creates an empty ResourceList with the correct envelope fields.
func NewResourceList() *ResourceList {
	return &ResourceList{
		APIVersion: "devopsmaestro.io/v1",
		Kind:       "List",
		Metadata:   map[string]any{},
		Items:      []any{},
	}
}

// BuildList builds a ResourceList from the given resources by serializing each
// resource via its registered handler's ToYAML() method and unmarshaling the
// result into map[string]any for inclusion in the Items array.
//
// The resources are expected to be passed in the desired order (typically
// DependencyOrder). BuildList preserves the input order; the caller
// (e.g., cmd/get_all.go) is responsible for collecting resources in
// dependency order.
//
// If ToYAML() fails for a resource, it is skipped with a warning log.
func BuildList(ctx Context, resources []Resource) (*ResourceList, error) {
	list := NewResourceList()

	for _, res := range resources {
		yamlBytes, err := ToYAML(res)
		if err != nil {
			slog.Warn("failed to serialize resource, skipping",
				"kind", res.GetKind(),
				"name", res.GetName(),
				"error", err)
			continue
		}

		var item map[string]any
		if err := yaml.Unmarshal(yamlBytes, &item); err != nil {
			slog.Warn("failed to unmarshal resource YAML, skipping",
				"kind", res.GetKind(),
				"name", res.GetName(),
				"error", err)
			continue
		}

		list.Items = append(list.Items, item)
	}

	return list, nil
}

// detectName extracts the metadata.name field from YAML data without full parsing.
func detectName(data []byte) string {
	var header struct {
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
	}
	if err := yaml.Unmarshal(data, &header); err != nil || header.Metadata.Name == "" {
		return "unknown"
	}
	return header.Metadata.Name
}

// ApplyList parses a List YAML document and applies each item via the
// registered handler for its kind, continuing on error (kubectl precedent).
// Returns all successfully applied resources and a combined error summary.
func ApplyList(ctx Context, data []byte) ([]Resource, error) {
	var list ResourceList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to parse List: %w", err)
	}

	if list.Kind != "List" {
		return nil, fmt.Errorf("expected kind 'List', got '%s'", list.Kind)
	}

	var applied []Resource
	var errs []error

	for i, item := range list.Items {
		// Marshal item back to YAML for handler consumption
		itemBytes, err := yaml.Marshal(item)
		if err != nil {
			errs = append(errs, fmt.Errorf("item %d: failed to marshal: %w", i, err))
			continue
		}

		kind, err := DetectKind(itemBytes)
		if err != nil {
			errs = append(errs, fmt.Errorf("item %d: %w", i, err))
			continue
		}

		name := detectName(itemBytes)

		handler, err := MustGetHandler(kind)
		if err != nil {
			errs = append(errs, fmt.Errorf("item %d (%s %q): %w", i, kind, name, err))
			continue
		}

		res, err := handler.Apply(ctx, itemBytes)
		if err != nil {
			errs = append(errs, fmt.Errorf("item %d (%s %q): %w", i, kind, name, err))
			continue
		}

		applied = append(applied, res)
	}

	if len(errs) > 0 {
		return applied, fmt.Errorf("%d of %d items failed to apply:\n%w",
			len(errs), len(list.Items), errors.Join(errs...))
	}

	return applied, nil
}
