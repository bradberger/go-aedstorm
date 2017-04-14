package examples

import (
	"encoding/json"
	"net/http"

	aedstorm "github.com/bradberger/go-aedstorm"
)

type MyData struct {
	ID, Value string
}

// Entity defines the entity type name of the model as it's stored
// in Google App Engine datastore. If you don't define this method,
// the name of the struct itself will be used.
func (d *MyData) Entity() string {
	return "my-data"
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	// Save a new entity.
	d := &MyData{ID: "foo", Value: "bar"}
	if err := aedstorm.NewModel(&d).WithContext(ctx).Save(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dd := &MyData{ID: "foo"}
	if err := aedstorm.NewModel(&d).WithContext(ctx).Load(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// dd will now equal d after the load.
	w.Header.Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&dd)
}
