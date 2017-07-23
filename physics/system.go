package physics

import (
	"engo.io/ecs"
	"github.com/bcokert/engo-test/metrics"
)

const (
	// MaxPhysicsIterations is the maximum number of simulations that will occur before giving up
	// Typically it is the render that causes more than 1 to be needed
	// But if the physics itself is taking more time that the time it is simulating, this prevents infinite loops
	MaxPhysicsIterations = 10
)

// particle is the entity that the ParticlePhysicsSystem operates on
type particle interface {
	BasicEntity() *ecs.BasicEntity
	ParticleComponent() *ParticleComponent
}

// The ParticlePhysicsSystem is the integration point between engo (the game engine) and a ParticleEngine (physics engine)
// It handles game loop related events (like delta time), and provides the game an interface to the physics engine.
// It doesn't have any logic for the game except what is needed to manage the physics engine
type ParticlePhysicsSystem struct {
	ParticleEngine *ParticleEngine // the physics engine for this system
	SimulationRate int             // updates per second
	simulationAcc  float32         // seconds since the last simulation. When greater than simulationStep, simulation occurs
	simulationStep float32         // seconds per update, the constant step of each simulation step
}

// Priority determines when the system will run relative to other systems, higher meaning sooner
func (s *ParticlePhysicsSystem) Priority() int {
	return -100
}

// Add adds a new entity to the system
func (s *ParticlePhysicsSystem) Add(entity particle) {
	s.ParticleEngine.Add(entity)
}

// Remove removes an entity from the system
func (s *ParticlePhysicsSystem) Remove(entity ecs.BasicEntity) {
	s.ParticleEngine.Remove(entity.ID())
}

// New is called every time the system is added to a world
func (s *ParticlePhysicsSystem) New(world *ecs.World) {
	s.simulationAcc = 0
	s.simulationStep = 1.0 / float32(s.SimulationRate)
}

// Update is called as often as the game loop can. If the accumulator has accumulated
// at least 1 step worth of time, then the physics engine is run for that amount. Then the accumulator is decremented
// by a step. This continues until there is less than a step left in the accumulator.
// If for some reason the physics simulation takes longer than a step, physics will be simulated a
// maximum number times, which will have an effect similar to time slowing down
func (s *ParticlePhysicsSystem) Update(dt float32) {
	s.simulationAcc += dt

	// Simulate physics in steps until we've caught up to real time or hit the limit
	// Any remainder less than the simulationStep can be interpolated by the renderer
	for i := 0; s.simulationAcc > s.simulationStep && i < MaxPhysicsIterations; i++ {
		defer metrics.Timed(metrics.Func("Engine.Total"))
		s.simulationAcc -= s.simulationStep
		s.ParticleEngine.Integrate(s.simulationStep)
		s.ParticleEngine.ResolveCollisions()
	}
}
