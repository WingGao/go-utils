package utils

import (
	"testing"
	"errors"
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

func checkInterface(m IModel) bool{
	return true
}

func TestGetIDs(t *testing.T) {
	m := &Model{}
	checkInterface(m)
}
