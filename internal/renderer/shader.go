package renderer

import (
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	vertexCode, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return nil, err
	}

	fragmentCode, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return nil, err
	}

	vertex := gl.CreateShader(gl.VERTEX_SHADER)
	csource, free := gl.Strs(string(vertexCode) + "\x00")
	gl.ShaderSource(vertex, 1, csource, nil)
	free()
	gl.CompileShader(vertex)

	fragment := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource, free = gl.Strs(string(fragmentCode) + "\x00")
	gl.ShaderSource(fragment, 1, csource, nil)
	free()
	gl.CompileShader(fragment)

	ID := gl.CreateProgram()
	gl.AttachShader(ID, vertex)
	gl.AttachShader(ID, fragment)
	gl.LinkProgram(ID)

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	return &Shader{ID: ID}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &value[0])
}

func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}
