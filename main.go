package main

import (
	"github.com/PetrusJPrinsloo/learnopengl/config"
	"github.com/PetrusJPrinsloo/learnopengl/graphics"
	"github.com/PetrusJPrinsloo/learnopengl/shape"
	mgl "github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var cnf *config.Config

/*
 * Temp variables for testing.
 * If variables stay here too long or ar very useful, move to config
 */
var angle = 0.0
var previousTime float64

// end temp variables

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources\\shaders\\vertex\\shader.glsl")
	fragmentShaderSource := getTextFileContents("resources\\shaders\\fragment\\shader.glsl")

	runtime.LockOSThread()

	window := graphics.InitGlfw(cnf)
	previousTime = glfw.GetTime()
	defer glfw.Terminate()
	program := graphics.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(program)

	//Transformation Matrices
	projection := mgl.Perspective(mgl.DegToRad(45.0), float32(cnf.Width)/float32(cnf.Height), 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl.LookAtV(mgl.Vec3{3, 3, 3}, mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	texture := graphics.MakeTexture("resources\\textures\\container.png")
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture\x00")), 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	vao := graphics.MakeVao(shape.Cube, program)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.126, 0.145, 0.2, 1.0)

	for !window.ShouldClose() {
		draw(vao, window, program, texture, model, modelUniform)
	}
}

// loop over cells and tell them to draw
func draw(vao uint32, window *glfw.Window, program uint32, texture *uint32, model mgl.Mat4, modelUniform int32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	time := glfw.GetTime()
	elapsed := time - previousTime
	previousTime = time

	angle += elapsed
	model = mgl.HomogRotate3D(float32(angle), mgl.Vec3{0, 1, 0})

	gl.UseProgram(program)

	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

	// Maintenance
	window.SwapBuffers()
	glfw.PollEvents()
}

func getTextFileContents(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	text := string(content)
	return text
}
