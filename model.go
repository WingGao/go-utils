package utils

import (
	"database/sql"
	"fmt"
	ucore "github.com/WingGao/go-utils/ucore"
	"github.com/fatih/structs"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"reflect"
	"strings"
	"time"
)

var (
	errorNoSetParent = errors.New("not set parent")
)

type IModel interface {
	New() interface{}
	GModel() *gorm.DB
	PrimaryKey() string
	PrimaryKeyZero() bool
	PrimaryKeyValue() interface{}
	GetTableName() string
	GetModel() *Model
	GetDB() *gorm.DB
	SetDB(g *gorm.DB)
	Begin() *gorm.DB
	Rollback() (err error)
	Commit() (err error)
	NewScope() (*gorm.Scope, error)
	Limit(limit interface{}) *gorm.DB
	IsLoaded() bool
	LoadAndSetId(id uint32) error
	LoadByPk(pk interface{}) error
	Exist(where ...interface{}) bool
	ExistPk() bool
	FetchColumnValue(keys ...string) (out interface{})
	Find(out interface{}, where ...interface{}) error
	RawFind(out interface{}, where ...interface{}) *gorm.DB
	//创建对应父Slice切片的地址,指针 *[]*ParentType
	MakePSlice() interface{}
	BatchInsertBad(items []*Model) (err error)
	Save() error
	FirstOrCreate(where ...interface{}) (err error)
	Update(attrs ...interface{}) error
	Updates(values interface{}, ignoreProtectedAttrs ...bool) (err error)
	Upsert() error
	GetParent() interface{}
	SetParent(p interface{})
	IsValid() error
	FormatError(err error) error
	Delete() error
	Where(query interface{}, args ...interface{}) *gorm.DB
	FormatSql(sql string, args ...interface{}) string
	SetDBOpt(name string, value interface{}) *gorm.DB
	//连贯操作
	Select(query interface{}, args ...interface{}) *gorm.DB
	// 永久生效
	SetSaveAssociations(v bool)
	//Association(column string)
	Table() *gorm.DB
	// 获得 select 的语句
	GetFieldsSql(ignore []string, prefix string, asprefix string) string
}

type IGetDB interface {
	GetDB() *gorm.DB
}
type IModelParent interface {
	//检测该对象是否符合规则
	IsValid() error
	// 用户格式化数据库错误
	FormatError(err error) error
	FormatFields(str string) string
	SetPrimaryKey() (key interface{}, err error)
	CheckUnique() error // 检查唯一性，在使用软删除的时候同时使用，默认BeforeSave的时候调用
	//Delete 操作前会自动调用，检测是否可以删除
	//BeforeDelete(scope *gorm.Scope) error
	//AfterDelete(scope *gorm.Scope) error
	//BeforeUpdate(scope *gorm.Scope) (err error)
}

/*

## 自定义主键

	type Order struct {
		utils.Model        `gorm:"-"`
		utils.ModelTime    `structs:",flatten"`
		Serial    uint64   `gorm:"primary_key;auto_increment:false"` //订单号,非自增
	}
	...
	// 自动创建主键
	func (m *Order) SetPrimaryKey() (key interface{}, err error) {
		key, err = uuid.NextID()
		return
	}

*/
type Model struct {
	ID uint32   `gorm:"primary_key" bson:"ID"`
	DB *gorm.DB `gorm:"-" json:"-" bson:"-" form:"-" es:"-"`
	// 指向父的指针
	parent            interface{} `gorm:"-"`
	associationColumn []string    `gorm:"-"`
	tx                *gorm.DB    `gorm:"-"` //事务，进行事务的时候暂存
	next              *gorm.DB    `gorm:"-"` //连贯操作需要
	OmitFields        []string    `gorm:"-" json:"-" bson:"-" form:"-"`
}

func (m *Model) PrimaryKey() string {
	scope, _ := m.NewScope()
	return scope.PrimaryKey()
}

func (m *Model) PrimaryKeyZero() bool {
	scope, _ := m.NewScope()
	return scope.PrimaryKeyZero()
}

func (m *Model) PrimaryKeyValue() interface{} {
	scope, _ := m.NewScope()
	return scope.PrimaryKeyValue()
}

