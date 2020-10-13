package elasticutil

import (
	"testing"
	"github.com/WingGao/go-utils"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/stretchr/testify/assert"
)

type TestPost struct {
	utils.Model
	Title       string
	PublishTime *time.Time
	Context     string `es:"ik:max"`
	IntPtr      *int
	LongA       uint32
	Ignore      string `es:"-"`
}

func TestMain(m *testing.M) {
	log.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

func TestNewElasticModel(t *testing.T) {
	post := &TestPost{}
	post.SetParent(post)
	esMod := NewElasticModel(post)
	assert.Equal(t, esMod.Doc.Properties["PublishTime"].Type, "date")
	assert.Equal(t, esMod.Doc.Properties["IntPtr"].Type, "integer")
	assert.Equal(t, esMod.Doc.Properties["LongA"].Type, "long")
	assert.NotContains(t, "Ignore", esMod.Doc.Properties)
}
