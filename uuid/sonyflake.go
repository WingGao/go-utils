package uuid

import (
	"github.com/sony/sonyflake"
	"time"
	"github.com/WingGao/go-utils/redis"
	"fmt"
)

var (
	client    *sonyflake.Sonyflake
	_config   Config
	machineId uint16
)

const (
	heartBeatDuration = 10 * time.Second
)

type Config struct {
	ProjectName string
	IgnoreExist bool
}

func heartbeat() {
	for {
		//5秒的冗余时间
		redis.MainClient.Expire(getMachineKey(_config.ProjectName, machineId), int(heartBeatDuration/time.Second)+5)
		time.Sleep(heartBeatDuration)
		//fmt.Println("sonyflake heartbeat")
	}
}

func getMachineKey(prefix string, u uint16) string {
	return fmt.Sprintf("%s_wing-utils-snoyflake-%d", prefix, u)
}

func Init(cnf Config) {
	if client != nil {
		return
	}
	_config = cnf
	client = sonyflake.NewSonyflake(sonyflake.Settings{
		//StartTime: time.Now(),
		CheckMachineID: func(u uint16) bool {
			if v, _ := redis.MainClient.Incr(getMachineKey(cnf.ProjectName, u)); v > 1 {
				if cnf.IgnoreExist {
					//忽律已存在
				} else {
					fmt.Println("[WingUtils] sonyflake machine id", u, "already exist")
					return false
				}
			}
			return true
		},
	})
	if client == nil {
		panic("sonyflake create failed")
	}
	id, err := client.NextID()
	if err != nil {
		panic(err)
	}
	dec := sonyflake.Decompose(id)
	machineId = uint16(dec["machine-id"])
	fmt.Printf("[WingUtils] sonyflake machine id: %d\n", machineId)
	go heartbeat()
}

// 销毁
func Destroy() {
	redis.MainClient.Del(getMachineKey(_config.ProjectName, machineId))
	machineId = 0
	client = nil
}

func NextID() (id uint64, err error) {
	return client.NextID()
}
