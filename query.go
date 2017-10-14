package aedstorm

import (
	"errors"
	"reflect"
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

var (
	nextResultData  interface{}
	nextResultError error
	nextResultKeys  []*datastore.Key
	mu              sync.Mutex
)

// SetMockQueryResult sets the query GetAll() result to be the values given
func SetMockQueryResult(data interface{}, err error, keys []*datastore.Key) {
	mu.Lock()
	defer mu.Unlock()
	nextResultData = data
	nextResultError = err
	nextResultKeys = keys
}

func clearNextResult() {
	mu.Lock()
	defer mu.Unlock()
	nextResultData = nil
	nextResultError = nil
	nextResultKeys = nil
}

func hasNextResult() bool {
	mu.Lock()
	defer mu.Unlock()
	return nextResultData != nil || nextResultError != nil || nextResultKeys != nil
}

// Query is a struct which implements a subset of the "datastore.Query" interface and is mockable
type Query struct {
	entity string
	dq     *datastore.Query
}

func (q *Query) Limit(num int) *Query {
	if q.dq == nil {
		q.dq = datastore.NewQuery(q.entity)
	}
	q.dq = q.dq.Limit(num)
	return q
}

// Filter implements the "datastore.Query".Filter interface
func (q *Query) Filter(filterStr string, value interface{}) *Query {
	if q.dq == nil {
		q.dq = datastore.NewQuery(q.entity)
	}
	q.dq = q.dq.Filter(filterStr, value)
	return q
}

// Order returns a derivative query with a field-based sort order. Orders are
// applied in the order they are added. The default order is ascending; to sort
// in descending order prefix the fieldName with a minus sign (-).
func (q *Query) Order(fieldName string) *Query {
	if q.dq == nil {
		q.dq = datastore.NewQuery(q.entity)
	}
	q.dq = q.dq.Order(fieldName)
	return q
}

// Count matches the "datastore.Query".Count interface
func (q *Query) Count(ctx context.Context) (int, error) {
	if q.dq == nil {
		q.dq = datastore.NewQuery(q.entity)
	}
	return q.dq.Count(ctx)
}

// GetAll matches the "datastore.Query".GetAll interface
func (q *Query) GetAll(ctx context.Context, out interface{}) ([]*datastore.Key, error) {

	// For purposes of mocking, this allows a one-time return value to be preset in advance
	if hasNextResult() {
		defer clearNextResult()
		if err := Copy(nextResultData, out); err != nil {
			return nil, err
		}
		return nextResultKeys, nextResultError
	}

	if q.dq == nil {
		q.dq = datastore.NewQuery(q.entity)
	}
	return q.dq.GetAll(ctx, out)
}

// NewQuery returns a new query based off the type of m. If m implements the EntityName interface, it uses
// that for an entity name, otherwise it uses the name of the struct itself
func NewQuery(m interface{}) *Query {
	entityKind, err := getEntityName(m)
	if err != nil {
		panic(err)
	}
	return &Query{entity: entityKind}
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

// Copy copies one interface into the other doing type checking to make sure
// it's safe. If it cannot be copied, an error is returned.
func Copy(srcVal interface{}, dstVal interface{}) error {
	curEl := reflect.ValueOf(srcVal)
	if curEl.Kind() == reflect.Ptr {
		curEl = curEl.Elem()
	}
	dstEl := reflect.ValueOf(dstVal)
	if dstEl.Kind() == reflect.Ptr {
		dstEl = dstEl.Elem()
	}
	if !dstEl.CanSet() {
		return errors.New("Destination value type is invalid")
	}
	if !curEl.Type().AssignableTo(dstEl.Type()) {
		return errors.New("Cannot assign value to destination")
	}
	dstEl.Set(curEl)
	return nil
}
