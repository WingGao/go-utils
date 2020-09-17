package orm

// 统一orm库
type IModel interface {
}

type IModelParent interface {
	TableName() string
	GetModel() IModel
}
