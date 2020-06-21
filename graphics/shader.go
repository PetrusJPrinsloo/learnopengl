package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"strings"
)

type Shader struct {
	Id uint32
}

func ShaderFactory(vertexShaderSource string, fragmentShaderSource string) *Shader {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	shader := Shader{Id: program}
	return &shader
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logMsg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logMsg))

		return 0, fmt.Errorf("failed to compile %v: %v", source, logMsg)
	}

	return shader, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.Id)
}

func (s *Shader) setBool(name string, value bool) {

	//fast convert bool to int32
	bitSetVar := int32(0)
	if value {
		bitSetVar = 1
	}

	gl.Uniform1i(gl.GetUniformLocation(s.Id, gl.Str(name+"\x00")), bitSetVar)
}

// Wrapping the ugly gl function calls
func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.Id, gl.Str(name)), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.Id, gl.Str(name)), value)
}

func (s *Shader) SetVec2(name string, value mgl.Vec2) {
	gl.Uniform3fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, &value[0])
}

func (s *Shader) SetVec3(name string, value mgl.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, &value[0])
}

func (s *Shader) SetVec4(name string, value mgl.Vec4) {
	gl.Uniform3fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, &value[0])
}

func (s *Shader) SetMat2(name string, value mgl.Mat2) {
	gl.UniformMatrix2fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, false, &value[0])
}

func (s *Shader) SetMat3(name string, value mgl.Mat3) {
	gl.UniformMatrix3fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, false, &value[0])
}

func (s *Shader) SetMat4(name string, value mgl.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.Id, gl.Str(name)), 1, false, &value[0])
}
