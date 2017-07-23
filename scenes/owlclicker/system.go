package owlclicker

import (
	"math/rand"

	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/owls"
	"github.com/bcokert/engo-test/physics"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

// Every scene has 1 game system that manages the actual scene and any concrete entities that the other systems manage
type system struct {
	Seed          int64
	OwlTexture    *common.Texture
	OwlInterval   float32
	log           logging.Logger
	timeToNextOwl float32
	entities      map[uint64]*owl
	rand          *rand.Rand
	world         *ecs.World
}

// Remove removes a entity from the system, by its entity id
func (s *system) Remove(entity ecs.BasicEntity) {
	if _, ok := s.entities[entity.ID()]; ok {
		delete(s.entities, entity.ID())
	}
}

// New is called every time the system is added to a world
func (s *system) New(world *ecs.World) {
	s.world = world
	if s.OwlInterval == 0 {
		s.OwlInterval = 5
	}
	s.timeToNextOwl = 0

	if s.rand == nil {
		s.rand = rand.New(rand.NewSource(s.Seed))
	}

	if s.entities == nil {
		s.entities = make(map[uint64]*owl, 10)
	}
}

// Update will periodically add new owls to the game
func (s *system) Update(dt float32) {
	s.timeToNextOwl -= dt
	if s.timeToNextOwl <= 0 {
		s.timeToNextOwl = s.OwlInterval

		scale := s.rand.Float32()/2 + 0.25
		maxWidth := int(engo.GameWidth() - s.OwlTexture.Width()*scale)
		maxHeight := int(engo.GameHeight() - s.OwlTexture.Height()*scale)
		position := engo.Point{
			float32(s.rand.Intn(maxWidth-int(s.OwlTexture.Width()))) + s.OwlTexture.Width(),
			float32(s.rand.Intn(maxHeight-int(s.OwlTexture.Height()))) + s.OwlTexture.Height(),
		}
		velocity := engo.Point{s.rand.Float32()*200 - 100, s.rand.Float32()*100 - 50}

		owl := newOwl(position, velocity, s.OwlTexture, scale, 1, float32(s.rand.Intn(5)+2))
		s.entities[owl.BasicEntity().ID()] = owl

		for _, worldSystem := range s.world.Systems() {
			switch targetSystem := worldSystem.(type) {
			case *common.RenderSystem:
				targetSystem.Add(owl.BasicEntity(), owl.RenderComponent(), owl.SpaceComponent())
			case *common.MouseSystem:
				targetSystem.Add(owl.BasicEntity(), owl.MouseComponent(), owl.SpaceComponent(), nil)
			case *owls.OwlSystem:
				targetSystem.Add(owl)
			case *physics.ParticlePhysicsSystem:
				targetSystem.Add(owl)
			}
		}
	}
}
