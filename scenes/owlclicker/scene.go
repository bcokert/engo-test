package owlclicker

import (
	"fmt"
	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/metrics"
	"github.com/bcokert/engo-test/owls"
	"github.com/bcokert/engo-test/physics"
)

type Scene struct {
	Log logging.Logger
}

// Type returns an identifying string for this system, primarily to differentiate systems
func (s *Scene) Type() string {
	return "OwlClicker"
}

// Preload runs exactly once before Setup is run, and its results can be used within Setup
func (s *Scene) Preload() {
	engo.Files.Load("textures/owl.png")
}

// Setup adds the systems to the world and initializes everything for the game
func (s *Scene) Setup(world *ecs.World) {
	owlTexture, err := common.LoadedSprite("textures/owl.png")
	if err != nil {
		panic("Error loading texture: " + err.Error())
	}

	// Priority -1000
	world.AddSystem(&common.RenderSystem{})

	// Priority 100
	world.AddSystem(&common.MouseSystem{})

	// Priority 0
	world.AddSystem(&owls.OwlSystem{
		Log: s.Log,
	})
	world.AddSystem(&system{OwlTexture: owlTexture, Seed: 312, OwlInterval: 2, log: s.Log})

	// Priority -100
	world.AddSystem(&physics.ParticlePhysicsSystem{
		ParticleEngine: physics.NewParticleEngine(
			engo.Point{0, 150},
			0.99,
			[]physics.Wall{
				physics.Wall{engo.Point{0, engo.GameHeight()}, engo.Point{engo.GameWidth(), engo.GameHeight()}},
				physics.Wall{engo.Point{0, 0}, engo.Point{0, engo.GameHeight()}},
				physics.Wall{engo.Point{engo.GameWidth(), 0}, engo.Point{engo.GameWidth(), engo.GameHeight()}},
			},
			s.Log),
		SimulationRate: 60,
	})
}

// Exit is run right before closing the game
func (s *Scene) Exit() {
	now := time.Now().Local()
	path := fmt.Sprintf("functionmetrics/owlicker.%s.metrics", now.Format("2006-01-02-15-04-05"))
	err := metrics.Output(path)

	if err != nil {
		s.Log.Error("An error ocurred writing the metrics file", logging.F{"error": err, "path": path})
		return
	}

	s.Log.Info("Created metrics file", logging.F{"path": path})
}
