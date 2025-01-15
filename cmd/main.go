package main

import (
	"bhole/internal/physics"
	"bhole/internal/renderer"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	window, err := renderer.NewWindow(1920, 1080, "Galaxy simulator by 0xE5")
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	centerPos := mgl32.Vec3{0, 0, 0}
	simulation := physics.NewSimulation(centerPos)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}

		if action == glfw.Press || action == glfw.Repeat {
			switch key {
			case glfw.KeyW:
				window.Camera.ProcessKeyboard(renderer.FORWARD, 0.016)
			case glfw.KeyS:
				window.Camera.ProcessKeyboard(renderer.BACKWARD, 0.016)
			case glfw.KeyA:
				window.Camera.ProcessKeyboard(renderer.LEFT, 0.016)
			case glfw.KeyD:
				window.Camera.ProcessKeyboard(renderer.RIGHT, 0.016)
			}
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		if window.Camera.FirstMouse {
			window.Camera.LastX = xpos
			window.Camera.LastY = ypos
			window.Camera.FirstMouse = false
		}

		xoffset := float32(xpos - window.Camera.LastX)
		yoffset := float32(window.Camera.LastY - ypos)

		window.Camera.LastX = xpos
		window.Camera.LastY = ypos

		window.Camera.ProcessMouseMovement(xoffset, yoffset, true)
	})

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	for !window.ShouldClose() {
		simulation.Update()
		window.Draw(simulation)
	}
}
