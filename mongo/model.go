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
	m.createdSessions.Add(s)
	c = s.DB(m.DbName).C(m.parent.(IMgParent).TableName())
	return
}

//注意，会完全覆盖
func (m *MgModel) Save() error {
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
	err := m.one(mc.FindId(ToObjectId(id)), m.parent)
	return m.pFormatError(err)
}

func (m *MgModel) UpdateId(update interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := mc.UpdateId(m.Id, update)
	return m.pFormatError(err)
}

func (m *MgModel) Find(q interface{}, out interface{}) error {
	mc, ms := m.C()
	defer ms.Close()
	err := m.one(mc.Find(q), out)
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
	defer ms.Close()
	cnt, err := mc.Find(q).Count()
	return cnt, m.pFormatError(err)
}

//由于mgo的赋值会替换全部属性，所以需要重新赋值
func (m *MgModel) one(q *mgo.Query, out interface{}) error {
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

//转换ObjectId, 支持 ObjectId, string
func ToObjectId(in interface{}) bson.ObjectId {
	if id, ok := in.(bson.ObjectId); ok {
		return id
	} else if s, ok := in.(string); ok {
		return bson.ObjectIdHex(s)
	}
	return bson.ObjectId("")
}

//func ScanOne(q *mgo.Query, out interface{}) (err error) {
//	out := funk.PtrOf(out.parent)
//	return
//}
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

func BIn(field string, v interface{}) (out bson.M) {
	return bson.M{field: bson.M{"$in": v}}
}

func MarshalJSONStr(m interface{}) string {
	js, _ := bson.MarshalJSON(m)
	return string(js)
}
