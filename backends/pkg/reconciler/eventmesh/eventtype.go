package eventmesh

import (
	"knative.dev/eventing/pkg/apis/eventing/v1beta3"
)

// NamespacedName returns the name and namespace of the event type in the format "<namespace>/<name>"
func (et EventType) NamespacedName() string {
	return NamespacedName(et.Namespace, et.Name)
}

// NamespacedType returns the type and namespace of the event type in the format "<namespace>/<type>"
func (et EventType) NamespacedType() string {
	return NamespacedName(et.Namespace, et.Type)
}

// convertEventType converts a Knative Eventing EventType to a simplified representation that is easier to consume by the Backstage plugin.
// see EventType.
func convertEventTypev1beta3(et *v1beta3.EventType) EventType {
	cet := EventType{
		Name:        et.Name,
		Namespace:   et.Namespace,
		Uid:         string(et.UID),
		Description: ToStrPtrOrNil(et.Spec.Description),
		Labels:      et.Labels,
		Annotations: FilterAnnotations(et.Annotations),
		Reference:   ToStrPtrOrNil(NamespacedRefName(et.Spec.Reference)),
		// this field will be populated later on, when we have process the triggers
		ConsumedBy: make([]string, 0),
	}

	if len(et.Spec.Attributes) == 0 {
		return cet
	}

	for _, attr := range et.Spec.Attributes {
		switch attr.Name {
		case "type":
			cet.Type = attr.Value
		case "schema":
			cet.SchemaURL = ToStrPtrOrNil(attr.Value)
		}
	}

	return cet
}
