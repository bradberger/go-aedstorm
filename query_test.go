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
