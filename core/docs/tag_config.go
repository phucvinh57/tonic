package docs

type TagConfig struct {
	SwaggerType string // default: "swaggertype"
	Format      string // default: "format"
	Example     string // default: "example"
}

func DefaultTagConfig() TagConfig {
	return TagConfig{
		SwaggerType: swaggerTypeTag,
		Format:      formatTag,
		Example:     exampleTag,
	}
}

func (tc TagConfig) withDefaults() TagConfig {
	d := DefaultTagConfig()
	if tc.SwaggerType == "" {
		tc.SwaggerType = d.SwaggerType
	}
	if tc.Format == "" {
		tc.Format = d.Format
	}
	if tc.Example == "" {
		tc.Example = d.Example
	}
	return tc
}
