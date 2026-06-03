package docs

// https://swagger.io/specification/#parameter-object
type ParameterObject struct {
	Name        string `json:"name,omitempty" validate:"required"`
	In          string `json:"in,omitempty" validate:"required,oneof=query querystring header path cookie"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
	AllowEmpty  bool   `json:"allowEmptyValue,omitempty"`

	Style         string                        `json:"style,omitempty" validate:"omitempty,oneof=matrix label form simple spaceDelimited pipeDelimited deepObject"`
	Explode       bool                          `json:"explode,omitempty"`
	AllowReserved bool                          `json:"allowReserved,omitempty"`
	Schema        *SchemaOrReference            `json:"schema,omitempty" validate:"excluded_with=Content"`
	Example       any                           `json:"example,omitempty" validate:"excluded_with=Examples Content"`
	Examples      map[string]ExampleOrReference `json:"examples,omitempty" validate:"excluded_with=Example Content"`

	Content map[string]MediaTypeOrReference `json:"content,omitempty" validate:"excluded_with=Schema Example Examples"`
}

type ParameterOrReference struct {
	*ParameterObject `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
