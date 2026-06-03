package schema_test

import (
	"testing"

	"github.com/TickLabVN/tonic/core"
	"github.com/stretchr/testify/assert"
)

// Comparision validate tags: https://github.com/go-playground/validator?tab=readme-ov-file#comparisons
func TestValidate_Comparision(t *testing.T) {
	assert := assert.New(t)

	t.Run("Eq", func(t *testing.T) {
		type TestEq struct {
			Age           int      `json:"age" validate:"eq=21 24"`
			Name          string   `json:"name" validate:"eq=John Alice"`
			IsAdmin       bool     `json:"isAdmin" validate:"eq=true"`
			RecoverEmails []string `json:"recoverEmails" validate:"min=1,max=6"`
			Score         float64  `json:"score" validate:"eq=3.14"`
		}
		schema, err := AssertParse(assert, core.Init(), TestEq{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"enum": [21, 24]
				},
				"name": {
					"type": "string",
					"enum": ["John", "Alice"]
				},
				"isAdmin": {
					"type": "boolean",
					"enum": [true]
				},
				"score": {
					"type": "number",
					"format": "float64",
					"enum": [3.14]
				},
				"recoverEmails": {
					"type": "array",
					"items": {"type": "string"},
					"minItems": 1,
					"maxItems": 6
				}
			}
		}`, schema)
	})
	t.Run("EqIgnoreCase", func(t *testing.T) {
		type TestEqIgnoreCase struct {
			Name string `json:"name" validate:"eq_ignore_case=John Doe"`
		}

		schema, err := AssertParse(assert, core.Init(), TestEqIgnoreCase{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"name": {
					"type": "string",
					"pattern": "^(?i)(John|Doe)$"
				}
			}
		}`, schema)
	})

	t.Run("Gt", func(t *testing.T) {
		type TestGt struct {
			Age            int      `json:"age" validate:"gt=21"`
			Name           string   `json:"name" validate:"gt=5"`
			RecoverEmails  []string `json:"recoverEmails" validate:"gt=1.2"`
			InterestedList []string `json:"interestedList" validate:"gt=2"`
		}
		schema, err := AssertParse(assert, core.Init(), TestGt{})
		assert.Nil(err)

		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"minimum": 21,
					"exclusiveMinimum": true
				},
				"name": {
					"type": "string",
					"minLength": 6
				},
				"recoverEmails": {
					"type": "array",
					"items": { "type": "string" },
					"minItems": 2
				},
				"interestedList": {
					"type": "array",
					"items": { "type": "string" },
					"minItems": 3
				}
			}
		}`, schema)
	})

	t.Run("Gte", func(t *testing.T) {
		type TestGte struct {
			Age           int      `json:"age" validate:"gte=21"`
			Name          string   `json:"name" validate:"gte=5"`
			RecoverEmails []string `json:"recoverEmails" validate:"gte=1"`
		}
		b, err := AssertParse(assert, core.Init(), TestGte{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"minimum": 21
				},
				"name": {
					"type": "string",
					"minLength": 5
				},
				"recoverEmails": {
					"type": "array",
					"items": { "type": "string" },
					"minItems": 1
				}
			}
		}`, string(b))
	})

	t.Run("Lt", func(t *testing.T) {
		type TestLt struct {
			Age            int      `json:"age" validate:"lt=21"`
			Name           string   `json:"name" validate:"lt=5"`
			RecoverEmails  []string `json:"recoverEmails" validate:"lt=1.2"`
			InterestedList []string `json:"interestedList" validate:"lt=2"`
		}
		b, err := AssertParse(assert, core.Init(), TestLt{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"maximum": 21,
					"exclusiveMaximum": true
				},
				"name": {
					"type": "string",
					"maxLength": 4
				},
				"recoverEmails": {
					"type": "array",
					"items": { "type": "string" },
					"maxItems": 1
				},
				"interestedList": {
					"type": "array",
					"items": { "type": "string" },
					"maxItems": 1
				}
			}
		}`, string(b))
	})

	t.Run("Lte", func(t *testing.T) {
		type TestLte struct {
			Age            int      `json:"age" validate:"lte=21"`
			Name           string   `json:"name" validate:"lte=5"`
			RecoverEmails  []string `json:"recoverEmails" validate:"lte=1.2"`
			InterestedList []string `json:"interestedList" validate:"lte=2"`
		}
		b, err := AssertParse(assert, core.Init(), TestLte{})
		assert.Nil(err)

		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"maximum": 21
				},
				"name": {
					"type": "string",
					"maxLength": 5
				},
				"recoverEmails": {
					"type": "array",
					"items": { "type": "string" },
					"maxItems": 1
				},
				"interestedList": {
					"type": "array",
					"items": { "type": "string" },
					"maxItems": 2
				}
			}
		}`, string(b))
	})

	t.Run("Ne", func(t *testing.T) {
		type TestEq struct {
			Age           int      `json:"age" validate:"ne=21"`
			Name          string   `json:"name" validate:"ne=John Doe"`
			IsAdmin       bool     `json:"isAdmin" validate:"ne=true"`
			RecoverEmails []string `json:"recoverEmails" validate:"ne=1"`
			Score         float64  `json:"score" validate:"ne=3.14"`
		}
		b, err := AssertParse(assert, core.Init(), TestEq{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"age": {
					"type": "integer",
					"format": "int32",
					"not": { "enum": [21] }
				},
				"name": {
					"type": "string",
					"not": { "enum": ["John", "Doe"] }
				},
				"isAdmin": {
					"type": "boolean",
					"not": { "enum": [true] }
				},
				"score": {
					"type": "number",
					"format": "float64",
					"not": { "enum": [3.14] }
				},
				"recoverEmails": {
					"type": "array",
					"items": {
						"type": "string"
					},
					"not": {
						"minItems": 1,
						"maxItems": 1
					}
				}
			}
		}`, string(b))
	})

	t.Run("NeIgnoreCase", func(t *testing.T) {
		type TestNeIgnoreCase struct {
			Name string `json:"name" validate:"ne_ignore_case=John Doe"`
		}
		b, err := AssertParse(assert, core.Init(), TestNeIgnoreCase{})
		assert.Nil(err)
		assert.JSONEq(`{
			"type": "object",
			"properties": {
				"name": {
					"type": "string",
					"not": { "pattern": "^(?i)(John|Doe)$" }
				}
			}
		}`, string(b))
	})
}
