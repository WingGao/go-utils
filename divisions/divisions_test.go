package divisions

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	LoadFromGithub()
	t.Logf("%#v", mainDivision)
	assert.NotNil(t, mainDivision)
	assert.Equal(t, uint32(1), mainDivision.Version)
	t.Logf("Provinces=%d, Cities=%d, Areas=%d", len(mainDivision.Provinces), len(mainDivision.Cities), len(mainDivision.Areas))
}
