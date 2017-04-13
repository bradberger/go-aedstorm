// Package aedstorm an ORM like functions which makes working with App Engine datastore entities in go a bit easier
package aedstorm

// Model is an interface for datastore entities. User-implemented structs must
// implement this interface for things to work.
type Model interface{}

// EntityName is an interface which defines the entity name of the data to be stored in the datastore.
// If a struct implements this interface, then it will be saved with this entity name. If not, the name
// of the struct itself will be used.
type EntityName interface {
	Entity() string
}

// EntityID is an interface returns the int64 ID for the datastore struct. If the supplied struct
// implements this interface, then it's result will be the ID of the struct in the datastore. Otherwise,
// a new random uuid v4 will be used
type EntityID interface {
	GetID() string
}

// EntityError is an interface which returns an error. When the model's struct implements this, it's
// called before saving the struct to the datastore. It can be used for verification of the data in the
// model, etc.
type EntityError interface {
	Error() error
}

// OnSave is an interface which defines a callback which is run after a entity is successfully saved.
// It's run parallel with the caching method, so there's no guarantee that the model is already in memcache
// when it's called. Instead, if you need to rely on it being in memcache, implement the OnCache interface.
type OnSave interface {
	Save() error
}

// OnCache is an interface which defines a callback which is run after a entity is successfully cached.
type OnCache interface {
	Cache() error
}

// OnUncache is an interface which defines a callback which is run after a entity is successfully removed from cache.
type OnUncache interface {
	Uncache() error
}

// OnDelete is an interface which defines a callback which is run after a entity is successfully deleted.
type OnDelete interface {
	Delete() error
}

// SetID is an interface, which if defined, allows the model to set it's own ID.
type SetID interface {
	SetID(string)
}
