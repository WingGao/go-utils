package wmongo

import (
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type mode struct {
	name string
	kind reflect.Kind
	val  interface{} //指向当前
	//parent interface{} //指向父级
	parent map[string]interface{}
}
type MapValueWriter struct {
	buf   map[string]interface{}
	stack *lls.Stack //*mode
}

//func (m *MapValueWriter) getParent() (interface{}, bool) {
//	if top, ok := m.stack.Peek(); ok {
//		return top.(*mode).parent
//	}
//	return nil
//}
func (m *MapValueWriter) peek() *mode {
	top, _ := m.stack.Peek()
	return top.(*mode)
}
func (m *MapValueWriter) pop() *mode {
	top, _ := m.stack.Pop()
	return top.(*mode)
}
func (m *MapValueWriter) WriteDocumentElement(name string) (bsonrw.ValueWriter, error) {
	top := m.peek()
	m.stack.Push(&mode{name: name, kind: reflect.Interface, val: nil, parent: top.val.(map[string]interface{})})
	return m, nil
}

func (m *MapValueWriter) WriteDocumentEnd() error {
	m.stack.Pop()
	return nil
}

// 对于数组，所有数组都是map[string]interface{}类型，最后再转换数组
func (m *MapValueWriter) WriteArray() (bsonrw.ArrayWriter, error) {
	return m, nil
}
func (m *MapValueWriter) WriteArrayElement() (bsonrw.ValueWriter, error) {
	return m, nil
}

func (m *MapValueWriter) WriteArrayEnd() error {
	return nil
}

func (m *MapValueWriter) WriteBinary(b []byte) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteBinaryWithSubtype(b []byte, btype byte) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteBoolean(bool) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteCodeWithScope(code string) (bsonrw.DocumentWriter, error) {
	return m, nil
}

func (m *MapValueWriter) WriteDBPointer(ns string, oid primitive.ObjectID) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteDateTime(dt int64) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteDecimal128(primitive.Decimal128) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteDouble(float64) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteInt32(int32) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteInt64(int64) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteJavascript(code string) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteMaxKey() error {
	panic("implement me")
}

func (m *MapValueWriter) WriteMinKey() error {
	panic("implement me")
}
func (m *MapValueWriter) writeVal(val interface{}) error {
	curr := m.pop()
	curr.parent[curr.name] = val
	return nil
}
func (m *MapValueWriter) WriteNull() error {
	return m.writeVal(nil)
}

func (m *MapValueWriter) WriteObjectID(v primitive.ObjectID) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteRegex(pattern, options string) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteString(v string) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteDocument() (bsonrw.DocumentWriter, error) {
	return m, nil
}

func (m *MapValueWriter) WriteSymbol(symbol string) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteTimestamp(t, i uint32) error {
	panic("implement me")
}

func (m *MapValueWriter) WriteUndefined() error {
	panic("implement me")
}

func NewMapValueWriter() *MapValueWriter {
	mvw := &MapValueWriter{
		buf:   make(map[string]interface{}),
		stack: lls.New(),
	}
	mvw.stack.Push(&mode{name: "$root", kind: reflect.Map, parent: nil, val: mvw.buf}) //root
	return mvw
}

func a() {
	mvw := NewMapValueWriter()
	enc := new(bson.Encoder)
	enc.Reset(mvw)
}
