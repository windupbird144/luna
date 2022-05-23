package stuff

import "math/rand"

const (
	DamageBaseline = 50
	DamageDelta    = 9
	CritOdds       = 0.125
	CritMultiplier = 2
)

type HyperBeam struct {
	Damage int
	Crit   bool
}

func (h HyperBeam) ActualDamage() int {
	if h.Crit {
		return h.Damage * CritMultiplier
	}
	return h.Damage
}

func NewHyperBeam() HyperBeam {
	// calculate base damage
	delta := rand.Intn(DamageDelta)
	if rand.Float32() < 0.5 {
		delta *= -1
	}
	// calculate crit
	crit := rand.Float32() < CritOdds
	return HyperBeam{
		Damage: DamageBaseline + delta,
		Crit:   crit,
	}
}
