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
package wmongo

import (
	"context"
	"github.com/WingGao/go-utils"
	"github.com/WingGao/go-utils/ucore"
	"github.com/globalsign/mgo/bson"
	"github.com/go-errors/errors"
	icontext "github.com/kataras/iris/context"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
	"time"
)

// MongoDB结构的通用
// 如果新增加域，需要在one()里也添加对应复制
type MgModel struct {
	Id primitive.ObjectID `bson:"_id"`
	//IdHex   string        `bson:"-"`
	Client *mongo.Client `bson:"-" json:"-"`
	DbName string        `bson:"-" json:"-"`
	// 指向父的指针
	parent     interface{} `bson:"-"`
	softDelete bool        `bson:"-"`
}

type IMgModel interface {
	GetModel() *MgModel
	SetModel(n *MgModel)
	New() interface{}
	SetParent(p interface{})
	GetParent() interface{}
	C() (c *MgCollection, s *mongo.Client)
	UpdateId(update interface{}) error
}
type IMgParent interface {
	TableName() string
	GetModel() *MgModel
	FormatError(err error) error
	BeforeDelete() error
	BeforeSave() error
}

func (m *MgModel) SetModel(n *MgModel) {
	m.Client = n.Client
	m.DbName = n.DbName
}

func (m *MgModel) GetModel() *MgModel {
	return m
}

//得到一个基础父类，可以被重写，值不复制
func (m *MgModel) New() interface{} {
	n := ucore.PtrOf(m.parent)
	reflect.ValueOf(n).Elem().FieldByName("MgModel").Set(reflect.ValueOf(MgModel{
		Client: m.Client, DbName: m.DbName, softDelete: m.softDelete, parent: n,
	}))
	return n
}

func (m *MgModel) GetParent() interface{} {
	return m.parent
}

//基本可以用作初始化
func (m *MgModel) SetParent(p interface{}) {
	if m.DbName == "" {
		m.DbName = utils.DefaultConfig.Mongodb.DBName
	}
	m.parent = p
}

//需要手动关闭session
func (m *MgModel) GetClient() *mongo.Client {
	if m.Client == nil {
		return nil
	}
	return m.Client
}

//关闭所有获取到的session
func (m *MgModel) CloseAllSession() {
	//if m.createdSessions != nil {
	//	for i := 0; i < m.createdSessions.Size(); i++ {
	//		if s, ok := m.createdSessions.Get(i); ok && s != nil {
	//			s.(*mgo.Session).Close()
	//		}
	//	}
	//	m.createdSessions.Clear()
	//}
}

func (m *MgModel) SetSoftDelete(sf bool) {
	m.softDelete = sf
}

func (m *MgModel) C() (c *MgCollection, s *mongo.Client) {
	s = m.GetClient()
	if s == nil {
		return
	}
	c = NewMgCollection(s.Database(m.DbName).Collection(m.parent.(IMgParent).TableName()))
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
	if m.Id.IsZero() {
		m.Id = primitive.NewObjectID()
	}
	if ti, ok := m.parent.(IMgTimeModel); ok { //添加时间
		ti.AutoNow()
	}
	mc, _ := m.C()
	res, err := mc.UpsertId(m.Id, m.parent)
	if err != nil {
		return m.parent.(IMgParent).FormatError(err)
	}
	if res.UpsertedID != nil { // 只有新插入的有
		m.Id = res.UpsertedID.(primitive.ObjectID)
	}
	return nil
}

func (m *MgModel) SetId(id interface{}) {
	m.Id = ToObjectId(id)
}

func (m *MgModel) LoadById(id interface{}) error {
	mc, _ := m.C()
	//err := m.One(mc.FindId(ToObjectId(id)), m.parent)
	err := mc.FindId(id, m.parent)
	return m.pFormatError(err)
}

//注意，如果局部更新，请传$set
func (m *MgModel) UpdateId(update interface{}) error {
	mc, _ := m.C()
	_, err := mc.UpdateId(m.Id, update)
	return m.pFormatError(err)
}

func (m *MgModel) UpdateIdSet(update interface{}) error {
	mc, _ := m.C()
	//defer ms.Close()
	_, err := mc.UpdateId(m.Id, BSet(update))
	return m.pFormatError(err)
}

func (m *MgModel) FindOne(q interface{}, out interface{}) error {
	mc, _ := m.C()
	//err := m.One(mc.Find(q), out)
	res := mc.FindOne(context.Background(), q)
	err := DecodeSingleRes(res, out)
	return m.pFormatError(err)
}

