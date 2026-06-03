package docs

// https://swagger.io/specification/#encoding-object
type EncodingObject struct {
	ContentType   string                       `json:"contentType,omitempty"`
	Headers       map[string]HeaderOrReference `json:"headers,omitempty"`
	Style         string                       `json:"style,omitempty" validate:"omitempty,oneof=matrix label form simple spaceDelimited pipeDelimited deepObject"`
	Explode       bool                         `json:"explode,omitempty"`
	AllowReserved bool                         `json:"allowReserved,omitempty"`
}
