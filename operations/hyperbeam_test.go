package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that ActualDamage takes into account whether the hyper beam is/isn't a critical hit
func TestActualDamage(t *testing.T) {
	assert.Equal(t, HyperBeam{10, false}.ActualDamage(), 10, "ActualDamage should just return the damage for a regular hit")
	assert.Equal(t, HyperBeam{10, true}.ActualDamage(), 10*CritMultiplier, "ActaulDamage did not take the critical hit into account")
}

// Test that damage is a positive number
func TestNewHyperBeam(t *testing.T) {
	assert.True(t, NewHyperBeam().Damage > 0, "NewHyperBeam should have positive damage")
}
