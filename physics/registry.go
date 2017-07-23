package physics

import (
	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/metrics"
)

// The ParticleRegistry stores all particles in some structure, and retrieves them for collision related logic
// This is where optimizations like quadtrees and other subdivisions would be done to rule out
// particles when checking for collisions
type ParticleRegistry struct {
	log       logging.Logger
	particles map[uint64]particle
}

func (r *ParticleRegistry) Add(p particle) {
	defer metrics.Timed(metrics.Func("Registry.Add"))
	r.log.Debug("Adding particle to registry", logging.F{"id": p.BasicEntity().ID(), "particleComponent": p.ParticleComponent()})
	r.particles[p.BasicEntity().ID()] = p
}

func (r *ParticleRegistry) Remove(id uint64) {
	defer metrics.Timed(metrics.Func("Registry.Remove"))
	p, ok := r.particles[id]
	if ok {
		r.log.Debug("Removing particle from registry", logging.F{"id": id, "particleComponent": p.ParticleComponent()})
		delete(r.particles, id)
	}
}
