package physics

import (
	"engo.io/engo"
	"github.com/bcokert/engo-test/logging"
)

type ParticleEngine struct {
	gravity           engo.Point     // the gravitational force, in meters/second^2
	dampingFactor     float32        // the damping factor that reduces the velocity slightly (sort of like an aerodynamic impulse due to drag) but is primarily for numerical stability
	log               logging.Logger // engine wide logger
	*ParticleRegistry                // stores all particles and efficiently finds them for the collision detector
}

func NewParticleEngine(gravity engo.Point, dampingFactor float32, logger logging.Logger) *ParticleEngine {
	if logger == nil {
		logger = logging.NewDefaultLogger(logging.INFO)
	}

	if dampingFactor < 0 || dampingFactor > 1 {
		dampingFactor = 0.99
	}

	logger.Info("Creating new ParticleEngine with safe configuration", logging.F{"gravity": gravity, "dampingFactor": dampingFactor})

	return &ParticleEngine{
		gravity:       gravity,
		dampingFactor: dampingFactor,
		log:           logger,
		ParticleRegistry: &ParticleRegistry{
			log:       logger,
			particles: map[uint64]particle{},
		},
	}
}
