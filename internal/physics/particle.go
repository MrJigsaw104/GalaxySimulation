package physics

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Particle struct {
	Position mgl32.Vec3
	Velocity mgl32.Vec3
	Mass     float32
}

type Simulation struct {
	CenterPos mgl32.Vec3
	Particles []Particle
}

func NewSimulation(centerPos mgl32.Vec3) *Simulation {
	rand.Seed(time.Now().UnixNano())

	particles := make([]Particle, 200000)
	for i := range particles {
		// Spiral başlangıç pozisyonları
		angle := float64(rand.Float32() * 2 * math.Pi)
		radius := float64(rand.Float32()*400 + 50)

		// Yükseklik dağılımı (disk şekli için)
		height := float32(rand.NormFloat64() * 10)

		particles[i] = Particle{
			Position: mgl32.Vec3{
				float32(radius * math.Cos(angle)),
				height,
				float32(radius * math.Sin(angle)),
			},
			Velocity: mgl32.Vec3{
				float32(-math.Sin(angle)),
				0,
				float32(math.Cos(angle)),
			}.Mul(float32(math.Sqrt(radius)) * 0.5),
			Mass: 1.0,
		}
	}
	return &Simulation{CenterPos: centerPos, Particles: particles}
}

func (s *Simulation) Update() {
	dt := float32(0.016)

	for i := range s.Particles {
		toCenter := s.CenterPos.Sub(s.Particles[i].Position)
		distance := toCenter.Len()

		G := float32(9.8)
		force := toCenter.Normalize().Mul(G * 1000 / (distance * distance))

		currentSpeed := s.Particles[i].Velocity.Len()
		normalizedVel := s.Particles[i].Velocity.Normalize()
		tangentialDir := mgl32.Vec3{-normalizedVel.Z(), 0, normalizedVel.X()}

		tangentialForce := tangentialDir.Mul(currentSpeed * currentSpeed / distance)

		totalForce := force.Add(tangentialForce)
		acceleration := totalForce.Mul(1.0 / s.Particles[i].Mass)

		s.Particles[i].Velocity = s.Particles[i].Velocity.Add(acceleration.Mul(dt))
		s.Particles[i].Position = s.Particles[i].Position.Add(s.Particles[i].Velocity.Mul(dt))

		minDist := float32(50.0)
		if distance < minDist {
			s.Particles[i].Position = s.CenterPos.Add(toCenter.Normalize().Mul(minDist))
		}
	}
}
