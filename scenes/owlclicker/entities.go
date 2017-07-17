package owlclicker

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/bcokert/engo-test/systems"
)

type owl struct {
	basicEntity          ecs.BasicEntity
	renderComponent      common.RenderComponent
	mouseComponent       common.MouseComponent
	particleComponent    systems.ParticleComponent
	basicHealthComponent systems.BasicHealthComponent
}

func (o *owl) BasicEntity() *ecs.BasicEntity {
	return &o.basicEntity
}

func (o *owl) RenderComponent() *common.RenderComponent {
	return &o.renderComponent
}

func (o *owl) MouseComponent() *common.MouseComponent {
	return &o.mouseComponent
}

func (o *owl) ParticleComponent() *systems.ParticleComponent {
	return &o.particleComponent
}

func (o *owl) SpaceComponent() *common.SpaceComponent {
	return &o.particleComponent.SpaceComponent
}

func (o *owl) BasicHealthComponent() *systems.BasicHealthComponent {
	return &o.basicHealthComponent
}

func newOwl(position, velocity engo.Point, texture *common.Texture, scale, mass float32, health float64) *owl {
	return &owl{
		basicEntity: ecs.NewBasic(),
		renderComponent: common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{scale, scale},
		},
		mouseComponent: common.MouseComponent{},
		particleComponent: systems.NewParticleComponent(
			texture.Width()*scale,
			texture.Height()*scale,
			mass,
			position,
			velocity,
		),
		basicHealthComponent: systems.BasicHealthComponent{
			Health:    health,
			MaxHealth: health,
		},
	}
}
