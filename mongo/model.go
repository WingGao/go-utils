// 用来方便MongoDB对象创建
// MgModel
//	type HugoScorePageRecord struct {
//		mongo.MgModel                                            `bson:",inline" structs:",flatten"`
//		PageID                  uint32                           `bson:"PageID"`                  //评分表ID
//		CourseHourInteractionID uint32                           `bson:"CourseHourInteractionID"` //互动ID course_hour_interaction_id
//		Result                  map[string]HugoScorePageQuestion `bson:"Result"`
//		ToUserID                uint32                           `bson:"ToUserID"` //被评价者
//		ToUser                  *Student                         `bson:"-" json:",omitempty"`
//		FromUserID              uint32                           `bson:"FromUserID"` //评论者
//	}
/*
	func (HugoScorePageRecord) TableName() string {
		return "px_hugo_page_record"
	}

	func NewHugoScorePageRecord(d interface{}) (m *HugoScorePageRecord) {
		m = &HugoScorePageRecord{}
		if d != nil {
			copier.Copy(m, d)
		}
		m.MgModel = mongo.MgModel{Session: _mgodb, DbName: mgoName}
		m.SetParent(m)
		return
	}
 */
package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/thoas/go-funk"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"reflect"
	"strings"
	"time"
	"github.com/go-errors/errors"
)

// MongoDB结构的通用
// 如果新增加域，需要在one()里也添加对应复制
type MgModel struct {
	Id bson.ObjectId `bson:"_id"`
	//IdHex   string        `bson:"-"`
	Session *mgo.Session `bson:"-" json:"-"`
	DbName  string       `bson:"-" json:"-"`
	// 指向父的指针
	parent          interface{} `bson:"-"`
	createdSessions *sll.List   `bson:"-"`
}

type IMgModel interface {
	GetModel() *MgModel
	SetModel(n *MgModel)
	SetParent(p interface{})
	GetParent() interface{}
	C() (c *mgo.Collection, s *mgo.Session)
	UpdateId(update interface{}) error
}
type IMgParent interface {
	TableName() string
	FormatError(err error) error
	BeforeDelete() error
	BeforeSave() error
}

func (m *MgModel) SetModel(n *MgModel) {
	m.Session = n.Session
	m.DbName = n.DbName
}

func (m *MgModel) GetModel() *MgModel {
	return m
}

func (m *MgModel) GetParent() interface{} {
	return m.parent
}

//基本可以用作初始化
func (m *MgModel) SetParent(p interface{}) {
	if m.createdSessions == nil {
		m.createdSessions = sll.New()
	}
	m.parent = p
}

//需要手动关闭session
func (m *MgModel) GetSession() *mgo.Session {
	if m.Session == nil {
		return nil
	}
	return m.Session.New()
}

//关闭所有获取到的session
func (m *MgModel) CloseAllSession() {
	if m.createdSessions != nil {
		for i := 0; i < m.createdSessions.Size(); i++ {
			if s, ok := m.createdSessions.Get(i); ok && s != nil {
				s.(*mgo.Session).Close()
			}
		}
		m.createdSessions.Clear()
	}
}

func (m *MgModel) C() (c *mgo.Collection, s *mgo.Session) {
	s = m.GetSession()
	if s == nil {
		return
	}
	m.createdSessions.Add(s)
	c = s.DB(m.DbName).C(m.parent.(IMgParent).TableName())
	return
}

//保存前置
func (m *MgModel) BeforeSave() error {
	return nil
}

//注意，会完全覆盖
func (m *MgModel) Save() error {
	if err := m.parent.(IMgParent).BeforeSave(); err != nil {
		return err
	}
	if !m.Id.Valid() {
		m.Id = bson.NewObjectId()
	}
	mc, ms := m.C()
	defer ms.Close()
	_, err := mc.UpsertId(m.Id, m.parent)
	if err != nil {
		return m.parent.(IMgParent).FormatError(err)
	}
	return nil
}

func (m *MgModel) LoadById(id interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := m.One(mc.FindId(ToObjectId(id)), m.parent)
	return m.pFormatError(err)
}

//注意，如果局部更新，请传$set
func (m *MgModel) UpdateId(update interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := mc.UpdateId(m.Id, update)
	return m.pFormatError(err)
}

func (m *MgModel) UpdateIdSet(update interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := mc.UpdateId(m.Id, BSet(update))
	return m.pFormatError(err)
}

func (m *MgModel) FindOne(q interface{}, out interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := m.One(mc.Find(q), out)
	return m.pFormatError(err)
}

func (m *MgModel) FindAll(q interface{}, arr interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := mc.Find(q).All(arr)
	return m.pFormatError(err)
}

func (m *MgModel) Count(q interface{}) (int, error) {
	mc, ms := m.C()
	if mc == nil {
		return 0, errors.New("not inited")
	}
	defer ms.Close()
	cnt, err := mc.Find(q).Count()
	return cnt, m.pFormatError(err)
}

