package docs

// https://swagger.io/specification/#server-object
type ServerObject struct {
	Name        string                          `json:"name,omitempty"`
	Url         string                          `json:"url" validate:"required"`
	Description string                          `json:"description,omitempty"`
	Variables   map[string]ServerVariableObject `json:"variables,omitempty"`
}
