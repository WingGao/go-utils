package uuid

import (
	"github.com/WingGao/go-utils"
	"os"
	"testing"
	"github.com/WingGao/go-utils/redis"
	"time"
)

func TestMain(m *testing.M) {
	testConf, _ := utils.LoadConfig(os.Getenv("NXPT_GO_CONF"))
	redis.LoadClient(testConf.Redis)
	Init()
	os.Exit(m.Run())
}

func TestNextID(t *testing.T) {
	cnt := 0
	for {
		time.Sleep(2 * time.Second)
		t.Log(NextID())
		cnt++
		if cnt > 15 {
			return
		}
	}
}
