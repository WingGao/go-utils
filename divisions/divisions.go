package divisions

import (
	"github.com/ungerik/go-dry"
	"github.com/WingGao/go-utils"
	"time"
	"github.com/json-iterator/go"
	"path/filepath"
)

// 中国省市数据
type Division struct {
	Provinces  map[string]DivisionItem
	Cities     map[string]DivisionItem
	Areas      map[string]DivisionItem
	isInited   bool
	UpdateTime time.Time
	Version    uint32
}

type DivisionItem struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	ParentCode string `json:"parent_code,omitempty"`
}

var (
	mainDivision = Division{}
)

func init() {
	errs := utils.NewErrorList()
	tempFile := filepath.Join(utils.BinPath(), "divisions-data.json")
	old, err := dry.FileGetBytes(tempFile)
	if err != nil {
		j1, err1 := dry.FileGetBytes("https://github.com/modood/Administrative-divisions-of-China/raw/master/dist/provinces.json", 10*time.Second)
		errs.AppendE(err1)
		mainDivision.Provinces, err1 = parseJson(j1)
		errs.AppendE(err1)
		j2, err2 := dry.FileGetBytes("https://github.com/modood/Administrative-divisions-of-China/raw/master/dist/cities.json", 10*time.Second)
		errs.AppendE(err2)
		mainDivision.Cities, err2 = parseJson(j2)
		errs.AppendE(err2)
		j3, err3 := dry.FileGetBytes("https://github.com/modood/Administrative-divisions-of-China/raw/master/dist/areas.json", 10*time.Second)
		errs.AppendE(err3)
		mainDivision.Areas, err3 = parseJson(j3)
		errs.AppendE(err3)
		mainDivision.UpdateTime = time.Now()
	} else {
		errs.AppendE(jsoniter.Unmarshal(old, &mainDivision))
	}

	errs.Panic()
	mainDivision.isInited = true
	mainDivision.Version = 1
	dry.FileSetJSON(tempFile, mainDivision)
}

func parseJson(bs []byte) (out map[string]DivisionItem, err error) {
	arr := []DivisionItem{}
	err = jsoniter.Unmarshal(bs, &arr)
	if err != nil {
		return
	}
	out = make(map[string]DivisionItem, len(arr))
	for _, v := range arr {
		out[v.Code] = v
	}
	return
}
