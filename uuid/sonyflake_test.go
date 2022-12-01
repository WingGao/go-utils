package uuid

import (
	"github.com/WingGao/go-utils/redis"
	"os"
	"testing"
	"time"
	"wingao.net/webproj/core"
)

func TestMain(m *testing.M) {
	testConf, _ := core.LoadConfig(os.Getenv("NXPT_GO_CONF"))
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
