package owls

import (
	"image/color"

	"github.com/bcokert/engo-test/logging"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type BasicHealthComponent struct {
	Health    float32
	MaxHealth float32
}

type HealthBarComponent struct {
	width       float32
	height      float32
	position    engo.Point
	emptySpace  common.SpaceComponent
	emptyRender common.RenderComponent
	emptyBasic  ecs.BasicEntity
	fullSpace   common.SpaceComponent
	fullRender  common.RenderComponent
	fullBasic   ecs.BasicEntity
}

func (c *HealthBarComponent) Update(percent float32, parentPosition engo.Point) {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	position := engo.Point{parentPosition.X, parentPosition.Y - c.height*2}

	c.fullSpace.Width = percent * c.width
	c.emptySpace.Width = c.width - c.fullSpace.Width

	c.emptySpace.Position = position
	c.fullSpace.Position = position
	c.fullSpace.Position.X += c.emptySpace.Width
}

func NewHealthBarComponent(width, height float32, parentPosition engo.Point) HealthBarComponent {
	position := engo.Point{parentPosition.X, parentPosition.Y - height*2}
	return HealthBarComponent{
		width:       width,
		height:      height,
		position:    position,
		emptySpace:  common.SpaceComponent{Position: position, Width: 0, Height: height},
		emptyRender: common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 0, 0, 255}},
		emptyBasic:  ecs.NewBasic(),
		fullSpace:   common.SpaceComponent{Position: position, Width: width, Height: height},
		fullRender:  common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{0, 255, 0, 255}},
		fullBasic:   ecs.NewBasic(),
	}
}

type owlEntity interface {
	BasicEntity() *ecs.BasicEntity
	RenderComponent() *common.RenderComponent
	MouseComponent() *common.MouseComponent
	BasicHealthComponent() *BasicHealthComponent
	HealthBarComponent() *HealthBarComponent
	SpaceComponent() *common.SpaceComponent
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
	healthbar := entity.HealthBarComponent()

	for _, worldSystem := range s.world.Systems() {
		switch targetSystem := worldSystem.(type) {
		case *common.RenderSystem:
			targetSystem.Add(&healthbar.emptyBasic, &healthbar.emptyRender, &healthbar.emptySpace)
			targetSystem.Add(&healthbar.fullBasic, &healthbar.fullRender, &healthbar.fullSpace)
		}
	}
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
		healthbar := owl.HealthBarComponent()
		mouse := owl.MouseComponent()

		// update health, removing if dead
		if mouse.Clicked {
			health.Health--
			if health.Health <= 0 {
				s.world.RemoveEntity(*owl.BasicEntity())
				s.world.RemoveEntity(owl.HealthBarComponent().emptyBasic)
				s.world.RemoveEntity(owl.HealthBarComponent().fullBasic)
				continue
			}
		}

		// update healthbar based on health
		percentHealthy := health.Health / health.MaxHealth
		healthbar.Update(percentHealthy, owl.SpaceComponent().Position)

		col := color.RGBA{255, 255, 255, 255}

		// if the mouse is over the owl, make it glow a bit blue
		if mouse.Hovered {
			col.R -= 40
			col.G -= 40
		}

		owl.RenderComponent().Color = col
	}
}
