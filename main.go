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

var cubePositions = []mgl.Vec3{
	{0.0, 0.0, 0.0},
	{2.0, 5.0, -15.0},
	{-1.5, -2.2, -2.5},
	{-3.8, -2.0, -12.3},
	{2.4, -0.4, -3.5},
	{-1.7, 3.0, -7.5},
	{1.3, -2.0, -2.5},
	{1.5, 2.0, -2.5},
	{1.5, 0.2, -1.5},
	{-1.3, 1.0, -1.5},
}

var camera = graphics.GetCamera()

var deltaTime = 0.0
var lastFrame = 0.0

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources\\shaders\\vertex\\shader.glsl")
	fragmentShaderSource := getTextFileContents("resources\\shaders\\fragment\\shader.glsl")

	camera.LastX = float64(cnf.Width) / 2.0
	camera.LastY = float64(cnf.Height) / 2.0

	runtime.LockOSThread()

	window := graphics.InitGlfw(cnf)
	defer glfw.Terminate()
	program := graphics.InitOpenGL(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(program)

	texture := graphics.MakeTexture("resources\\textures\\container.png")
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture\x00")), 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	vao := graphics.MakeVao(shape.Cube, program)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.126, 0.145, 0.2, 1.0)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(camera.MouseCallback)
	window.SetScrollCallback(camera.ScrollCallback)

	for !window.ShouldClose() {
		draw(vao, window, program, texture)
	}

	window.Destroy()
}

// loop over cells and tell them to draw
func draw(vao uint32, window *glfw.Window, program uint32, texture *uint32) {
	// per-frame time logic
	// --------------------
	currentFrame := glfw.GetTime()
	deltaTime = currentFrame - lastFrame
	lastFrame = currentFrame

	// input
	// -----
	processInput(window)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *texture)

	gl.UseProgram(program)

	//Transformation Matrices
	projection := mgl.Perspective(mgl.DegToRad(float32(camera.Fov)), float32(cnf.Width)/float32(cnf.Height), 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// camera/view transformation
	view := mgl.Ident4()
	view = view.Mul4(mgl.LookAtV(camera.CameraPos, camera.CameraPos.Add(camera.CameraFront), camera.CameraUp))
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	gl.BindVertexArray(vao)

	// Render a bunch of cubes
	for _, cube := range cubePositions {
		model := mgl.Ident4()
		model = model.Mul4(mgl.Translate3D(cube.X(), cube.Y(), cube.Z()))
		//angle := mgl.DegToRad(20.0 * float32(i))
		//model = model.Mul4(mgl.HomogRotate3D(angle, mgl.Vec3{1, 0.3, 0.5}))
		modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

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

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		log.Println("Escape key pressed")
		window.SetShouldClose(true)
	}

	cameraSpeed := float32(2.5 * deltaTime)
	// Forward
	if window.GetKey(glfw.KeyW) == glfw.Press {
		log.Println("W key pressed")
		camera.CameraPos = camera.CameraPos.Add(camera.CameraFront.Mul(cameraSpeed))
	}

	// Backward
	if window.GetKey(glfw.KeyS) == glfw.Press {
		log.Println("S key pressed")
		camera.CameraPos = camera.CameraPos.Sub(camera.CameraFront.Mul(cameraSpeed))
	}

	// Left
	if window.GetKey(glfw.KeyA) == glfw.Press {
		log.Println("A key pressed")
		camera.CameraPos = camera.CameraPos.Sub(camera.CameraFront.Cross(camera.CameraUp).Normalize().Mul(cameraSpeed))
	}

	// Right
	if window.GetKey(glfw.KeyD) == glfw.Press {
		log.Println("D key pressed")
		camera.CameraPos = camera.CameraPos.Add(camera.CameraFront.Cross(camera.CameraUp).Normalize().Mul(cameraSpeed))
	}
}
