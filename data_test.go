package aedstorm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"google.golang.org/appengine/aetest"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

var (
	ctx  context.Context
	done func()
)

func TestMain(m *testing.M) {
	var err error
	ctx, done, err = aetest.NewContext()
	if err != nil {
		log.Fatalf("Could not get test context: %v", err)
	}
	defer done()
	os.Exit(m.Run())
}

type testModelErr struct {
	ID string
}

func (tme *testModelErr) Error() error {
	return errors.New("custom error")
}

type testModel struct {
	ID string
}

func (tm *testModel) GetID() string {
	if tm.ID == "" {
		tm.ID = fmt.Sprintf("%d", time.Now().Unix())
	}
	return tm.ID
}

func (tm *testModel) Entity() string {
	return "testModel"
}

func (tm *testModel) Save(ctx context.Context) error {
	return nil
}

type testModelWithIDTag struct {
	Email string `datastore:"id"`
}

func (tm *testModelWithIDTag) GetID() string {
	if tm.Email == "" {
		tm.Email = fmt.Sprintf("%d", time.Now().Unix())
	}
	return tm.Email
}

func (tm *testModelWithIDTag) Entity() string {
	return "testModelWithIDTag"
}

func (tm *testModelWithIDTag) Save(ctx context.Context) error {
	return nil
}

type testModelWithNoID struct {
	Str string
}

func (tm *testModelWithNoID) GetID() string {
	return ""
}

func (tm *testModelWithNoID) Entity() string {
	return "testModelWithNoID"
}

func (tm *testModelWithNoID) Save(ctx context.Context) error {
	return nil
}

func TestNewModel(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, dm.model, tm)
}

type nonStructModel bool

func (n *nonStructModel) GetID() string {
	return ""
}

func (n *nonStructModel) Save(ctx context.Context) error {
	return nil
}

func TestNewModelWithNoPtr(t *testing.T) {
	assert.Panics(t, func() {
		tm := testModel{}
		NewModel(tm)
	})
	assert.Panics(t, func() {
		nsm := nonStructModel(true)
		NewModel(&nsm)
	})
}

func TestNewModelVerifyFlag(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	dm.verified = true
	assert.NoError(t, dm.verify())
}

func TestNewModelVerifyContext(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, ErrNoContext, dm.verify())
}

func TestNewModelVerifyModel(t *testing.T) {
	dm := &DataModel{ctx: ctx}
	assert.Equal(t, ErrNoModel, dm.verify())
}

func TestNewModelVerify(t *testing.T) {
	dm := &DataModel{ctx: ctx, model: &testModel{}}
	if !assert.NoError(t, dm.verify()) {
		return
	}
	assert.True(t, dm.verified)
	assert.Equal(t, dm.idFieldName, "ID")
}

func TestNewModelIDFieldName(t *testing.T) {
	dm := &DataModel{ctx: ctx, model: &testModelWithNoID{}}
	assert.EqualError(t, dm.verify(), "Type testModelWithNoID has no ID field")
}

func TestGetIDField(t *testing.T) {
	tm := &testModelWithIDTag{}
	dm := NewModel(tm)
	idField, err := dm.getIDField()
	assert.NoError(t, err)
	assert.Equal(t, "Email", idField)
}

func TestDataModelUse(t *testing.T) {
	tm := &testModelWithIDTag{}
	dm := NewModel(tm)
	dm.WithContext(ctx)
	assert.NotNil(t, dm.ctx)
	assert.NotNil(t, dm.Context())
}

func TestModelCacheKeyName(t *testing.T) {
	u := fmt.Sprintf("%d", time.Now().Unix())
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, "model.testModel."+u, dm.cacheKey())
}

func TestFromCacheWithNoContext(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, ErrNoContext, dm.fromCache())
}

func TestFromCache(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	assert.NoError(t, dm.Cache())
	assert.NoError(t, dm.fromCache())
}

func TestCacheWithNoContext(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, ErrNoContext, dm.Cache())
}

func TestSaveWithNoContext(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, ErrNoContext, dm.Save())
}

func TestSaveWithErrorInterface(t *testing.T) {
	tm := &testModelErr{}
	dm := NewModel(tm).WithContext(ctx)
	dm.verified = true
	assert.EqualError(t, dm.Save(), "custom error")
}

func TestDelete(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	assert.NoError(t, dm.Save())
	assert.NoError(t, dm.Delete())
}

func TestDeleteWithNoCache(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	assert.NoError(t, dm.Save())
	assert.NoError(t, dm.Uncache())
	assert.Error(t, dm.Delete())
}

func TestDeleleteUnsaved(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	assert.Error(t, dm.Delete())
}

type ondeleteTestModel struct {
	ID string
}

func (*ondeleteTestModel) Delete() error {
	return errors.New("delete me")
}

func TestDeleteWithOnDelete(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	assert.NoError(t, dm.Save())
	assert.EqualError(t, dm.Delete(), "delete me")
}

func TestGetKey(t *testing.T) {
	tm := &testModel{ID: "foobar"}
	dm := NewModel(tm)
	dm.WithContext(ctx)
	k := dm.Key()
	assert.NotNil(t, k)
	assert.Nil(t, k.Parent())
	assert.Equal(t, k.StringID(), tm.ID)
	assert.Equal(t, "testModel", k.Kind())
	assert.Equal(t, int64(0), k.IntID())
}

func TestEntityName(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Equal(t, "testModel", dm.getEntityName())
}

func TestEntityNamePreset(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	dm.entity = "foobar"
	assert.Equal(t, "foobar", dm.getEntityName())
}

func TestLoadWhenVerifyFails(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.Error(t, dm.Load())
}

func TestLoadFromCache(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	dm.verified = true
	assert.NoError(t, dm.Cache())
	assert.NoError(t, dm.Load())
}

func TestLoadFromDatastore(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm).WithContext(ctx)
	dm.verified = true
	assert.NoError(t, dm.Save())
	assert.NoError(t, dm.Uncache())
	assert.NoError(t, dm.Load())
}

type modelWithNoEntity struct{}

func TestEntityNameWithNoInterface(t *testing.T) {
	m := &modelWithNoEntity{}
	dm := NewModel(m)
	assert.Equal(t, "modelWithNoEntity", dm.getEntityName())
}

func TestEntityID(t *testing.T) {
	tm := &testModel{}
	dm := NewModel(tm)
	assert.NotEmpty(t, dm.ID())
}

type testModelWithIDField struct {
	ID string
}

func TestSetEntityID(t *testing.T) {
	tm := &testModelWithIDField{}
	dm := NewModel(tm)
	assert.NotEmpty(t, dm.ID())
	assert.Equal(t, dm.ID(), tm.ID)
}
