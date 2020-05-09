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

	nowt := time.Now()
	t.Log(nowt, nowt.YearDay())
	bt, _ := time.Parse("2006-01-02 15:04:05 MST", "1959-05-10 01:00:00 CST")
	assert.Equal(t, 60, int(CalcAge(bt)))
	bt, _ = time.Parse("2006-01-02 15:04:05 MST", "1959-05-09 01:00:00 CST")
	assert.Equal(t, 61, int(CalcAge(bt)))
	bt, _ = time.Parse("2006-01-02 15:04:05 MST", "1959-05-08 01:00:00 CST")
	assert.Equal(t, 61, int(CalcAge(bt)))
}
