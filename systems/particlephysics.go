package systems

import (
	"github.com/bcokert/engo-test/logging"

	"github.com/engoengine/math"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

// ParticleComponent represents a physical particle object - one with position, mass, and velocity, but no rotation
// It's sufficiently described by one point in space for Newtonian Physics, except collisions which are done on bounding boxes
// It must be made legal - use NewParticleComponent to guarantee this
type ParticleComponent struct {
	InvMass          float32               // 0 means infinite mass
	SpaceComponent   common.SpaceComponent // Contains Position and Rotation
	Velocity         engo.Point            // per second
	ForceAccumulator engo.Point            // the sum of all forces on this particle since the last integration step (eg: collisions, etc). Doesn't include environment forces, like gravity
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
		Velocity: velocity,
	}
}

// particleEntity is the entity that the ParticlePhysicsSystem operates on
type particleEntity interface {
	BasicEntity() *ecs.BasicEntity
	ParticleComponent() *ParticleComponent
}

// The ParticlePhysicsSystem is the root physics engine. It receives updates as often as
// possible from the game loop, which will also render as often as possible. However the
// physics system will only run at a specified rate - it will run multiple times if the
// rendering has caused the simulation to fall behind
// The physics system takes all entities in the system and integrates them, handling
// collisions and constraints as necessary.
type ParticlePhysicsSystem struct {
	TopLeft        engo.Point // built in bounding box for all entities
	BottomRight    engo.Point // built in bounding box for all entities
	Gravity        engo.Point // the gravitational force, in meters/second^2
	DampingFactor  float32    // the damping factor that reduces the velocity slightly (sort of like an aerodynamic impulse due to drag) but is primarily for numerical stability
	SimulationRate int        // updates per second
	Log            logging.Logger
	simulationStep float32 // seconds per update, the constant step of each simulation step
	simulationAcc  float32 // seconds since the last simulation. When greater than simulationStep, simulation occurs
	entities       map[uint64]particleEntity
}

// Add adds a new entity to the system
func (s *ParticlePhysicsSystem) Add(entity particleEntity) {
	s.Log.Info("Adding particle physics object", map[string]interface{}{"entityID": entity.BasicEntity().ID(), "particleComponent": entity.ParticleComponent()})
	s.entities[entity.BasicEntity().ID()] = entity
}

// Remove removes a entity from the system, by its entity id
func (s *ParticlePhysicsSystem) Remove(entity ecs.BasicEntity) {
	if _, ok := s.entities[entity.ID()]; ok {
		s.Log.Info("Removing particle physics object", map[string]interface{}{"entity": entity})
		delete(s.entities, entity.ID())
	}
}

// New is called every time the system is added to a world
func (s *ParticlePhysicsSystem) New(world *ecs.World) {
	s.simulationAcc = 0
	s.simulationStep = 1.0 / float32(s.SimulationRate)
	if s.entities == nil {
		s.entities = make(map[uint64]particleEntity, 10)
	}
}

// Update is called as often as the game loop can. If the accumulator has accumulated
// at least 1 step worth of time, then steps are simulated, decrementing the accumulator each time,
// until the accumulator is less than the step.
// If for some reason the physics simulation takes longer than a step, physics will be simulated a
// maximum of 10 times before yielding control, which will have an effect similar to time slowing down
func (s *ParticlePhysicsSystem) Update(dt float32) {
	s.simulationAcc += dt

	// Simulate physics in steps until we've caught up to real time or hit the limit of 10
	// Any remainder less than the simulationStep can be interpolated by the renderer
	for i := 0; s.simulationAcc > s.simulationStep && i < 10; i++ {
		for _, entity := range s.entities {
			s.integrate(entity)
		}
	}
}

