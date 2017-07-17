package physics

import "engo.io/engo"

type ParticleCollisionManifold struct {
	a                particle   // the primary object in the collision
	b                particle   // can be nil for constraint-based collisions, like walls
	penetrationDepth float32    // how much along the contactNormal we have penetrated
	contactNormal    engo.Point // towards a
}

func (e *ParticleEngine) ResolveCollisions() {
}

func (e *ParticleEngine) detectCollisions() []*ParticleCollisionManifold {
	return []*ParticleCollisionManifold{}
}

// width := body.SpaceComponent.Width
// height := body.SpaceComponent.Height
// left := body.SpaceComponent.Position.X
// right := body.SpaceComponent.Position.X + width
// top := body.SpaceComponent.Position.Y
// bottom := body.SpaceComponent.Position.Y + height

// if p.X+dp.X <= s.BoxLeft {
// 	p.X = 0
// 	dp.X = dp.X * -1 * (1 - s.Elasticity.X)
// 	dp.Y *= (1 - s.Friction.Y)
// } else if p.X+dp.X+w >= s.BoxRight {
// 	p.X = s.BoxRight - w
// 	dp.X = dp.X * -1 * (1 - s.Elasticity.X)
// 	dp.Y *= (1 - s.Friction.Y)
// } else {
// 	p.X = p.X + dp.X
// }
// dp.X = dp.X + ddp.X + s.Gravity.X
// if math.Abs(float64(dp.X)) < 0.001 {
// 	dp.X = 0
// }

// accelY := ddp.Y + s.Gravity.Y
// if entity.BasicEntity().ID() == 1 {
// 	log.Debug("Tick: y: %v   dy: %v   ddy: %v", p.Y, dp.Y, accelY)
// }
// if p.Y+dp.Y <= s.BoxTop {
// 	p.Y = 0
// 	dp.Y = dp.Y * -1 * (1 - s.Elasticity.Y)
// 	if dp.X != 0 {
// 		dp.X *= (1 - s.Friction.X)
// 	}
// } else if p.Y+dp.Y+h >= s.BoxBottom {
// 	if entity.BasicEntity().ID() == 1 {
// 		log.Printf("Collision with floor")
// 	}
// 	p.Y = s.BoxBottom - h
// 	dp.Y = dp.Y * -1 * (1 - s.Elasticity.Y)
// 	if entity.BasicEntity().ID() == 1 {
// 		log.Debug("After adjusting y and dp.Y: y: %v   dy: %v", p.Y, dp.Y)
// 	}
// 	if accelY > dp.Y*-1 {
// 		// the entity did not bounce off fast enough to overcome acceleration
// 		dp.Y = 0
// 	}
// 	if dp.X != 0 {
// 		dp.X *= (1 - s.Friction.X)
// 	}
// } else {
// 	p.Y = p.Y + dp.Y
// 	dp.Y += accelY
// }
// if math.Abs(float64(dp.Y)) < 0.001 {
// 	dp.Y = 0
// }
