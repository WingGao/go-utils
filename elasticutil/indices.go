package elasticutil

type IndexRep struct {
	Mappings ElasticModel `json:"mappings"`
	Settings IndexSetting `json:"settings"`
}

type IndexSetting struct {
	Blocks map[string]interface{} `json:"blocks"`
}

type IndexSettingBlocks struct {
	ReadOnlyAllowDelete bool `json:"read_only_allow_delete"`
}
