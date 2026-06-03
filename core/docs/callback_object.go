package docs

// Field Pattern	Type	Description
// {expression}	Path Item Object | Reference Object	A Path Item Object, or a reference to one, used to define a callback request and expected responses. A complete example is available.
type CallbackObject map[string]PathItemOrReference

type CallbackOrReference struct {
	*CallbackObject  `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
