package physics

import (
	"engo.io/engo"
	"engo.io/engo/common"
)

// ParticleComponent contains the particle-physics related properties of an entity.
// It's sufficiently described by one point in space for Newtonian Physics, except collisions which are done on bounding spheres
// It must be made legal - use NewParticleComponent to guarantee this
type ParticleComponent struct {
	InvMass          float32               // 0 means infinite mass
	SpaceComponent   common.SpaceComponent // Contains Position and Rotation
	Velocity         engo.Point            // per second
	ForceAccumulator engo.Point            // the sum of all forces on this particle since the last integration step (eg: collisions, etc). Doesn't include environment forces, like gravity
	Restitution      float32               // the coefficient of restitution is a factor for the percentange of velocity this object retains after a collision
}

// NewParticleComponent constructs a legal component and provides some helpers
func NewParticleComponent(width, height, mass float32, position, velocity engo.Point) ParticleComponent {
	if mass < 0 {
		mass = 0
	}

	inverseMass := float32(0)
	if mass > 0 {
		inverseMass = 1 / mass
	}

	return ParticleComponent{
		InvMass: inverseMass,
		SpaceComponent: common.SpaceComponent{
			Position: position,
			Width:    width,
			Height:   height,
			Rotation: 0,
		},
		Velocity:         velocity,
		ForceAccumulator: engo.Point{},
		Restitution:      0.7,
	}
}