func (s *ParticlePhysicsSystem) integrate(e particleEntity) {
	dt := s.simulationStep
	body := e.ParticleComponent()
	// width := body.SpaceComponent.Width
	// height := body.SpaceComponent.Height
	// left := body.SpaceComponent.Position.X
	// right := body.SpaceComponent.Position.X + width
	// top := body.SpaceComponent.Position.Y
	// bottom := body.SpaceComponent.Position.Y + height

	invm := body.InvMass

	s.Log.Debug("Before Integration: ", logging.F{"position": body.SpaceComponent.Position, "velocity": body.Velocity, "invm": invm})

	// Calculate the net force on the object
	var netF engo.Point

	// Environment forces
	netF.Add(s.Gravity)

	// Add other external forces on the object, then clear them
	netF.Add(body.ForceAccumulator)
	body.ForceAccumulator.Set(0, 0)

	// Acceleration from net force
	// a = F/m
	acceleration := netF
	acceleration.MultiplyScalar(invm)
	s.Log.Debug("Integration Forces", logging.F{"gravity": s.Gravity, "externalforces": body.ForceAccumulator, "netF": netF, "acceleration": acceleration})

	// Update the position
	// p = p + v*t + 0.5*a*t^2
	v := body.Velocity                  // v
	v.MultiplyScalar(dt)                // v*t
	body.SpaceComponent.Position.Add(v) // p = p + v*t
	a := acceleration                   // a
	a.MultiplyScalar(0.5 * dt * dt)     // 0.5*a*t^2
	body.SpaceComponent.Position.Add(a) // p = p + v*t + 0.5*a*t^2

	// Update the velocity
	// v = v*damp^t + a*t
	dampt := math.Pow(s.DampingFactor, dt) // damp^t
	body.Velocity.MultiplyScalar(dampt)    // v = v*damp^t
	a = acceleration                       // a
	a.MultiplyScalar(dt)                   // a*t
	body.Velocity.Add(a)                   // v = v*damp^t + a*t
	s.Log.Debug("After Integration: ", logging.F{"position": body.SpaceComponent.Position, "velocity": body.Velocity})

	// if p.X+dp.X <= s.BoxLeft {
	// 	p.X = 0
	// 	dp.X = dp.X * -1 * (1 - s.Elasticity.X)
	// 	dp.Y *= (1 - s.Friction.Y)
	// } else if p.X+dp.X+w >= s.BoxRight {
	// 	p.X = s.BoxRight - w
	// 	dp.X = dp.X * -1 * (1 - s.Elasticity.X)
	// 	dp.Y *= (1 - s.Friction.Y)
	// } else {
	// 	p.X = p.X + dp.X
	// }
	// dp.X = dp.X + ddp.X + s.Gravity.X
	// if math.Abs(float64(dp.X)) < 0.001 {
	// 	dp.X = 0
	// }

	// accelY := ddp.Y + s.Gravity.Y
	// if entity.BasicEntity().ID() == 1 {
	// 	log.Debug("Tick: y: %v   dy: %v   ddy: %v", p.Y, dp.Y, accelY)
	// }
	// if p.Y+dp.Y <= s.BoxTop {
	// 	p.Y = 0
	// 	dp.Y = dp.Y * -1 * (1 - s.Elasticity.Y)
	// 	if dp.X != 0 {
	// 		dp.X *= (1 - s.Friction.X)
	// 	}
	// } else if p.Y+dp.Y+h >= s.BoxBottom {
	// 	if entity.BasicEntity().ID() == 1 {
	// 		log.Printf("Collision with floor")
	// 	}
	// 	p.Y = s.BoxBottom - h
	// 	dp.Y = dp.Y * -1 * (1 - s.Elasticity.Y)
	// 	if entity.BasicEntity().ID() == 1 {
	// 		log.Debug("After adjusting y and dp.Y: y: %v   dy: %v", p.Y, dp.Y)
	// 	}
	// 	if accelY > dp.Y*-1 {
	// 		// the entity did not bounce off fast enough to overcome acceleration
	// 		dp.Y = 0
	// 	}
	// 	if dp.X != 0 {
	// 		dp.X *= (1 - s.Friction.X)
	// 	}
	// } else {
	// 	p.Y = p.Y + dp.Y
	// 	dp.Y += accelY
	// }
	// if math.Abs(float64(dp.Y)) < 0.001 {
	// 	dp.Y = 0
	// }
}
