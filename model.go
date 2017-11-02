package utils

import (
	"time"
	"github.com/jinzhu/gorm"
	"reflect"
	"errors"
	"github.com/thoas/go-funk"
	"github.com/fatih/structs"
	"strings"
	"fmt"
)

var (
	errorNoSetParent = errors.New("not set parent")
)

type IModel interface {
	GModel() *gorm.DB
	GetPK() interface{}
	GetTableName() string
	GetModel() *Model
	GetDB() *gorm.DB
	SetDB(g *gorm.DB)
	Begin() *gorm.DB
	Rollback() (err error)
	Commit() (err error)
	NewScope() (*gorm.Scope, error)
	Limit(limit interface{}) *gorm.DB
	Exist(where ...interface{}) bool
	ExistID() bool
	FetchColumnValue(keys ... string) (out interface{})
	Find(out interface{}, where ...interface{}) *gorm.DB
	MakePSlice() interface{}
	BatchInsertBad(items []*Model) (err error)
	Save() error
	FirstOrCreate(where ...interface{}) (err error)
	Update(attrs ...interface{}) error
	Upsert(attrs ...interface{}) error
	GetParent() interface{}
	SetParent(p interface{})
	IsValid() error
	FormatError(err error) error
	Delete() error
	Where(query interface{}, args ...interface{}) *gorm.DB
	FormatSql(sql string, args ... interface{}) string
	SetDBOpt(name string, value interface{}) *gorm.DB
	//连贯操作
	Select(query interface{}, args ...interface{}) *gorm.DB
	// 永久生效
	SetSaveAssociations(v bool)
	//Association(column string)
	Table() *gorm.DB
}

type IGetDB interface {
	GetDB() *gorm.DB
}
type IModelParent interface {
	//检测该对象是否符合规则
	IsValid() error
	// 用户格式化数据库错误
	FormatError(err error) error
	//Delete 操作前会自动调用，检测是否可以删除
	//BeforeDelete(scope *gorm.Scope) error
	//AfterDelete(scope *gorm.Scope) error
	//BeforeUpdate(scope *gorm.Scope) (err error)
}

type Model struct {
	ID uint32   `gorm:"primary_key"`
	DB *gorm.DB `gorm:"-" json:"-"`
	// 指向父的指针
	parent            interface{} `gorm:"-"`
	associationColumn []string    `gorm:"-"`
	tx                *gorm.DB    `gorm:"-"` //事务，进行事务的时候暂存
	next              *gorm.DB    `gorm:"-"` //连贯操作需要
	OmitFields        []string    `gorm:"-" json:"-"`
}

type ModelTime struct {
	CreatedAt *time.Time `json:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty"`
	//DeletedAt *time.Time `sql:"index"`
}

func (m *Model) GetPK() interface{} {
	return m.ID
}

func (m *Model) GetTableName() string {
	scope, _ := m.NewScope()
	return scope.TableName()
}

func (m *Model) GetModel() *Model {
	return m
}

//Model.DB 是原始数据库，修改后的都应该通过这个函数读取，列入事务
func (m *Model) GetDB() *gorm.DB {
	if m.tx != nil {
		return m.tx
	}
	//if m.next != nil {
	//	return m.next
	//}
	return m.DB
}

func (m *Model) SetDB(g *gorm.DB) {
	m.DB = g
}

func (m *Model) GModel() *gorm.DB {
	return m.GetDB().Model(m.parent)
}

func (m *Model) Begin() *gorm.DB {
	m.tx = m.DB.Begin()
	return m.tx
}
func (m *Model) Rollback() (err error) {
	if m.tx != nil {
		err = m.tx.Rollback().Error
		m.tx = nil
	}
	return
}
func (m *Model) Commit() (err error) {
	if m.tx != nil {
		err = m.tx.Commit().Error
		m.tx = nil
	}
	return
}

//TODO 默认应该关闭关联存储
func (m *Model) SetSaveAssociations(v bool) {
	m.DB = m.DB.Set("gorm:save_associations", v)
}

func (m *Model) SetDBOpt(name string, value interface{}) *gorm.DB {
	return m.DB.Set(name, value)
}