func (m *MgModel) Exist(q bson.M) bool {
	mc, ms := m.C()
	defer ms.Close()
	cnt, _ := mc.Find(q).Limit(1).Count()
	return cnt > 0
}

//删除前置
func (m *MgModel) BeforeDelete() error {
	return nil
}
func (m *MgModel) DeleteId() error {
	if err := m.parent.(IMgParent).BeforeDelete(); err != nil {
		return err
	}
	mc, ms := m.C()
	defer ms.Close()
	return mc.RemoveId(m.Id)
}

//由于mgo的赋值会替换全部属性，所以需要重新赋值
func (m *MgModel) One(q *mgo.Query, out interface{}) error {
	var oldM interface{}
	if om, ok := out.(IMgModel); ok {
		oldM = funk.PtrOf(*om.GetModel())
	}
	err := q.One(out)
	if err == nil && oldM != nil {
		m := out.(IMgModel).GetModel()
		m.parent = out
		m.Session = oldM.(*MgModel).Session
		m.DbName = oldM.(*MgModel).DbName
		m.createdSessions = oldM.(*MgModel).createdSessions
	}
	return err
}
func (m *MgModel) pFormatError(err error) error {
	if err != nil {
		if p, ok := m.parent.(IMgParent); ok {
			return p.FormatError(err)
		}
	}
	return err
}

// 格式化错误
func (m *MgModel) FormatError(err error) error {
	return err
}

type MgTimeModel struct {
	UpdatedAt *time.Time `bson:"UpdatedAt"`
}

func (m *MgTimeModel) AutoNow() {
	now := time.Now()
	m.UpdatedAt = &now
}

//转换ObjectId, 支持 ObjectId, string
func ToObjectId(in interface{}) bson.ObjectId {
	//_id是24位
	if s, ok := in.(string); ok && len(s) == 24 {
		return bson.ObjectIdHex(s)
	} else if id, ok := in.(bson.ObjectId); ok {
		if len(string(id)) == 24 {
			return bson.ObjectIdHex(string(id))
		}
		return id
	}
	return bson.ObjectId("")
}

//func ScanOne(q *mgo.Query, out interface{}) (err error) {
//	out := funk.PtrOf(out.parent)
//	return
//}
//
func BEq(v interface{}) (out bson.M) {
	return bson.M{"$eq": v}
}
func BNe(v interface{}) (out bson.M) {
	return bson.M{"$ne": v}
}

//query logical
func BOr(items ...bson.M) (bson.M) {
	return bson.M{"$or": items}
}
func BAnd(items ...bson.M) (bson.M) {
	return bson.M{"$and": items}
}

func BSet(v interface{}) (out bson.M) {
	return bson.M{"$set": v}
}
func BUnset(v interface{}) (out bson.M) {
	return bson.M{"$unset": v}
}
func BCount(v interface{}) (out bson.M) {
	return bson.M{"$count": v}
}
func BSum(v interface{}) (out bson.M) {
	return bson.M{"$sum": v}
}
func BAvg(v interface{}) (bson.M) {
	return bson.M{"$avg": v}
}

func BAddFields(field string, v interface{}) (out bson.M) {
	return bson.M{"$addFields": bson.M{field: v}}
}

func BMatch(v interface{}) (bson.M) {
	return bson.M{"$match": v}
}
func BGroup(v interface{}) (bson.M) {
	return bson.M{"$group": v}
}

func BIn(v interface{}) (out bson.M) {
	return bson.M{"$in": v}
}

func BInField(field string, v interface{}) (out bson.M) {
	return bson.M{field: bson.M{"$in": v}}
}

func BExists(v interface{}) (bson.M) {
	return bson.M{"$exists": v}
}

//array
func BElemMatch(v interface{}) (out bson.M) {
	return bson.M{"$elemMatch": v}
}

//忽略某些
func GetMSetIgnore(obj interface{}, bsonFields ...string) (bm bson.M) {
	setM := bson.M{}
	objt := reflect.TypeOf(obj)
	objv := reflect.ValueOf(obj)
	if objt.Kind() == reflect.Ptr {
		objt = objt.Elem()
		objv = objv.Elem()
	}
	info, err1 := bson.GetStructInfo(objt)
	if err1 != nil {
		return
	}
	ignoreMap := make(map[string]bool, len(bsonFields))
	for _, f := range bsonFields {
		ignoreMap[f] = true
	}

	for _, v := range info.FieldsList {
		if v.Key == "_id" { //忽略_id
			continue
		}
		if _, ok := ignoreMap[v.Key]; !ok {
			setv := objv
			if len(v.Inline) > 0 {
				//inline
				for _, inlineNum := range v.Inline {
					setv = setv.Field(inlineNum)
				}
			} else {
				setv = setv.Field(v.Num)
			}
			setM[v.Key] = setv.Interface()
		}
	}
	bm = BSet(setM)

	return
}

//有些版本的MongoDB会报错，可以使用该方法忽律
func IgnoreDuplicateKey(err error) error {
	if err != nil {
		if strings.HasPrefix(err.Error(), "E11000 duplicate key error") {
			return nil
		}
	}
	return err
}
