package schema_test

import (
	"testing"

	"github.com/TickLabVN/tonic/core"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerTypeBytesToString(t *testing.T) {
	assert := assert.New(t)
	type CertificateKeyPair struct {
		Crt []byte `json:"crt" swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
		Key []byte `json:"key" swaggertype:"string" format:"base64" example:"U3dhZ2dlciByb2Nrcw=="`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, CertificateKeyPair{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"crt": { "type": "string", "format": "base64", "examples": ["U3dhZ2dlciByb2Nrcw=="] },
			"key": { "type": "string", "format": "base64", "examples": ["U3dhZ2dlciByb2Nrcw=="] }
		}
	}`, schema)
}

func TestExampleIntegerCoercion(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Age int `json:"age" example:"123"`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"age": { "type": "integer", "format": "int32", "examples": [123] }
		}
	}`, schema)
}

func TestExampleBooleanCoercion(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Active bool `json:"active" example:"true"`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"active": { "type": "boolean", "examples": [true] }
		}
	}`, schema)
}

func TestFormatTagOnPlainString(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Email string `json:"email" format:"email"`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"email": { "type": "string", "format": "email" }
		}
	}`, schema)
}

func TestSwaggerTypeArrayOfNumber(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Scores string `json:"scores" swaggertype:"array,number"`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"scores": { "type": "array", "items": { "type": "number" } }
		}
	}`, schema)
}

func TestNoOverrideTagsIsUnchanged(t *testing.T) {
	assert := assert.New(t)
	type Payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	spec := core.Init()
	schema, err := AssertParse(assert, spec, Payload{})
	assert.Nil(err)
	assert.JSONEq(`{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "integer", "format": "int32" }
		}
	}`, schema)
}
