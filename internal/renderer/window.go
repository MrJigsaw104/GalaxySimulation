package renderer

import (
	"bhole/internal/physics"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	window      *glfw.Window
	Camera      *Camera
	shader      *Shader
	particleVAO uint32
}

func NewWindow(width, height int, title string) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}

	// Derinlik testi ve blend'i etkinleştir
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Arka plan rengini siyah yap
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	camera := NewCamera(
		mgl32.Vec3{200, 100, 200},
		mgl32.Vec3{0, 1, 0},
	)
	shader, err := NewShader("assets/shaders/vertex.glsl", "assets/shaders/fragment.glsl")
	if err != nil {
		return nil, err
	}

	w := &Window{
		window: window,
		Camera: camera,
		shader: shader,
	}
	w.setupParticleVAO()

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		camera.ProcessMouseMovement(float32(xpos), float32(ypos), true)
	})

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.Enable(gl.PROGRAM_POINT_SIZE)

	return w, nil
}

func (w *Window) setupParticleVAO() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	w.particleVAO = vao

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
}

func (w *Window) Draw(sim *physics.Simulation) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	w.shader.Use()
	w.shader.SetMat4("view", w.Camera.GetViewMatrix())
	w.shader.SetMat4("projection", w.Camera.GetProjectionMatrix())
	w.shader.SetVec3("centerPos", sim.CenterPos)

	w.drawParticles(sim.Particles)

	w.window.SwapBuffers()
	glfw.PollEvents()
}

func (w *Window) Destroy() {
	w.window.Destroy()
}

func (w *Window) SetKeyCallback(cb glfw.KeyCallback) {
	w.window.SetKeyCallback(cb)
}

func (w *Window) SetShouldClose(value bool) {
	w.window.SetShouldClose(value)
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *Window) drawParticles(particles []physics.Particle) {
	w.shader.Use()
	gl.PointSize(12.0)
	model := mgl32.Ident4()
	w.shader.SetMat4("model", model)
	w.shader.SetMat4("view", w.Camera.GetViewMatrix())

	gl.PointSize(8.0)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	data := make([]float32, len(particles)*6)
	for i, p := range particles {
		data[i*6] = p.Position.X()
		data[i*6+1] = p.Position.Y()
		data[i*6+2] = p.Position.Z()
		data[i*6+3] = p.Velocity.X()
		data[i*6+4] = p.Velocity.Y()
		data[i*6+5] = p.Velocity.Z()
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STREAM_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(w.particleVAO)
	gl.DrawArrays(gl.POINTS, 0, int32(len(particles)))

	gl.DeleteBuffers(1, &vbo)
}

func (w *Window) SetCursorPosCallback(cb glfw.CursorPosCallback) {
	w.window.SetCursorPosCallback(cb)
}

func (w *Window) SetInputMode(mode glfw.InputMode, value int) {
	w.window.SetInputMode(mode, value)
}
