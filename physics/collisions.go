package physics

import (
	"engo.io/engo"
	"github.com/bcokert/engo-test/logging"
	"github.com/bcokert/engo-test/metrics"
)

type Wall struct {
	P1 engo.Point
	P2 engo.Point
}

type ParticleCollisionManifold struct {
	a                particle   // the primary object in the collision
	b                particle   // can be nil for constraint-based collisions, like walls
	penetrationDepth float32    // how much along the contactNormal we have penetrated
	contactNormal    engo.Point // towards a
}

func (e *ParticleEngine) ResolveCollisions() {
	collisions := e.detectCollisions()

	defer metrics.Timed(metrics.Func("Engine.ResolveCollisions"))
	for _, collision := range collisions {
		restitution := collision.a.ParticleComponent().Restitution
		if collision.b != nil {
			if collision.b.ParticleComponent().Restitution < restitution {
				restitution = collision.b.ParticleComponent().Restitution
			}
		}

		totalVelocity := collision.a.ParticleComponent().Velocity
		if collision.b != nil {
			totalVelocity.Subtract(collision.b.ParticleComponent().Velocity)
		}
		separatingVelocity := engo.DotProduct(totalVelocity, collision.contactNormal)

		// if the objects are moving away from eachother already, don't resolve collision
		if separatingVelocity <= 0 {
			continue
		}

		deltaVelocity := separatingVelocity*-restitution - separatingVelocity

		// divide the deltaVelocity amongst the participants so that each gets velocity
		// inversely proportional to its mass.
		// aka the lighter objects gets more velocity change
		// deltaVx = (minvx/minvtotal) * deltaV * normal
		// if there is only 1 object, then minvx / minvtotal == 1, which simplifies the calculation
		if collision.b == nil {
			deltaVa := collision.contactNormal
			deltaVa.MultiplyScalar(deltaVelocity)

			collision.a.ParticleComponent().Velocity.Add(deltaVa)
		} else {
			totalInverseMass := collision.a.ParticleComponent().InvMass + collision.b.ParticleComponent().InvMass
			percentA := collision.a.ParticleComponent().InvMass / totalInverseMass
			percentB := collision.b.ParticleComponent().InvMass / totalInverseMass

			// direction first
			deltaVa := collision.contactNormal
			deltaVb := collision.contactNormal
			deltaVb.MultiplyScalar(-1)

			// then magnitudes
			deltaVa.MultiplyScalar(deltaVelocity * percentA)
			deltaVb.MultiplyScalar(deltaVelocity * percentB)

			collision.a.ParticleComponent().Velocity.Add(deltaVa)
			collision.b.ParticleComponent().Velocity.Add(deltaVb)
		}
	}
}

func (e *ParticleEngine) detectCollisions() []*ParticleCollisionManifold {
	defer metrics.Timed(metrics.Func("Engine.detectCollisions"))
	collisions := make([]*ParticleCollisionManifold, 0, len(e.ParticleRegistry.particles))

	// Detect collisions with walls
	for _, p := range e.ParticleRegistry.particles {
		body := p.ParticleComponent()
		width := body.SpaceComponent.Width
		height := body.SpaceComponent.Height
		r := width
		if height > width {
			r = height
		}
		r = r / 2

		for _, wall := range e.walls {
			// Set P to the center of the entity
			P := body.SpaceComponent.Position
			P.X += width / 2
			P.Y += height / 2

			L := wall.P1         // L = P1
			L.Subtract(wall.P2)  // L = P1 - P2
			L, _ = L.Normalize() // L = (P1 - P2).Normalize()
			PL := P              // PL = P
			PL.Subtract(wall.P2) // PL = P - P2

			PLprojL := L
			PLprojL.MultiplyScalar(engo.DotProduct(PL, L))
			nearestPoint := wall.P2
			nearestPoint.Add(PLprojL)
			wallToP := P
			wallToP.Subtract(nearestPoint)

			normal, distanceToWall := wallToP.Normalize()
			normal.MultiplyScalar(-1) // since we always give the normal from a's perspective (aka a thinks b hit it)

			if distanceToWall <= r {
				// We're always going to treat it as if only the edge of the circle collided; thus we bounce off the wall
				// reflected over the normal, and the collision point is the same as the line from the sphere center to the wall
				manifold := ParticleCollisionManifold{
					a:                p,
					b:                nil,                // the wall will not react to the collision
					penetrationDepth: r - distanceToWall, // how much of the sphere is on the other side of the wall
					contactNormal:    normal,
				}
				collisions = append(collisions, &manifold)

				e.log.Debug("Detected collision", logging.F{"wall": wall, "manifold": manifold})
			}
		}
	}

	return collisions
}