//生成一个新的主键，一般用于自定义主键
func (m *Model) SetPrimaryKey() (key interface{}, err error) {
	return nil, nil
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

// 开始事务处理
// 后续结束事务需要手动调用其他方法
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

func (m *Model) AutoEnd(commit bool) (err error) {
	if commit {
		err = m.Commit()
	} else {
		err = m.Rollback()
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

func (m *Model) FormatColumns(keys ...string) []string {
	scope, _ := m.NewScope()
	rkeys := make([]string, len(keys))
	for i, v := range keys {
		if ucore.StrHasLowerPrefix(v) {
			rkeys[i] = v
		} else if f, ok := scope.FieldByName(v); ok {
			rkeys[i] = f.DBName
		}
	}
	return rkeys
}

// 创建索引
// 约定 uix_xxx 是唯一索引 其他为普通索引
func (m *Model) CreateIndexes(indexList map[string][]string) (err error) {
	for name, columns := range indexList {
		if strings.HasPrefix(name, "uix_") {
			err = m.GModel().AddUniqueIndex(name, columns...).Error
		} else {
			err = m.GModel().AddIndex(name, columns...).Error
		}
		if err != nil {
			return
		}
	}
	return
}

func (m *Model) Limit(limit interface{}) *gorm.DB {
	return m.GetDB().Limit(limit)
}

//只返回第一个
func (m *Model) FetchColumnValue(keys ...string) (out interface{}) {
	if m.ID == 0 || m.parent == nil {

	} else {
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
	}
	if len(keys) > 0 {
		out = funk.Get(m.parent, keys[0])
	}
	return
}

func (m *Model) Find(out interface{}, where ...interface{}) (err error) {
	err = m.RawFind(out, where...).Error
	return m.parent.(IModelParent).FormatError(err)
}

//用了scan的方法, 没有limit, 会自动判断软删除
func (m *Model) RawFind(out interface{}, where ...interface{}) (db *gorm.DB) {
	if len(where) > 0 {
		db = m.GModel().Where(where[0], where[1:]...).Scan(out)
	} else {
		db = m.GetDB().Find(out)
	}
	pv := reflect.ValueOf(out)
	if pv.Kind() == reflect.Ptr {
		pv = pv.Elem()
	}
	if pv.Kind() == reflect.Slice {
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
	err := m.RawFind(list, where...).Error
	return list, err
}

//返回 []*ParentType
func (m *Model) FindList2(where ...interface{}) (interface{}, error) {
	list := m.MakePSlice()
	err := m.RawFind(list, where...).Error
	val := reflect.ValueOf(list)
	return val.Elem().Interface(), err
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

// 更新局部，只逻辑判断主键持否为nil，并不会去查数据库
func (m *Model) Save() (err error) {
	if m.GetDB() == nil {
		err = errors.New("Model.DB is null")
	} else if err == nil {
		scope, _ := m.NewScope()
		if !scope.PrimaryKeyZero() {
			//更新
			err = scope.DB().Omit(append(m.OmitFields, scope.PrimaryKey(), COL_CREATED_AT)...).Updates(m.parent).Error
		} else { //创建
			_, err = m.parent.(IModelParent).SetPrimaryKey()
			if err == nil {
				err = scope.DB().Create(m.parent).Error
			}
		}
	}
	return m.parent.(IModelParent).FormatError(err)
}

// 会更新全部flied, 强制调用Save方法
func (m *Model) FSave() (err error) {
	if m.GetDB() == nil {
		err = errors.New("Model.DB is null")
	} else if err == nil {
		err = m.GetDB().Save(m.parent).Error
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Create() (err error) {
	err = m.GetDB().Create(m.parent).Error
	return m.parent.(IModelParent).FormatError(err)
}

// 主动更具pk去查询数据库来判断是创建还是更新
func (m *Model) Upsert() (err error) {
	if m.ExistPk() {
		err = m.GetDB().Save(m.parent).Error
	} else {
		err = m.GetDB().Create(m.parent).Error
	}
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Updates(values interface{}, ignoreProtectedAttrs ...bool) (err error) {
	d := m.GetDB().Model(m.parent).Omit("id").Updates(values, ignoreProtectedAttrs...)
	err = d.Error
	return m.parent.(IModelParent).FormatError(err)
}

// 更新单个属性
func (m *Model) Update(attrs ...interface{}) (err error) {
	if m.PrimaryKeyZero() {
		return errors.New("pk is nil")
	}

	d := m.GetDB().Model(m.parent).Omit("id").Update(attrs...)
	err = d.Error

	return m.parent.(IModelParent).FormatError(err)
}

// 增加一个数据库列，但是这个col一定不为null
func (m *Model) Increase(col string, step int) (err error) {
	if m.PrimaryKeyZero() {
		return errors.New("pk is nil")
	}
	d := m.GetDB().Model(m.parent).Update(col, gorm.Expr(fmt.Sprintf("`%s` + ?", col), step))
	err = d.Error
	return m.parent.(IModelParent).FormatError(err)
}

//只能删除自己
func (m *Model) Delete() error {
	scope, _ := m.NewScope()
	if scope.PrimaryKeyZero() {
		return errors.New("Delete require PK")
	}
	err := scope.DB().Delete(m.parent).Error
	return m.parent.(IModelParent).FormatError(err)
}

//更具id删除
func (m *Model) DeleteByIDs(ids []uint32) (err error) {
	tx := m.DB.Begin()
	mod := m.New().(IModel).GetModel()
	for _, id := range ids {
		mod.ID = id
		err = mod.Delete()
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) DeleteBy(where ...interface{}) error {
	err := m.GetDB().Delete(m.parent, where...).Error
	return m.parent.(IModelParent).FormatError(err)
}

func (m *Model) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.GetDB().Model(m.parent).Where(query, args...)
}

//判断是否加载，需要重写
func (m *Model) IsLoaded() bool {
	return false
}

func (m *Model) LoadById() error {
	d := m.GetDB().First(m.parent)
	return m.parent.(IModelParent).FormatError(d.Error)
}

func (m *Model) LoadAndSetId(id uint32) error {
	d := m.GetDB().First(m.parent, "id = ?", id)
	return m.parent.(IModelParent).FormatError(d.Error)
}

func (m *Model) LoadByPk(pk interface{}) error {
	d := m.GetDB().First(m.parent, fmt.Sprintf("%s = ?", m.PrimaryKey()), pk)
	return m.parent.(IModelParent).FormatError(d.Error)
}

// key是struct里的Field，不是数据库的列名
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

func (m *Model) ToJSON() string {
	j, _ := jsoniter.MarshalToString(m)
	return j
}

//最好ID必须设置，不然会查询全部;如果没有定义的时候没有ID，则无法生效
func (m *Model) Exist(where ...interface{}) bool {
	item := m.New()
	scope, _ := m.NewScope()
	e := scope.DB().Select(scope.PrimaryKey()).First(item, where...).Error
	return e == nil && !item.(IModel).PrimaryKeyZero()
}

func (m *Model) ExistPk() bool {
	if m.PrimaryKeyZero() {
		return false
	}
	return m.Exist(m.FormatSql("$PK = ?"), m.PrimaryKeyValue())
}

//格式化sql，添加自定义变量
// $MTABLE = 当前表名
// $PK = 主键
func (m *Model) FormatSql(sql string, args ...interface{}) (out string) {
	scope, _ := m.NewScope()
	if len(args) > 0 {
		sql = fmt.Sprintf(sql, args...)
	}
	out = sql
	out = strings.Replace(out, "$MTABLE", scope.TableName(), -1)
	out = strings.Replace(out, "$PK", scope.PrimaryKey(), -1)
	return
}

//被JOIN
// selectKeysMap = nil | map[string]string k:to | string
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
	} else if sq, ok := selectKeysMap.(string); ok {
		attrs = []string{sq}
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

func (m *Model) FormatFields(str string) string {
	return ""
}

// 格式化错误
// IMPORTANT: 记得最后调用 err = p.Model.FormatError(err)
func (m *Model) FormatError(err error) error {
	if err != nil {
		errStr := err.Error()
		if errStr == "record not found" {
			err = errors.Wrap("不存在", 1)
		} else if strings.HasPrefix(errStr, "Error 1062: Duplicate entry") {
			if nl := m.parent.(IModelParent).FormatFields(errStr); nl != "" {
				err = errors.Wrap(fmt.Errorf("%s已存在", nl), 1)
			}
		}
	}
	return err
}

func (m *Model) BeforeSave(scope *gorm.Scope) error {
	return m.parent.(IModelParent).CheckUnique()
}

// 用来手动检查unique的情况，在软删除使用的情况下需要用到
// 默认会在`BeforeSave`的时候调用
func (m *Model) CheckUnique() error {
	return nil
}

func (m *Model) BeforeDelete(scope *gorm.Scope) error {
	// 将isactive=NULL 是为了让unique_index生效
	// 所有使用unique的索引应该都与isactive组合
	_, hasActField := scope.FieldByName("IsActive")
	if hasActField {
		scope.Raw(fmt.Sprintf(
			"UPDATE %v SET `isactive`=NULL%v",
			scope.QuotedTableName(),
			addExtraSpaceIfExist(scope.CombinedConditionSql()),
		)).Exec()
		//TODO 更好的方式
		// 清空scope不然后续无法使用，
		scope.SQL = ""
		scope.SQLVars = []interface{}{}
	}
	return nil
}
func (m *Model) AfterDelete(scope *gorm.Scope) error {
	return nil
}

//得到一个基础父类，可以被重写，值不复制
func (m *Model) New() interface{} {
	n := ucore.PtrOf(m.parent)
	reflect.ValueOf(n).Elem().FieldByName("Model").Set(reflect.ValueOf(Model{parent: n, DB: m.GetDB()}))
	return n
}

func (m *Model) GetFieldsSql(ignore []string, table string, asprefix string) string {
	scope, _ := m.NewScope()
	fields := scope.Fields()
	sb := ucore.StringBuilder{}
	for _, f := range fields {
		if f.IsIgnored {
			continue
		}
		if funk.ContainsString(ignore, f.DBName) {
			continue
		}
		if sb.Len() > 0 {
			sb.Write(", ")
		}
		if table != "" {
			sb.Write(table, ".")
		}
		sb.Write("`", f.DBName, "`")
		if asprefix != "" {
			sb.Write(" AS `", asprefix, f.DBName, "`")
		}
	}
	return sb.String()
}

// TODO 更好的自定义方式
const (
	COL_CREATED_AT = "inserttime"
	COL_UPDATED_AT = "updattime"
)

type IModelTime interface {
	UnsetTime()
}

type ModelTime struct {
	CreatedAt *time.Time `gorm:"Column:inserttime;index:idx_inserttime;default:CURRENT_TIMESTAMP;comment:'插入时间'" json:",omitempty"`
	UpdatedAt *time.Time `gorm:"Column:updatetime;index:idx_updatetime;default:CURRENT_TIMESTAMP;comment:'更新时间'" json:",omitempty"` //updated_at
	//DeletedAt *time.Time `sql:"index"`
}

func (m *ModelTime) UnsetTime() {
	m.CreatedAt = nil
	m.UpdatedAt = nil
}
func (m *ModelTime) ColNameCreatedAt() string {
	return COL_CREATED_AT
}
func (m *ModelTime) ColNameUpdateAt() string {
	return COL_UPDATED_AT
}

type ModelSoftDelete struct {
	DeletedAt *time.Time `gorm:"index:idx_deleted_at" json:",omitempty"` //deleted_at
	IsActive  *bool      `gorm:"Column:isactive;index:idx_isactive;DEFAULT:1;COMMENT:'逻辑删除(1:保留,0:删除)'"`
}

//func (m ModelSoftDelete) GetDeleteWhere() string {
//	return "deleted_at IS NULL"
//}

func (m ModelSoftDelete) GetActiveWhere() string {
	return "deleted_at IS NULL"
}

func (ModelSoftDelete) UnDelete(mi IModel) error {
	err := mi.Table().Unscoped().Model(mi.GetParent()).Updates(map[string]interface{}{"deleted_at": nil, "isactive": true}, false).Error
	return mi.FormatError(err)
}

func GetIDs(ms interface{}) []uint32 {
	return funk.Map(ms, func(v IModel) uint32 {
		if v != nil {
			if m := v.GetModel(); m != nil {
				return m.ID
			}
		}
		return 0
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
			} else if strings.HasPrefix(v, "DISTINCT") {
				at := strings.Split(v, " ")
				attrs[i] = fmt.Sprintf("DISTINCT %s.%s", scope.QuotedTableName(), at[1])
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

func GormIn(g *gorm.DB, field string, arg interface{}) *gorm.DB {
	if funk.IsEmpty(arg) {
		return g
	}
	return g.Where(fmt.Sprintf("%s IN (?)", field), arg)
}
func GormNotIn(g *gorm.DB, field string, arg interface{}) *gorm.DB {
	if funk.IsEmpty(arg) {
		return g
	}
	return g.Where(fmt.Sprintf("%s NOT IN (?)", field), arg)
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
		} else if strings.HasPrefix(errStr, "Error 1062: Duplicate entry") {
			err = errors.New("已存在")
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

// sql.Rows的长度
func RowsLength(rows *sql.Rows) (l int) {
	for rows.Next() {
		l++
	}
	return
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
