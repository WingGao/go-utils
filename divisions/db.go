package divisions

type AddressBase struct {
	Province string `gorm:"comment:'省市'"` //
	City     string `gorm:"comment:'城市'"` //
	Area     string `gorm:"comment:'地区'"`
	Line1    string `gorm:"size:255;comment:'详细地址'"`
}

func (m AddressBase) ProvinceCn() string {
	if p, ok := mainDivision.Provinces[m.Province]; ok {
		return p.Name
	}
	return ""
}
func (m AddressBase) CityCn() string {
	if p, ok := mainDivision.Cities[m.City]; ok {
		return p.Name
	}
	return ""
}
func (m AddressBase) AreaCn() string {
	if p, ok := mainDivision.Areas[m.Area]; ok {
		return p.Name
	}
	return ""
}
