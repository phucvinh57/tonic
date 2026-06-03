package docs

// https://swagger.io/specification/#xml-object
type XmlObject struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	NodeType  string `json:"nodeType,omitempty" validate:"omitempty,oneof=element attribute text cdata none"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}
