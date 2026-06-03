package docs

type Discriminator struct {
	PropertyName   string            `json:"propertyName" validate:"required"`
	Mapping        map[string]string `json:"mapping,omitempty"`
	DefaultMapping string            `json:"defaultMapping,omitempty"`
}
