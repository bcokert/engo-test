package owls

import (
	"image/color"

	"github.com/bcokert/engo-test/logging"

	"engo.io/ecs"
	"engo.io/engo/common"
)

type BasicHealthComponent struct {
	Health    float64
	MaxHealth float64
}

type owlEntity interface {
	BasicEntity() *ecs.BasicEntity
	RenderComponent() *common.RenderComponent
	MouseComponent() *common.MouseComponent
	BasicHealthComponent() *BasicHealthComponent
}

// The OwlSystem manages a group of owls; their creation, mouse interaction, and so on.
// Everything but physics related properties are managed by the Owl System.
// An owl system also automatically adds owls to itself over time
type OwlSystem struct {
	entities map[uint64]owlEntity
	world    *ecs.World
	Log      logging.Logger
}

// Add adds a new entity to the system
func (s *OwlSystem) Add(entity owlEntity) {
	s.entities[entity.BasicEntity().ID()] = entity
}

// Remove removes an entity from the system, by its entity id
func (s *OwlSystem) Remove(entity ecs.BasicEntity) {
	if _, ok := s.entities[entity.ID()]; ok {
		delete(s.entities, entity.ID())
	}
}

// New is called every time the system is added to a world
func (s *OwlSystem) New(world *ecs.World) {
	s.world = world

	if s.entities == nil {
		s.entities = make(map[uint64]owlEntity, 10)
	}
}

// Update processes the user interactions with owls, updating or removing them as necessary
// It also internally manages the list of owls, adding them at certain intervals
func (s *OwlSystem) Update(dt float32) {
	for _, owl := range s.entities {
		health := owl.BasicHealthComponent()
		mouse := owl.MouseComponent()

		// update health, removing if dead
		if mouse.Clicked {
			health.Health--
			if health.Health <= 0 {
				s.world.RemoveEntity(*owl.BasicEntity())
				continue
			}
		}

		// calculate the color based on health, redder is deader
		percentHealthy := health.Health / health.MaxHealth
		col := color.RGBA{200, uint8(200.0 * percentHealthy), uint8(200.0 * percentHealthy), 255}

		// if the mouse is over the owl, make it glow whiter
		if mouse.Hovered {
			col.R += 55
			col.G += 55
			col.B += 55
		}

		owl.RenderComponent().Color = col
	}
}
