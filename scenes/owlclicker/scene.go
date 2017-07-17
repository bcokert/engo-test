package owlclicker

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/systems"
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

	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	world.AddSystem(&systems.ParticlePhysicsSystem{
		TopLeft:        engo.Point{0, 0},
		BottomRight:    engo.Point{engo.GameWidth(), engo.GameHeight()},
		Gravity:        engo.Point{0, 2},
		DampingFactor:  0.99,
		SimulationRate: 60,
		Log:            s.Log,
	})
	world.AddSystem(&systems.OwlSystem{
		Log: s.Log,
	})
	world.AddSystem(&system{OwlTexture: owlTexture, Seed: 312, OwlInterval: 3})
}
