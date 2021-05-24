package oswinutil

import (
	"golang.org/x/sys/windows/registry"
)

// RegistryDeleteAllInKey 删除注册表项下的所有内容
func RegistryDeleteAllInKey(key registry.Key, sub string) error {
	k, err := registry.OpenKey(key, sub, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	kStat, _ := k.Stat()
	//fmt.Println("删除" + sub)
	if kStat.SubKeyCount > 0 {
		subKeys, _ := k.ReadSubKeyNames(int(kStat.SubKeyCount))
		for _, skey := range subKeys {
			if err = RegistryDeleteAllInKey(k, skey); err != nil {
				return err
			}
		}
	}
	// 删除自己
	return registry.DeleteKey(key, sub)
}
