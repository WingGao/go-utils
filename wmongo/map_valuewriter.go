package wmongo

import (
	"fmt"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strconv"
	"time"
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

var debug = false

//func (m *MapValueWriter) getParent() (interface{}, bool) {
//	if top, ok := m.stack.Peek(); ok {
//		return top.(*mode).parent
//	}
//	return nil
//}
func (m *MapValueWriter) log(format string, args ...interface{}) {
	if debug {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

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
	if top, ok := m.stack.Peek(); ok {
		top.(*mode).kind = reflect.Array
		top.(*mode).val = make(map[string]interface{})
	} else {
		//这个是顶部元素
		m.stack.Push(&mode{name: "$root", kind: reflect.Array, parent: nil, val: m.buf}) //root
	}
	return m, nil
}

func (m *MapValueWriter) WriteArrayElement() (bsonrw.ValueWriter, error) {
	top := m.peek()
	m.log("WriteArrayElement %v", top)
	arr := top.val.(map[string]interface{})
	mod := &mode{name: strconv.Itoa(len(arr)), kind: reflect.Interface, parent: arr,
		val: make(map[string]interface{})}
	arr[mod.name] = mod.val
	m.stack.Push(mod)
	return m, nil
}

func (m *MapValueWriter) WriteArrayEnd() error {
	// 转化为array
	top := m.pop()
	vmap := top.val.(map[string]interface{})
	arr := make([]interface{}, len(vmap))
	for i := 0; i < len(vmap); i++ {
		ai := strconv.Itoa(i)
		arr[i] = vmap[ai]
	}
	top.parent[top.name] = arr
	return nil
}

func (m *MapValueWriter) WriteBinary(b []byte) error {
	return m.writeVal(b)
}

func (m *MapValueWriter) WriteBinaryWithSubtype(b []byte, btype byte) error {
	return m.writeVal(primitive.Binary{Data: b, Subtype: btype})
}

func (m *MapValueWriter) WriteBoolean(v bool) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteCodeWithScope(code string) (bsonrw.DocumentWriter, error) {
	return m, m.writeVal(code)
}

func (m *MapValueWriter) WriteDBPointer(ns string, oid primitive.ObjectID) error {
	return m.writeVal(primitive.DBPointer{DB: ns, Pointer: oid})
}

func (m *MapValueWriter) WriteDateTime(dt int64) error {
	return m.writeVal(time.Unix(0, dt*int64(time.Millisecond)))
}

func (m *MapValueWriter) WriteDecimal128(v primitive.Decimal128) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteDouble(v float64) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteInt32(v int32) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteInt64(v int64) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteJavascript(code string) error {
	return m.writeVal(primitive.JavaScript(code))
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
	return m.writeVal(primitive.Regex{Pattern: pattern, Options: options})
}

func (m *MapValueWriter) WriteString(v string) error {
	return m.writeVal(v)
}

func (m *MapValueWriter) WriteDocument() (bsonrw.DocumentWriter, error) {
	if top, ok := m.stack.Peek(); ok {
		mod := top.(*mode)
		mod.kind = reflect.Map
		mod.val = make(map[string]interface{})
		mod.parent[mod.name] = mod.val
	} else {
		//这个是顶部元素
		m.stack.Push(&mode{name: "$root", kind: reflect.Map, parent: nil, val: m.buf}) //root
	}
	return m, nil
}

func (m *MapValueWriter) WriteSymbol(symbol string) error {
	return m.writeVal(primitive.Symbol(symbol))
}

func (m *MapValueWriter) WriteTimestamp(t, i uint32) error {
	return m.writeVal(primitive.Timestamp{T: t, I: i})
}

func (m *MapValueWriter) WriteUndefined() error {
	return m.writeVal(nil)
}

func NewMapValueWriter() *MapValueWriter {
	mvw := &MapValueWriter{
		buf:   make(map[string]interface{}),
		stack: lls.New(),
	}
	return mvw
}

func StructToBsonMap(s interface{}) (map[string]interface{}, error) {
	mvw := NewMapValueWriter()
	ec := bsoncodec.EncodeContext{Registry: bson.DefaultRegistry}
	enc := &bson.Encoder{}
	enc.SetContext(ec)
	enc.Reset(mvw)
	err := enc.Encode(s)
	return mvw.buf, err
}
