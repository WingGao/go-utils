package wcluster

import (
	"github.com/WingGao/go-utils/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	internal.LoadFromFile()
	internal.LoadRedis()
	os.Exit(m.Run())
}
func TestInit(t *testing.T) {
	w1, _ := NewCluster("1", "test")
	err1 := w1.Register()
	w2, _ := NewCluster("2", "test")
	err2 := w2.Register()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.True(t, w1.isMaster)
	assert.False(t, w2.isMaster)
}

func TestInit2(t *testing.T) {
	w1, _ := NewCluster("1", "test")
	err1 := w1.Register()
	assert.NoError(t, err1)
	assert.False(t, w1.isMaster)
	<-time.After(30 * time.Second)
	assert.Equal(t, true, w1.isMaster)
}