func (m *Model) NewScope() (*gorm.Scope, error) {
	if m.parent == nil {
		return nil, errorNoSetParent
	}
	return m.GetDB().NewScope(m.parent), nil
}

func (m *Model) FormatColumns(keys ... string) []string {
	scope, _ := m.NewScope()
	rkeys := make([]string, len(keys))
	for i, v := range keys {
		f, ok := scope.FieldByName(v)
		if ok {
			rkeys[i] = f.DBName
		}
	}
	return rkeys
}

func (m *Model) Limit(limit interface{}) *gorm.DB {
	return m.GetDB().Limit(limit)
}

//只返回第一个
func (m *Model) FetchColumnValue(keys ... string) (out interface{}) {
	if m.ID == 0 || m.parent == nil {
		return nil
	}
	selectKeys := []string{}
	for _, key := range keys {
		v := funk.Get(m.parent, key)
		if funk.IsZero(v) {
			dbkey := m.FormatColumns(key)[0]
			selectKeys = append(selectKeys, dbkey)
		}
	}
	if len(selectKeys) > 0 {
		m.Select(selectKeys).First(m.parent)
	}
	if len(keys) > 0 {
		out = funk.Get(m.parent, keys[0])
	}
	return
}

//用了scan的方法
func (m *Model) Find(out interface{}, where ...interface{}) (db *gorm.DB) {
	if len(where) > 0 {
		db = m.Table().Where(where[0], where[1:]...).Scan(out)
	} else {
		db = m.GetDB().Find(out)
	}
	pv := reflect.ValueOf(out)
	if pv.Kind() == reflect.Ptr {
		pv = pv.Elem()
	}
	if pv.Kind() != reflect.Struct {
		// 多个结果的情况
		for i := 0; i < pv.Len(); i++ {
			outv := pv.Index(i)
			//fmt.Println("outv", i, outv.Type().Name(), outv.Interface())
			//for j := 0; j < outv.NumField(); j++ {
			//	fmt.Println("field", outv.Field(j).Type().Name())
			//}
			//vt := reflect.TypeOf(outv.Interface())
			//for j := 0; j < vt.NumMethod(); j++ {
			//	fmt.Println("method", vt.Method(i).Name)
			//}
			if _, ok := outv.Interface().(IModel); ok {
				outv.MethodByName("SetParent").Call([]reflect.Value{outv})
				outv = outv.Elem()
				outv.FieldByName("DB").Set(reflect.ValueOf(m.GetDB()))
			}
		}
	}
	return
}

//返回 *[]*ParentType
func (m *Model) FindList(where ...interface{}) (interface{}, error) {
	list := m.MakePSlice()
	err := m.Find(list, where...).Error
	return list, err
}

//创建对应父Slice切片的地址,指针 *[]*ParentType
func (m *Model) MakePSlice() interface{} {
	t := reflect.TypeOf(m.parent)
	slice := reflect.MakeSlice(reflect.SliceOf(t), 100, 100)
	arr := reflect.New(slice.Type())
	arr.Elem().Set(slice)
	return arr.Interface()
}

