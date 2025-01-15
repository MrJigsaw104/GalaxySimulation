package renderer

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraMovement int

const (
	FORWARD CameraMovement = iota
	BACKWARD
	LEFT
	RIGHT
)

type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed float32
	MouseSens     float32
	Zoom          float32

	LastX      float64
	LastY      float64
	FirstMouse bool
}

func NewCamera(position, worldUp mgl32.Vec3) *Camera {
	camera := &Camera{
		Position:      mgl32.Vec3{0, 200, 400},
		WorldUp:       worldUp,
		Yaw:           -90.0,
		Pitch:         -30.0,
		MovementSpeed: 50.0,
		MouseSens:     0.1,
		Zoom:          45.0,
		FirstMouse:    true,
	}

	camera.updateCameraVectors()
	return camera
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(c.Zoom), 1920.0/1080.0, 0.1, 1000.0)
}

func (c *Camera) ProcessKeyboard(direction CameraMovement, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime

	switch direction {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}

	minDistance := float32(100.0)
	if c.Position.Len() < minDistance {
		c.Position = c.Position.Normalize().Mul(minDistance)
	}
}

func (c *Camera) updateCameraVectors() {
	direction := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
	}
	c.Front = direction.Normalize()

	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

func (c *Camera) ProcessMouseMovement(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= c.MouseSens
	yoffset *= c.MouseSens

	c.Yaw += xoffset
	c.Pitch += yoffset

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	c.updateCameraVectors()
}
