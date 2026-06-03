package docs

// https://swagger.io/specification/#request-body-object
type RequestBodyObject struct {
	Description string                          `json:"description,omitempty"`
	Content     map[string]MediaTypeOrReference `json:"content,omitempty" validate:"required"`
	Required    bool                            `json:"required,omitempty"`
}

type RequestBodyOrReference struct {
	*RequestBodyObject `json:",inline,omitempty"`
	*ReferenceObject   `json:",inline,omitempty"`
}
