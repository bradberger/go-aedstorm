[![GoDoc](https://godoc.org/github.com/bradberger/go-aedstorm?status.svg)](https://godoc.org/github.com/bradberger/go-aedstorm)
[![Build Status](https://semaphoreci.com/api/v1/brad/go-aedstorm/branches/master/shields_badge.svg)](https://semaphoreci.com/brad/go-aedstorm)
[![codecov](https://codecov.io/gh/bradberger/go-aedstorm/branch/master/graph/badge.svg)](https://codecov.io/gh/bradberger/go-aedstorm)

This is an ORM like package which makes working with App Engine datastore
entities in go a bit easier. The name `aedstorm` stands for App Engine DataSTore
ORM.

It makes working with the datastore and basic data structs quite simple.

```
import (
	"encoding/json"
	"net/http"

	aedstorm "github.com/bradberger/go-aedstorm"
)

type MyData struct {
	ID, Value string
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

	// dd.Value  will now equal d.Value after the load.
	w.Header.Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&dd)
}
```

For more examples, see the `examples` subdirectory.

I'm planning on adding more documentation in the future, including some
more advanced interfaces which you can implement on a model, and queries,
but for now this should give you an idea how easy it is to get going.


Planned improvements

- [ ] Better documentation, more basic examples
- [ ] More documentation around current interfaces which can be implemented
- [ ] Add support for loading/saving related entities