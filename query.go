package aedstorm

import (
	"reflect"

	"google.golang.org/appengine/datastore"
)

// NewQuery returns a new query based off the type of m. If m implements the EntityName interface, it uses
// that for an entity name, otherwise it uses the name of the struct itself
func NewQuery(m interface{}) *datastore.Query {
	entityKind, err := getEntityName(m)
	if err != nil {
		panic(err)
	}
	return datastore.NewQuery(entityKind)
}

// getEntityName returns the name of a struct type based on the EntityName interface value
func getEntityName(m interface{}) (string, error) {
	if m == nil {
		return "", ErrNilModel
	}
	t := reflect.ValueOf(m).Type()
	// If a pointer, get it's underlying type now.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Make sure it's a struct
	if t.Kind() != reflect.Struct {
		return "", ErrModelInvalid
	}
	// Try to get the entity name from a Entity() method. If not, fall back to the name of the type itself
	entityKind := t.Name()
	if obj, ok := m.(EntityName); ok {
		if n := obj.Entity(); n != "" {
			entityKind = n
		}
	}
	return entityKind, nil
}
