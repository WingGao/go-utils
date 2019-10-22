package divisions

import (
	"github.com/WingGao/go-utils/ucore"
	"github.com/WingGao/go-utils/wlog"
	"github.com/json-iterator/go"
	"github.com/ungerik/go-dry"
	"time"
)

// 中国省市数据
type Division struct {
	Provinces  map[string]DivisionItem
	Cities     map[string]DivisionItem
	Areas      map[string]DivisionItem
	isInited   bool
	UpdateTime time.Time `json:",omitempty"`
	Version    uint32    `json:",omitempty"`
}

type DivisionItem struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	CityCode     string `json:"cityCode,omitempty"`
	ProvinceCode string `json:"provinceCode,omitempty"`
}

var (
	mainDivision = Division{}
)

func GetDivisions() Division {
	return mainDivision
}

// 从github获取地区信息
func LoadFromGithub() (Division, error) {
	errs := ucore.NewErrorList()

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

	if err := errs.FirstError(); err != nil {
		return mainDivision, err
	}
	mainDivision.isInited = true
	mainDivision.Version = 1
	wlog.S().Infof("LoadFromGithub 加载成功 Provinces=%d", len(mainDivision.Provinces))
	return mainDivision, nil
}

func LoadFromJson(j string) error {
	err := jsoniter.UnmarshalFromString(j, &mainDivision)
	if err != nil {
		return err
	}
	mainDivision.isInited = true
	wlog.S().Infof("LoadFromJson 加载成功 Provinces=%d", len(mainDivision.Provinces))
	return nil
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
