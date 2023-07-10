package aedstorm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"

	gocache "github.com/rubanbydesign/gocache/cache"
	"github.com/rubanbydesign/rest/cache"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// Standardized error messages
var (
	ErrNoModel      = errors.New("Model is not loaded")
	ErrNoID         = errors.New("Model has no ID")
	ErrNoContext    = errors.New("No net/context was loaded")
	ErrNilModel     = errors.New("Model is nil")
	ErrModelInvalid = errors.New("Model must be a struct pointer")
)

const (
	// TagName is the tag name where we look for custom tag values, like "id"
	TagName = "datastore"
)

// DataModel is a ORM styled structure for saving and loading entities
type DataModel struct {
	model       Model
	ctx         context.Context
	verified    bool
	entity      string
	idFieldName string
	uid         string
	sync.Mutex
}

// NewModel returns an initialized data model
func NewModel(m Model) *DataModel {
	t := reflect.ValueOf(m).Type()
	if t.Kind() != reflect.Ptr {
		panic(ErrModelInvalid)
	}
	if t.Elem().Kind() != reflect.Struct {
		panic(ErrModelInvalid)
	}
	return &DataModel{model: m}
}

// verify checks that the DataModel has an ID field, context, and non-nil model.
// It should be called before any function that needs one of those things. The
// results are stored in memory for a bit better performance.
func (dm *DataModel) verify() (err error) {
	if dm.verified {
		return
	}
	if dm.ctx == nil {
		return ErrNoContext
	}
	if dm.model == nil {
		return ErrNoModel
	}
	dm.Lock()
	defer dm.Unlock()

	// See if implements the EntityID interface, if not then try to guess the field
	_, ok := dm.model.(EntityID)

	// Try to get the ID field, but only fail if the model doesn't implement EntityID
	if dm.idFieldName, err = dm.getIDField(); !ok && err != nil {
		return err
	}

	dm.verified = true
	return nil
}

// getIDField gets the name of the struct field which serves as an ID for the given model.
func (dm *DataModel) getIDField() (string, error) {

	t := reflect.ValueOf(dm.model).Type().Elem()

	// Iterate over all available fields and read the tag value or return if the field name is ID
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(TagName)
		if field.Name == "ID" || tag == "id" {
			return field.Name, nil
		}
	}

	return "", fmt.Errorf("Type %s has no ID field", t.Name())
}

func (dm *DataModel) getEntityName() string {
	if dm.entity != "" {
		return dm.entity
	}

	// Try to get the entity name from a Entity() method. If not, fall back to the name of the type itself
	if obj, ok := dm.model.(EntityName); ok {
		if dm.entity = obj.Entity(); dm.entity != "" {
			return dm.entity
		}
	}

	dm.entity = reflect.ValueOf(dm.model).Type().Elem().Name()
	return dm.entity
}

func (dm *DataModel) fromCache() error {
	if err := dm.verify(); err != nil {
		return err
	}
	return cache.New(dm.Context()).Get(dm.cacheKey(), dm.model)
}

// Load loads the entity from the datastore. Must have an ID for this to work.
func (dm *DataModel) Load() error {

	if err := dm.verify(); err != nil {
		return err
	}

	/*
		if err := dm.fromCache(); err == nil {
			return nil
		}
	*/
	gcpProjectID := os.Getenv("RBD_GCP_PROJECT_ID")
	dsClient, err := datastore.NewClient(dm.Context(), gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to create DataStore client: %v", err)
		return err
	}
	if err := dsClient.Get(dm.Context(), dm.Key(), dm.model); err != nil {
		return err
	}
	// If successful, then cache so we'll have it next time
	/*
		if err := dm.Cache(); err != nil {
			return err
		}
	*/
	return nil
}

// Key returns the datastore key
func (dm *DataModel) Key() *datastore.Key {
	return &datastore.Key{Kind: dm.getEntityName(), ID: 0, Name: dm.ID(), Parent: nil}
}

// ID returns the underlying data struct's unique ID. If the supplied struct
// implements this interface, then it's result will be that of the model's EntityID()
// function. Otherwise, a new random uuid v4 will be used.
func (dm *DataModel) ID() string {

	if dm.uid != "" {
		return dm.uid
	}

	// Try to get the entity name from a Entity() method. If not, fall back to a new UUID v4
	if obj, ok := dm.model.(EntityID); ok {
		if id := obj.GetID(); id != "" {
			dm.uid = id
			return dm.uid
		}
	}

	uuid, err := NewUUID()
	if err != nil {
		panic(err)
	}

	dm.uid = uuid.String()

	// If can set the id field, do it now.
	setter, ok := dm.model.(SetID)
	if ok {
		setter.SetID(dm.uid)
	}

	if fieldName, err := dm.getIDField(); err == nil {
		reflect.ValueOf(dm.model).Elem().FieldByName(fieldName).SetString(dm.uid)
	}

	return dm.uid
}

// Save writes the entity to the datastore
func (dm *DataModel) Save() error {

	if err := dm.verify(); err != nil {
		return err
	}

	// Check if the struct has en Error() method, and use it if it does.
	if obj, ok := dm.model.(EntityError); ok {
		if err := obj.Error(); err != nil {
			return err
		}
	}
	gcpProjectID := os.Getenv("RBD_GCP_PROJECT_ID")
	dsClient, err := datastore.NewClient(dm.Context(), gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to create DataStore client: %v", err)
		return err
	}
	if _, err := dsClient.Put(dm.Context(), dm.Key(), dm.model); err != nil {
		return err
	}

	// var eg errgroup.Group
	// eg.Go(dm.Cache)
	// eg.Wait()

	var eg2 errgroup.Group
	if obj, ok := dm.model.(OnSave); ok {
		eg2.Go(obj.Save)
	}
	waitRoutine := eg2.Wait()
	return waitRoutine
}

// Context returns the internal net/context
func (dm *DataModel) Context() context.Context {
	return dm.ctx
}

func (dm *DataModel) cacheKey() string {
	return fmt.Sprintf("model.%s.%s", dm.getEntityName(), dm.ID())
}

// Cache caches the entity in memcache
func (dm *DataModel) Cache() error {
	if err := dm.verify(); err != nil {
		return err
	}
	if err := cache.New(dm.Context()).Set(dm.cacheKey(), dm.model, 0); err != nil {
		return err
	}
	if obj, ok := dm.model.(OnCache); ok {
		if err := obj.Cache(); err != nil {
			return err
		}
	}
	return nil
}

// WithContext sets the internal context for use in future operations
func (dm *DataModel) WithContext(ctx context.Context) *DataModel {
	dm.ctx = ctx
	return dm
}

// Uncache removes the cached model from cache
func (dm *DataModel) Uncache() error {
	if err := cache.New(dm.Context()).Del(dm.cacheKey()); err != nil {
		return err
	}
	if obj, ok := dm.model.(OnUncache); ok {
		if err := obj.Uncache(); err != nil {
			return err
		}
	}
	return nil
}

// Delete deletes the entity from the datastore and cache
func (dm *DataModel) Delete() error {
	gcpProjectID := os.Getenv("RBD_GCP_PROJECT_ID")
	dsClient, err := datastore.NewClient(dm.Context(), gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to create DataStore client: %v", err)
		return err
	}
	if err := dsClient.Delete(dm.Context(), dm.Key()); err != nil && err != gocache.ErrCacheMiss {
		return err
	}
	var eg errgroup.Group
	eg.Go(dm.Uncache)
	if obj, ok := dm.model.(OnDelete); ok {
		eg.Go(obj.Delete)
	}
	return eg.Wait()
}
