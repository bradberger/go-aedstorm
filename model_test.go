package aedstorm

// import (
// 	"fmt"
// 	"net/http"
//
// 	"golang.org/x/net/context"
//
// 	"google.golang.org/appengine"
// )
//
// type MyData struct {
// 	ID string
// }
//
// func (d *MyData) GetID() {
// 	if d.ID == "" {
// 		uuid, _ := NewUUID()
// 		d.ID = fmt.Sprintf("%#v", uuid)
// 	}
// 	return d.ID
// }
//
// func (d *MyData) Save(ctx context.context) error {
// 	return nil
// }
//
// func (d MyData) Entity() string {
// 	return "MyData"
// }
//
// func myHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := appengine.NewContext(r)
//
// }
//
// func ModelExample() {
//
// }
