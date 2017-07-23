package owlclicker

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/bcokert/engo-test/owls"
	"github.com/bcokert/engo-test/physics"
)

type owl struct {
	basicEntity          ecs.BasicEntity
	renderComponent      common.RenderComponent
	mouseComponent       common.MouseComponent
	particleComponent    physics.ParticleComponent
	basicHealthComponent owls.BasicHealthComponent
	healthBarComponent   owls.HealthBarComponent
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

func (o *owl) ParticleComponent() *physics.ParticleComponent {
	return &o.particleComponent
}

func (o *owl) SpaceComponent() *common.SpaceComponent {
	return &o.particleComponent.SpaceComponent
}

func (o *owl) BasicHealthComponent() *owls.BasicHealthComponent {
	return &o.basicHealthComponent
}

func (o *owl) HealthBarComponent() *owls.HealthBarComponent {
	return &o.healthBarComponent
}

func newOwl(position, velocity engo.Point, texture *common.Texture, scale, mass, health float32) *owl {
	return &owl{
		basicEntity: ecs.NewBasic(),
		renderComponent: common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{scale, scale},
		},
		mouseComponent: common.MouseComponent{},
		particleComponent: physics.NewParticleComponent(
			texture.Width()*scale,
			texture.Height()*scale,
			mass,
			position,
			velocity,
		),
		basicHealthComponent: owls.BasicHealthComponent{
			Health:    health,
			MaxHealth: health,
		},
		healthBarComponent: owls.NewHealthBarComponent(texture.Width()*scale, 6, position),
	}
}
