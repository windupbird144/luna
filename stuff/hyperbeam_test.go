package stuff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActualDamage(t *testing.T) {
	assert.Equal(t, HyperBeam{10, true}.ActualDamage(), 10*CritMultiplier)
	assert.Equal(t, HyperBeam{10, false}.ActualDamage(), 10)
}

func TestNewHyperBeam(t *testing.T) {
	assert.True(t, NewHyperBeam().Damage > 0)
}
