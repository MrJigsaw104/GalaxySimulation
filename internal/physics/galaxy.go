package physics

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Galaxy struct {
	Position            mgl32.Vec3
	Mass                float32
	Radius              float32
	SchwarzschildRadius float32
	AccretionDiskColor  mgl32.Vec3
}

func NewGalaxy(x, y, z float32, mass float32) *Galaxy {
	c := float32(299792458)
	G := float32(6.67430e-11)
	schwarzschildRadius := 2 * G * mass / (c * c)

	return &Galaxy{
		Position:            mgl32.Vec3{x, y, z},
		Mass:                mass,
		Radius:              schwarzschildRadius * 1.5,
		SchwarzschildRadius: schwarzschildRadius,
		AccretionDiskColor:  mgl32.Vec3{1.0, 0.6, 0.0},
	}
}

func (bh *Galaxy) CalculateGravitationalForce(pos mgl32.Vec3) mgl32.Vec3 {
	direction := bh.Position.Sub(pos)
	distance := direction.Len()

	eventHorizonFactor := float32(1.0)
	if distance < bh.SchwarzschildRadius*3 {
		eventHorizonFactor = float32(math.Pow(float64(bh.SchwarzschildRadius/distance), 4))
	}

	G := float32(6.67430e-11)
	c := float32(299792458)

	relativisticFactor := float32(1.0 / math.Sqrt(1.0-math.Min(0.99, float64(2*G*bh.Mass/(c*c*distance)))))

	force := G * bh.Mass / (distance * distance)

	return direction.Normalize().Mul(force * relativisticFactor * eventHorizonFactor)
}
