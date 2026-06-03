package docs

// https://swagger.io/specification/#media-type-object
type MediaTypeObject struct {
	Schema         *SchemaOrReference            `json:"schema,omitempty" validate:"excluded_with=ItemSchema"`
	ItemSchema     *SchemaOrReference            `json:"itemSchema,omitempty" validate:"excluded_with=Schema"`
	Example        any                           `json:"example,omitempty" validate:"excluded_with=Examples"`
	Examples       map[string]ExampleOrReference `json:"examples,omitempty" validate:"excluded_with=Example"`
	Encoding       map[string]EncodingObject     `json:"encoding,omitempty"`
	ItemEncoding   *EncodingObject               `json:"itemEncoding,omitempty"`
	PrefixEncoding []EncodingObject              `json:"prefixEncoding,omitempty"`
}

type MediaTypeOrReference struct {
	*MediaTypeObject `json:",inline,omitempty"`
	*ReferenceObject `json:",inline,omitempty"`
}
