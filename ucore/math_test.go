package ucore

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalcAge(t *testing.T) {
	ut := time.Unix(637593080, 0)
	t.Log(ut)
	age := CalcAge(ut)
	assert.Equal(t, uint32(30), age)
}
