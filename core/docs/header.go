package docs

// https://swagger.io/specification/#header-object
type HeaderObject struct {
	Description string `json:"description,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`

	Style         string                        `json:"style,omitempty" validate:"omitempty,oneof=simple"`
	Explode       bool                          `json:"explode,omitempty"`
	AllowReserved bool                          `json:"allowReserved,omitempty"`
	Schema        *SchemaOrReference            `json:"schema,omitempty" validate:"excluded_with=Content"`
	Example       any                           `json:"example,omitempty" validate:"excluded_with=Examples Content"`
	Examples      map[string]ExampleOrReference `json:"examples,omitempty" validate:"excluded_with=Example Content"`

	Content map[string]MediaTypeOrReference `json:"content,omitempty" validate:"excluded_with=Schema Example Examples"`
}

type HeaderOrReference struct {
	*HeaderObject    `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
