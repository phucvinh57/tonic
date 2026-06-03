package utils_test

import (
	"testing"

	"github.com/TickLabVN/tonic/core/utils"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestMergeStructs_OverrideEmptyField(t *testing.T) {
	assert := assert.New(t)
	a := TestStruct{Field1: "Hello"}
	b := TestStruct{Field1: "World", Field2: 0}

	result := utils.MergeStructs(a, b)

	assert.Equal("World", result.Field1)
	assert.Equal(0, result.Field2)
}

func TestMergeStructs_OverrideNonEmptyField(t *testing.T) {
	assert := assert.New(t)
	a := TestStruct{Field1: "Hello", Field2: 42}
	b := TestStruct{Field1: "World", Field2: 0}

	result := utils.MergeStructs(a, b)

	assert.Equal("World", result.Field1)
	assert.Equal(0, result.Field2)
}

func TestMergeStructs_EmptyStruct(t *testing.T) {
	assert := assert.New(t)
	a := TestStruct{}
	b := TestStruct{Field1: "World", Field2: 0}

	result := utils.MergeStructs(a, b)

	assert.Equal("World", result.Field1)
	assert.Equal(0, result.Field2)
}

func TestMergeStructs_EmptyFields(t *testing.T) {
	assert := assert.New(t)
	a := TestStruct{Field1: "World", Field2: 0}
	b := TestStruct{}

	result := utils.MergeStructs(a, b)
	assert.Equal("World", result.Field1)
	assert.Equal(0, result.Field2)
}

func TestMergeStructs_MergeMultipleStruct(t *testing.T) {
	assert := assert.New(t)
	a := TestStruct{Field1: "Hello"}
	b := TestStruct{Field1: "World", Field2: 42}
	c := TestStruct{Field2: 100}

	result := utils.MergeStructs(a, b, c)

	assert.Equal("World", result.Field1)
	assert.Equal(100, result.Field2)
}
