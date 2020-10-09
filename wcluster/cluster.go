package wcluster

import (
	"fmt"
	"github.com/WingGao/errors"
	"github.com/WingGao/go-utils/redis"
	"github.com/WingGao/go-utils/wlog"
	"github.com/ungerik/go-dry"
	"time"
)

var (
	Self *WCluster
)
// 主从工具类
type WCluster struct {
	id              string // id
	uDefId          string //用户给的名称，作为后缀用
	GroupName       string //群组名称
	IP              string
	isMaster        bool
	onMasterChanges map[string]func() // 如果状态发生改变
	keepCnt         uint64
	OnCheckMaster   func(w *WCluster) bool
}

func NewCluster(id, groupName string) (w *WCluster, err error) {
	w = &WCluster{GroupName: groupName, isMaster: false}
	w.IP = dry.RealNetIP()
	w.id = fmt.Sprintf("%s|%s", w.IP, id)
	w.onMasterChanges = make(map[string]func())
	return
}

func (w *WCluster) ID() string {
	return w.id
}
func (w *WCluster) IsMaster() bool {
	if w.OnCheckMaster != nil {
		return w.OnCheckMaster(w)
	}
	return w.isMaster
}
func (w *WCluster) IsMasterRaw() bool {
	return w.isMaster
}

// 注册
func (w *WCluster) Register() error {
	err := w.register()
	if err != nil {
		return err
	}
	go w.keepAlive()
	return nil
}

func (w *WCluster) keepAlive() {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			//wlog.S().Debugf("keep alive %#v", w)
			if w.isMaster {
				// 保活
				redis.MainClient.CtxExpire(fmt.Sprintf("util:cluster:%s_master", w.GroupName), 30*time.Second)
			} else {
				w.register()
			}
		}
	}
}

func (w *WCluster) register() error {
	isMaster, err := redis.MainClient.CtxSetNX(fmt.Sprintf("util:cluster:%s_master", w.GroupName), w.IP, 30*time.Second).Result()
	if err != nil {
		return err
	}
	if w.isMaster != isMaster || w.keepCnt == 0 { //第一次的时候也要调用
		wlog.S().Infof("状态改变 isMaster=%v", isMaster)
		for key, fn := range w.onMasterChanges {
			if fn != nil {
				wlog.S().Infof("onMasterChanges call %s", key)
				fn()
			}
		}
	}
	w.keepCnt += 1
	w.isMaster = isMaster
	if isMaster {
		err = w.registerMaster()
	} else {
		err = w.registerSlave()
	}
	return err
}

// 注册成为master
func (w *WCluster) registerMaster() error {
	return nil
}

// 注册成为slave
func (w *WCluster) registerSlave() error {
	return nil
}

// 添加到回调，这里的函数都要保证重复调用有效，而且添加的时候直接调用1次
func (w *WCluster) AddOnMasterChange(key string, fn func()) error {
	if _, ok := w.onMasterChanges[key]; ok {
		return errors.Errorf("AddOnMasterChange %s exist", key)
	}
	w.onMasterChanges[key] = fn
	fn()
	return nil
}

// 默认情况下 应该每个实例只有1个
func Init(id, group string) (err error) {
	if Self == nil {
		Self, err = NewCluster(id, group)
		err = Self.Register()
	} else {
		err = errors.New("已经初始化过")
	}
	return
}
