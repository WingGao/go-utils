package elasticutil

import (
	"github.com/fatih/structs"
	"github.com/WingGao/go-utils"
	"reflect"
	"time"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func init() {
	//log.Level = logrus.DebugLevel
}

type ElasticModel struct {
	Doc ElasticMapping `json:"doc"`
}

type ElasticMapping struct {
	Dynamic    string                     `json:"dynamic,omitempty"`
	Properties map[string]ElasticProperty `json:"properties"`
}
type ElasticProperty struct {
	Type  string `json:"type"`
	Index bool   `json:"index"` //searchable
	// text
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
}

func NewStruct(mod interface{}) (s *structs.Struct) {
	s = structs.New(mod)
	s.TagName = "es"
	return
}

/*
	直接标签
	{
		FieldA string `es:"type:text,analyzer:xxx,search_analyzer:xxx"`
	}
	标签列表：
		type 类型
		analyzer
		search_analyzer
		ik:(max|smart) 表示直接使用ik_max_word或ik_smart
		index:(true|false)
 */
func NewElasticModel(mod utils.IModel) *ElasticModel {
	esMod := &ElasticModel{}
	mapping := &ElasticMapping{
		Dynamic:    "false",
		Properties: map[string]ElasticProperty{},
	}
	//pkName := mod.PrimaryKey()
	esStruct := structs.New(mod.GetParent())
	esStruct.TagName = "es"

	for _, f := range esStruct.Fields() {
		log.Debugf("Name: %s, Kind %v", f.Name(), f.Kind())
		prop := ElasticProperty{Index: true}
		sType := f.GoValue().Type()
		if sType.Kind() == reflect.Ptr {
			sType = sType.Elem()
		}
		fieldValue := reflect.Indirect(reflect.New(sType))
		switch sType.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Uint8:
			prop.Type = "short"
		case reflect.Int, reflect.Int32, reflect.Uint16:
			prop.Type = "integer"
		case reflect.Uint, reflect.Uint32, reflect.Int64:
			prop.Type = "long"
		case reflect.Uint64: //not support
		case reflect.String:
			prop.Type = "text"
		case reflect.Struct:
			if _, ok := fieldValue.Interface().(time.Time); ok {
				prop.Type = "date"
			}
		}
		tags := structs.ParseTag2(f.Tag("es"))

		if t, ok := tags["type"]; ok {
			prop.Type = t
		}
		if prop.Type == "" { //跳过
			continue
		}

		if t, ok := tags["analyzer"]; ok {
			prop.Analyzer = t
		}
		if t, ok := tags["search_analyzer"]; ok {
			prop.SearchAnalyzer = t
		}
		if t, ok := tags["ik"]; ok {
			if t == "max" {
				prop.Analyzer = "ik_max_word"
				prop.SearchAnalyzer = "ik_max_word"
			}
		}
		if t, ok := tags["index"]; ok {
			prop.Index = t == "true"
		}

		mapping.Properties[f.Name()] = prop
	}
	esMod.Doc = *mapping
	return esMod
}
