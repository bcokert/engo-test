package main

import (
	"engo.io/ecs"
	"engo.io/engo/common"
)

type Owl struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}
