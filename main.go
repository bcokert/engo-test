package main

import (
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type MainScene struct{}

func (s *MainScene) Type() string {
	return "SceneType"
}

func (s *MainScene) Preload() {
	engo.Files.Load("textures/owl.png")
}

func (s *MainScene) Setup(world *ecs.World) {
	world.AddSystem(&common.RenderSystem{})

	owlTexture, err := common.LoadedSprite("textures/owl.png")
	if err != nil {
		log.Println("Error loading texture: ", err)
		return
	}

	owls := make([]*Owl, 0, 10)
	for i := 0; i < 10; i++ {
		owl := Owl{
			BasicEntity: ecs.NewBasic(),
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{float32(i*100%400 + 20), float32((i / 4) * 130)},
				Width:    99,
				Height:   128,
			},
			RenderComponent: common.RenderComponent{
				Drawable: owlTexture,
				Scale:    engo.Point{1, 1},
			},
		}

		owls = append(owls, &owl)
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, owl := range owls {
				sys.Add(&owl.BasicEntity, &owl.RenderComponent, &owl.SpaceComponent)
			}
		}
	}
}

func main() {
	options := engo.RunOptions{
		Title:  "MyGame",
		Width:  500,
		Height: 500,
	}

	engo.Run(options, &MainScene{})
}