func (m *MgModel) FindAll(q interface{}, out interface{}, opts ...*options.FindOptions) error {
	mc, _ := m.C()
	c, err := mc.Find(context.Background(), q, opts...)
	if err != nil {
		return err
	} else {
		err = c.All(context.Background(), out)
	}
	return m.pFormatError(err)
}

func (m *MgModel) Count(q interface{}) (int, error) {
	//mc, _ := m.C()
	//if mc == nil {
	//	return 0, errors.New("not inited")
	//}
	//defer ms.Close()
	//cnt, err := mc.Find(q).Count()
	//return cnt, m.pFormatError(err)
	return 0, errors.New("not implement")
}

func (m *MgModel) Exist(q bson.M) bool {
	//mc, ms := m.C()
	//defer ms.Close()
	//cnt, _ := mc.Find(q).Limit(1).Count()
	//return cnt > 0
	panic("not implement")
	return false
}

//删除前置
func (m *MgModel) BeforeDelete() error {
	return nil
}

func (m *MgModel) DeleteId() error {
	if err := m.parent.(IMgParent).BeforeDelete(); err != nil {
		return err
	}
	mc, _ := m.C()
	if m.softDelete { //如果开启了软删除，则将文档移动至一个_del集合
		old := m.New().(IMgModel)
		old.GetModel().LoadById(m.Id)
		dc := NewMgCollection(m.Client.Database(m.DbName).Collection(m.parent.(IMgParent).TableName() + "_del"))
		res, err1 := dc.UpsertId(old.GetModel().Id, old)
		if res.UpsertedCount != 1 || err1 != nil {
			return err1
		}
	}
	_, err := mc.RemoveId(m.Id)
	return err
}
func (m *MgModel) DeleteOne(q interface{}) (*mongo.DeleteResult, error) {
	mc, _ := m.C()
	return mc.RemoveOne(q)
}

//由于mgo的赋值会替换全部属性，所以需要重新赋值
//func (m *MgModel) One(q *mgo.Query, out interface{}) error {
//	var oldM interface{}
//	if om, ok := out.(IMgModel); ok {
//		oldM = funk.PtrOf(*om.GetModel())
//	}
//	err := q.One(out)
//	if err == nil && oldM != nil {
//		m := out.(IMgModel).GetModel()
//		m.parent = out
//		m.Client = oldM.(*MgModel).Client
//		m.DbName = oldM.(*MgModel).DbName
//	}
//	return err
//}
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

// 获取创建时间
func (m *MgModel) CreatedAt() *time.Time {
	if !m.Id.IsZero() {
		t := m.Id.Timestamp()
		return &t
	}
	return nil
}

// 从iris的请求中合并对象
func (m *MgModel) MergeFromIris(ctx icontext.Context) error {
	if bs, err := ctx.GetBody(); err == nil {
		return mergeFromJson(m.parent, bs)
	} else {
		return err
	}
}

type MgTimeModel struct {
	UpdatedAt *time.Time `bson:"UpdatedAt"`
}

type IMgTimeModel interface {
	AutoNow()
}

func (m *MgTimeModel) AutoNow() {
	now := time.Now()
	m.UpdatedAt = &now
}

type MgSoftDelete struct {
	DeletedAt *time.Time `gorm:"index:idx_deleted_at" json:",omitempty"` //deleted_at
}

//转换ObjectId, 支持 ObjectId, string
func ToObjectId(in interface{}) (oid primitive.ObjectID) {
	//_id是24位
	if s, ok := in.(string); ok && len(s) == 24 {
		oid, _ = primitive.ObjectIDFromHex(s)
		return
	} else if id, ok := in.(primitive.ObjectID); ok {
		return id
	}
	return primitive.NilObjectID
}
func ToObjectIds(arr interface{}) []primitive.ObjectID {
	ids := []primitive.ObjectID{}
	funk.ForEach(arr, func(v interface{}) {
		ids = append(ids, ToObjectId(v))
	})
	return ids
}

func DecodeSingleRes(sr *mongo.SingleResult, out interface{}) (err error) {
	//var bs []byte
	err = sr.Err()
	if err == nil {
		//bs, err = sr.DecodeBytes()
		//if err == nil {
		//	err = bson.Unmarshal(bs, out)
		//}
		err = sr.Decode(out)
	}
	return
}

//func ScanOne(q *mgo.Query, out interface{}) (err error) {
//	out := funk.PtrOf(out.parent)
//	return
//}
//

//有些版本的MongoDB会报错，可以使用该方法忽律
func IgnoreDuplicateKey(err error) error {
	if err != nil {
		if strings.HasPrefix(err.Error(), "E11000 duplicate key error") {
			return nil
		}
	}
	return err
}
