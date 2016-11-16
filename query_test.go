package aedstorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEntityName(t *testing.T) {
	tm := testModel{}
	entityName, err := getEntityName(tm)
	assert.NoError(t, err)
	assert.Equal(t, "testModel", entityName)
}

func TestGetEntityNameNil(t *testing.T) {
	entityName, err := getEntityName(nil)
	assert.Equal(t, ErrNilModel, err)
	assert.Equal(t, "", entityName)
}

func TestGetEntityNameNonStruct(t *testing.T) {
	var b bool
	entityName, err := getEntityName(&b)
	assert.Equal(t, ErrModelInvalid, err)
	assert.Equal(t, "", entityName)
}

func TestGetEntityNamePtr(t *testing.T) {
	tm := &testModel{}
	entityName, err := getEntityName(tm)
	assert.NoError(t, err)
	assert.Equal(t, "testModel", entityName)
}

func TestModelWithNoEntity(t *testing.T) {
	m := modelWithNoEntity{}
	entityName, err := getEntityName(m)
	assert.NoError(t, err)
	assert.Equal(t, "modelWithNoEntity", entityName)
}

func TestModelWithNoEntityPtr(t *testing.T) {
	m := &modelWithNoEntity{}
	entityName, err := getEntityName(m)
	assert.NoError(t, err)
	assert.Equal(t, "modelWithNoEntity", entityName)
}

func TestNewQuery(t *testing.T) {
	assert.NotPanics(t, func() {
		NewQuery(&modelWithNoEntity{})
	})
}

func TestNewQueryInvalid(t *testing.T) {
	assert.Panics(t, func() {
		NewQuery(nil)
	})
}
