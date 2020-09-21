package elasticutil

import (
	"bytes"
	"context"
	"github.com/WingGao/errors"
	"github.com/WingGao/go-utils/ucore"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	"reflect"
)

type IEsModel interface {
	GetModel() *EsModel
	//SetModel(n *EsModel)
	SetParent(p interface{})
	GetParent() IEsParent
	//C() (c *MgCollection, s *mongo.Client)
	UpdateId(update interface{}) error
}
type IEsParent interface {
	TableName() string
	GetModel() *EsModel
	New() interface{}
	//FormatError(err error) error
	//BeforeDelete() error
	//BeforeSave() error
}

type EsModel struct {
	Id     string
	Client *elasticsearch.Client `json:"-"`
	// 指向父的指针
	parent interface{} `json:"-"`
}

func (m *EsModel) GetModel() *EsModel {
	return m
}

func (m *EsModel) GetParent() IEsParent {
	return m.parent.(IEsParent)
}

func (m *EsModel) SetParent(p interface{}) {
	m.parent = p
}

// out *[]*Parent{}
func (m *EsModel) FindAll(out interface{}, body interface{}, q ...func(*esapi.SearchRequest)) (rep *SearchResult, err error) {
	var buf bytes.Buffer
	jsoniter.NewEncoder(&buf).Encode(body)
	q = append(q, m.Client.Search.WithContext(context.Background()), m.Client.Search.WithIndex(m.GetParent().TableName()),
		m.Client.Search.WithBody(&buf))
	res, err1 := m.Client.Search(q...)
	defer res.Body.Close()
	if err1 != nil {
		err = errors.WithStack(err1)
		return
	}
	if res.IsError() {
		err = errors.New(res.String())
	} else {
		rep = &SearchResult{}
		err = jsoniter.NewDecoder(res.Body).Decode(rep)
		ouL := m.SearchResultToList(rep)
		if ouL != nil {
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(ouL))
		}
	}
	return rep, errors.WrapSkip(err, 0)
}

func (m *EsModel) SearchResultToList(rep *SearchResult) interface{} {
	if rep.Hits == nil || rep.Hits.Hits == nil || len(rep.Hits.Hits) == 0 {
		return nil
	}
	outL := ucore.MakeSlice(m.parent, len(rep.Hits.Hits))
	outR := reflect.ValueOf(outL).Elem()
	for _, hit := range rep.Hits.Hits {
		if hit.Source == nil {
			continue
		}
		v := m.GetParent().New()
		v.(IEsParent).GetModel().Id = hit.Id
		if err := jsoniter.Unmarshal(hit.Source, v); err == nil {
			outR = reflect.Append(outR, reflect.ValueOf(v))
		}
	}
	return outR.Interface()
}

type EsError struct {
	Error string
}

type SearchResult elastic.SearchResult
