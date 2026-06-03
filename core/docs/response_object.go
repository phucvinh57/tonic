package docs

// https://swagger.io/specification/#response-object
type ResponseObject struct {
	Description string                          `json:"description,omitempty" validate:"required"`
	Summary     string                          `json:"summary,omitempty"`
	Headers     map[string]HeaderOrReference    `json:"headers,omitempty"`
	Content     map[string]MediaTypeOrReference `json:"content,omitempty"`
	Links       map[string]LinkOrReference      `json:"links,omitempty"`
}

type ResponseOrReference struct {
	*ResponseObject  `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
