package utils

import (
	"testing"
	"errors"
	"reflect"
	"github.com/stretchr/testify/assert"
)

type testMod2 struct {
	Model
}
type testMod struct {
	Model
}

func newMod() *testMod {
	m := &testMod{}
	m.SetParent(m)
	return m
}
func newMod2() *testMod2 {
	m := &testMod2{}
	m.SetParent(m)
	return m
}

func (m *testMod) IsValid() error {
	return errors.New("error for test")
}

func TestModel_IsValid(t *testing.T) {
	mod := newMod()
	mod2 := newMod2()

	if mod.Save() == nil {
		t.Fatal("check mod.IsValid failed")
	}
	err := mod2.Save()
	if err != nil {
		t.Fatal("check mod2.IsValid failed")
	}
	t.Log("IsVaild Pass")
}

func checkInterface(m IModel) bool {
	return true
}

func TestGetIDs(t *testing.T) {
	m := &Model{}
	checkInterface(m)
}

type TestA struct {
	Model
	FieldB uint32
}

func TestModel_New(t *testing.T) {
	a := &TestA{}
	a.ID = 1
	a.SetParent(a)
	ta := reflect.TypeOf(a).Elem()
	n := reflect.New(ta)
	nele := n.Elem()
	for i := 0; i < nele.NumField(); i++ {
		t.Log(nele.Field(i).Type().String())
	}
	nele.FieldByName("Model").Set(reflect.ValueOf(Model{ID: 2}))
	b := PtrOf(a)
	reflect.ValueOf(b).Elem().FieldByName("Model").Set(reflect.ValueOf(Model{ID: 3}))
	assert.Equal(t, uint32(2), nele.Addr().Interface().(*TestA).ID)
	assert.Equal(t, uint32(1), a.ID)
	assert.Equal(t, uint32(3), b.(*TestA).ID)
}
