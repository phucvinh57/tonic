package docs

// / https://swagger.io/specification/#openapi-object
type OpenApi struct {
	OpenAPI           string                         `json:"openapi" validate:"required"`
	Self              string                         `json:"$self,omitempty"`
	Info              InfoObject                     `json:"info" validate:"required"`
	JsonSchemaDialect string                         `json:"jsonSchemaDialect,omitempty"`
	Servers           []ServerObject                 `json:"servers,omitempty"`
	Paths             Paths                          `json:"paths,omitempty"`
	Webhooks          map[string]PathItemOrReference `json:"webhooks,omitempty"`
	Components        ComponentsObject               `json:"components,omitempty"`
	Security          []SecurityRequirement          `json:"security,omitempty"`
	Tags              []Tag                          `json:"tags,omitempty"`
	ExternalDocs      *ExternalDocumentationObject   `json:"externalDocs,omitempty"`
}
