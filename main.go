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

var lightPosition = mgl.Vec3{1.2, 1.0, 2.0}

var cubePosition = mgl.Vec3{0.0, 0.0, 0.0}

var camera = graphics.GetCamera()

var deltaTime = 0.0
var lastFrame = 0.0

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources\\shaders\\vertex\\colors.glsl")
	fragmentShaderSource := getTextFileContents("resources\\shaders\\fragment\\colors.glsl")
	vertexShaderSource_light := getTextFileContents("resources\\shaders\\vertex\\light_cube.glsl")
	fragmentShaderSource_light := getTextFileContents("resources\\shaders\\fragment\\light_cube.glsl")

	camera.LastX = float64(cnf.Width) / 2.0
	camera.LastY = float64(cnf.Height) / 2.0

	runtime.LockOSThread()

	window := graphics.InitGlfw(cnf)
	graphics.InitOpenGL()
	defer glfw.Terminate()
	objectShader := graphics.ShaderFactory(vertexShaderSource, fragmentShaderSource)
	lightShader := graphics.ShaderFactory(vertexShaderSource_light, fragmentShaderSource_light)
	objectShader.Use()

	texture := graphics.MakeTexture("resources\\textures\\container.png")
	objectShader.SetInt("texture", 0)
	gl.BindFragDataLocation(objectShader.Id, 0, gl.Str("outputColor\x00"))

	vao, vbo := graphics.MakeObjectVao(shape.Cube, objectShader.Id)
	lightVao := graphics.MakeLightVao(shape.Cube, lightShader.Id, vbo)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.126, 0.145, 0.2, 1.0)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(camera.MouseCallback)
	window.SetScrollCallback(camera.ScrollCallback)

	for !window.ShouldClose() {
		draw(vao, lightVao, window, &objectShader, &lightShader, texture)
	}

	window.Destroy()
}

// draw function called from application loop
func draw(vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader, texture *uint32) {
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

	objectShader.Use()
	objectShader.SetVec3("objectColor", mgl.Vec3{1.0, 0.5, 0.31})
	objectShader.SetVec3("lightColor", mgl.Vec3{1.0, 1.0, 1.0})

	//Transformation Matrices
	projection := mgl.Perspective(mgl.DegToRad(float32(camera.Fov)), float32(cnf.Width)/float32(cnf.Height), 0.1, 100.0)
	objectShader.SetMat4("projection", projection)

	// camera/view transformation
	view := mgl.Ident4()
	view = view.Mul4(mgl.LookAtV(camera.CameraPos, camera.CameraPos.Add(camera.CameraFront), camera.CameraUp))
	objectShader.SetMat4("view", view)

	gl.BindVertexArray(vao)

	// Render a bunch of cubes
	model := mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(cubePosition.X(), cubePosition.Y(), cubePosition.Z()))
	objectShader.SetMat4("model", model)

	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	// also draw the lamp object
	lightCubeShader.Use()
	lightCubeShader.SetMat4("projection", projection)
	lightCubeShader.SetMat4("view", view)
	model = mgl.Ident4()
	model = model.Mul4(mgl.Translate3D(lightPosition.X(), lightPosition.Y(), lightPosition.Z()))
	model = model.Mul4(mgl.Scale3D(0.3, 0.3, 0.3)) // a smaller cube
	lightCubeShader.SetMat4("model", model)

	gl.BindVertexArray(lightVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

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