//deprecated
//最蛋疼的多插入
func (m *Model) BatchInsertBad(items []*Model) (err error) {
	tx := m.DB.Begin()
	for _, v := range items {
		v.DB = tx
		err = v.Save()
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func (m *Model) First(out interface{}, where ...interface{}) (err error) {
	err = m.GModel().First(out, where...).Error
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) FirstOrCreate(where ...interface{}) (err error) {
	err = m.GModel().FirstOrCreate(m.parent, where...).Error
	return m.parent.(IModelParent).FormatError(err)
}

// 会更新全部flied
func (m *Model) Save() (err error) {
	err = m.parent.(IModelParent).IsValid()
	if m.GetDB() == nil {
		err = errors.New("Model.DB is null")
	} else if err == nil {
		if m.ID > 0 {
			//更新
			err = m.GModel().Omit(append(m.OmitFields, "id", "created_at")...).Updates(m.parent).Error
		} else {
			err = m.GModel().Save(m.parent).Error
		}
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) BeforeUpdate(scope *gorm.Scope) (err error) {
	return
}

// 会更新全部flied, 强制调用Save方法
func (m *Model) FSave() (err error) {
	err = m.parent.(IModelParent).IsValid()
	if m.GetDB() == nil {
		err = errors.New("Model.DB is null")
	} else if err == nil {
		err = m.GetDB().Save(m.parent).Error
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Upsert(attrs ...interface{}) (err error) {
	err = m.parent.(IModelParent).IsValid()
	if err == nil {
		if m.ID > 0 {
			err = m.Update(attrs...)
			return
		} else {
			err = m.GetDB().Save(m.parent).Error
		}
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Updates(values interface{}, ignoreProtectedAttrs ...bool) error {
	//因为updates是部分修改，目前所以不需要检查
	//但是其实应该做
	//TODO Updates加入IsValid()操作
	//err := m.parent.(IModelParent).IsValid()
	var err error
	if err == nil {
		d := m.GetDB().Model(m.parent).Omit("id").Updates(values, ignoreProtectedAttrs ...)
		err = d.Error
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Update(attrs ...interface{}) (error) {
	err := m.parent.(IModelParent).IsValid()
	if err == nil {
		d := m.GetDB().Model(m.parent).Omit("id").Update(attrs...)
		err = d.Error
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Delete() error {
	err := m.GetDB().Delete(m.parent).Error
	return m.parent.(IModelParent).FormatError(err)
}

//批量删除
func (m *Model) BatchDelete(ids []uint32) (err error) {
	tx := m.DB.Begin()
	err = tx.Where("id IN (?)", ids).Delete(m.parent).Error
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.GetDB().Model(m.parent).Where(query, args...)
}

func (m *Model) LoadById() error {
	d := m.GetDB().First(m.parent)
	return m.parent.(IModelParent).FormatError(d.Error)
}

func (m *Model) LoadAndSetId(id uint32) error {
	d := m.GetDB().First(m.parent, "id = ?", id)
	return m.parent.(IModelParent).FormatError(d.Error)
}
func (m *Model) LoadByKey(key string, val interface{}) error {
	col := m.FormatColumns(key)[0]
	m.ID = 0
	d := m.GetDB().First(m.parent, col+" = ?", val)
	return m.parent.(IModelParent).FormatError(d.Error)
}

//连贯操作
func (m *Model) Select(query interface{}, args ...interface{}) *gorm.DB {
	m.next = m.GModel().Select(query, args...)
	return m.next
}

//TODO 更好的关联模式
//func (m *Model) checkAssociation() {
//	if len(m.associationColumn) > 0 {
//
//	}
//}
//
//func (m *Model) Association(column ...string) {
//	m.associationColumn = column
//}

func (m *Model) Table() *gorm.DB {
	s, _ := m.NewScope()
	return m.GetDB().Table(s.TableName())
}

func (m *Model) Related(value interface{}, foreignKeys ...string) error {
	d := m.GetDB().Model(m.parent).Related(value, foreignKeys...)
	return m.parent.(IModelParent).FormatError(d.Error)
}

func (m *Model) SetParent(p interface{}) {
	m.parent = p
}

func (m *Model) GetParent() interface{} {
	return m.parent
}

//使用的时候需要注意标签的导出 `structs:",flatten"`
func (m *Model) ToMap() map[string]interface{} {
	return structs.Map(m.GetParent())
}

//最好ID必须设置，不然会查询全部
func (m *Model) Exist(where ...interface{}) bool {
	item := funk.PtrOf(m.parent)
	e := m.GModel().Select("id").First(item, where...).Error
	return e == nil && item.(IModel).GetModel().ID > 0
}

func (m *Model) ExistID() bool {
	if m.ID <= 0 {
		return false
	}
	return m.Exist("id = ?", m.ID)
}

//格式化sql，添加自定义变量
// $MTABLE = 当前表名
func (m *Model) FormatSql(sql string, args ... interface{}) string {
	scope, _ := m.NewScope()
	if len(args) > 0 {
		sql = fmt.Sprintf(sql, args...)
	}
	return strings.Replace(sql, "$MTABLE", scope.TableName(), -1)
}

//被JOIN
// selectKeysMap = nil | map[string]string k:to
// SELECT <map> FROM ... JOIN <m.TableName> ON <m.TableName>.<rkey> == <g.Table>.<lkey>
func (m *Model) JoinBy(g *gorm.DB, selectKeysMap interface{}, lkey, rkey, jtype string) *gorm.DB {
	table := m.GetTableName()
	scope := GetScope(g)
	attrs := GetSelectAttrs(g)
	if selectKeysMap == nil {
		attrs = append(attrs, fmt.Sprintf("`%v`.*", table))
	} else if sm, ok := selectKeysMap.(map[string]string); ok {
		for k, to := range sm {
			if strings.Contains(k, ".") || strings.Contains(k, "(") {
				attrs = append(attrs, fmt.Sprintf("`%v` as `%v`", k, to))
			} else {
				attrs = append(attrs, fmt.Sprintf("`%v`.`%v` as `%v`", table, k, to))
			}
		}
	}
	g = g.Select(attrs).Joins(
		fmt.Sprintf("%s JOIN `%s` ON `%s`.`%s` = `%s`.`%s`", jtype, table, table, rkey, scope.TableName(), lkey))
	return g
}

//inner
func (m *Model) JoinIBy(g *gorm.DB, selectKeysMap interface{}, lkey, rkey string) *gorm.DB {
	return m.JoinBy(g, selectKeysMap, lkey, rkey, "INNER")
}

//left
func (m *Model) JoinLBy(g *gorm.DB, selectKeysMap interface{}, lkey, rkey string) *gorm.DB {
	return m.JoinBy(g, selectKeysMap, lkey, rkey, "LEFT")
}

/*
用法
func (p *Term) IsValid() error {
	errs := make([]error, 1)
	if cutils.IsValidSlug(p.Slug) {
		errs = append(errs, errors.New("Slug not valid"))
	}
	return utils.FirstError(errs...)
}
*/
func (m *Model) IsValid() (err error) {
	return nil
}

// 格式化错误
// IMPORTANT: 记得最后调用 err = p.Model.FormatError(err)
func (m *Model) FormatError(err error) error {
	return err
}

func (m *Model) BeforeDelete(scope *gorm.Scope) error {
	return nil
}
func (m *Model) AfterDelete(scope *gorm.Scope) error {
	return nil
}

func GetIDs(ms interface{}) []uint32 {
	return funk.Map(ms, func(v IModel) uint32 {
		return v.GetModel().ID
	}).([]uint32)
}

//gorm
func GetScope(g *gorm.DB) *gorm.Scope {
	return g.NewScope(g.Value)
}

func GetSelectAttrs(g *gorm.DB) []string {
	scope := GetScope(g)
	attrs := scope.SelectAttrs()
	if len(attrs) == 0 {
		return []string{fmt.Sprintf("%s.*", scope.QuotedTableName())}
	} else {
		//自动补全表名
		for i, v := range attrs {
			if strings.Contains(v, ".") || strings.HasPrefix(v, scope.TableName()+".") ||
				strings.HasPrefix(v, scope.QuotedTableName()+".") {

			} else {
				attrs[i] = fmt.Sprintf("%s.%s", scope.QuotedTableName(), v)
			}
		}
	}
	return attrs
}
func GormAddSelect(g *gorm.DB, fields ...string) *gorm.DB {
	attrs := GetSelectAttrs(g)
	attrs = append(attrs, fields...)
	return g.Select(attrs)
}

func ScopeOmitFields(scope *gorm.Scope, fields ...string) {
	if updateAttrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
		for _, v := range fields {
			delete(updateAttrs.(map[string]interface{}), v)
		}
		scope.InstanceSet("gorm:update_attrs", updateAttrs)
	}
}

func FormatSqlError(err error) (error) {
	if err != nil {
		errStr := err.Error()
		if strings.HasPrefix(errStr, "Error 1451: Cannot delete or update a parent row") {
			err = errors.New("还有其他项目在使用")
		}
	}
	return err
}

func SqlEscape(source string) (string) {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}
