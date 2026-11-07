package schema_test

import (
	"testing"

	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/stretchr/testify/assert"
)

func TestCustomTagNames(t *testing.T) {
	assert := assert.New(t)
	type CertificateKeyPair struct {
		Crt []byte `json:"crt" swtype:"string" fmt:"base64" eg:"U3dhZ2dlciByb2Nrcw=="`
	}

	spec := core.Init(core.WithTagConfig(docs.TagConfig{
		SwaggerType: "swtype",
		Format:      "fmt",
		Example:     "eg",
	}))
	schema, err := AssertParse(assert, spec, CertificateKeyPair{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"crt": { "type": "string", "format": "base64", "examples": ["U3dhZ2dlciByb2Nrcw=="] }
		}
	}`, schema)
}

func TestCustomTagNamesIgnoreDefaults(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Crt []byte `json:"crt" swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
	}

	spec := core.Init(core.WithTagConfig(docs.TagConfig{
		SwaggerType: "swtype",
		Format:      "fmt",
		Example:     "eg",
	}))
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"crt": { "type": "array", "items": { "type": "integer", "format": "uint8" } }
		}
	}`, schema)
}

func TestPartialTagOverride(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Crt []byte `json:"crt" swaggertype:"string" fmt:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
	}

	spec := core.Init(core.WithTagConfig(docs.TagConfig{Format: "fmt"}))
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"crt": { "type": "string", "format": "base64", "examples": ["U3dhZ2dlciByb2Nrcw=="] }
		}
	}`, schema)
}

func TestDefaultTagNamesUnchanged(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Crt []byte `json:"crt" swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"crt": { "type": "string", "format": "base64", "examples": ["U3dhZ2dlciByb2Nrcw=="] }
		}
	}`, schema)
}
