package main

import (
	"github.com/PetrusJPrinsloo/learnopengl/config"
	"github.com/PetrusJPrinsloo/learnopengl/graphics"
	"github.com/PetrusJPrinsloo/learnopengl/scripting"
	"github.com/PetrusJPrinsloo/learnopengl/shape"
	mgl "github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"math"
	"runtime"

	"github.com/PetrusJPrinsloo/assimp-go"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var cnf *config.Config

var lightDirection = mgl.Vec3{-0.2, -1.0, -0.3}

//var cubePosition = mgl.Vec3{0.0, 0.0, 0.0}

var cubePositions = []mgl.Vec3{
	//{ 0.0,  0.0,  0.0},
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

var pointLightPositions = []mgl.Vec3{
	{0.7, 0.2, 2.0},
	{2.3, -3.3, -4.0},
	{-4.0, 2.0, -12.0},
	{0.0, 0.0, -3.0},
}

var camera = graphics.GetCamera()

var deltaTime = 0.0
var lastFrame = 0.0

func main() {
	scripting.HelloWorld()

	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources\\shaders\\vertex\\colors.glsl")
	fragmentShaderSource := getTextFileContents("resources\\shaders\\fragment\\colors.glsl")
	vertexShaderSourceLight := getTextFileContents("resources\\shaders\\vertex\\light_cube.glsl")
	fragmentShaderSourceLight := getTextFileContents("resources\\shaders\\fragment\\light_cube.glsl")

	camera.LastX = float64(cnf.Width) / 2.0
	camera.LastY = float64(cnf.Height) / 2.0

	runtime.LockOSThread()

	srcMeshes, err := assimp.ParseFile()

	window := graphics.InitGlfw(cnf)
	graphics.InitOpenGL()
	defer glfw.Terminate()
	objectShader := graphics.ShaderFactory(vertexShaderSource, fragmentShaderSource)
	lightShader := graphics.ShaderFactory(vertexShaderSourceLight, fragmentShaderSourceLight)
	objectShader.Use()

	vao, vbo := graphics.MakeObjectVao(shape.Cube, objectShader.Id)
	lightVao := graphics.MakeLightVao(shape.Cube, lightShader.Id, vbo)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.126, 0.145, 0.2, 1.0)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(camera.MouseCallback)
	window.SetScrollCallback(camera.ScrollCallback)

	diffuseMap := graphics.MakeTexture("resources\\textures\\container2.png")
	objectShader.SetInt("material.diffuse", 0)
	specularMap := graphics.MakeTexture("resources\\textures\\container2_specular.png")
	objectShader.SetInt("material.specular", 1)
	objectShader.SetFloat("light.constant", 1.0)
	objectShader.SetFloat("light.linear", 0.09)
	objectShader.SetFloat("light.quadratic", 0.032)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	for !window.ShouldClose() {
		draw(vao, lightVao, window, &objectShader, &lightShader, diffuseMap, specularMap)
	}

	window.Destroy()
}

// draw function called from application loop
func draw(vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader, texture uint32, specularMap uint32) {
	// per-frame time logic
	// --------------------
	currentFrame := glfw.GetTime()
	deltaTime = currentFrame - lastFrame
	lastFrame = currentFrame

	// input
	// -----
	processInput(window)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	objectShader.Use()
	objectShader.SetVec3("objectColor", mgl.Vec3{1.0, 0.5, 0.31})
	objectShader.SetVec3("lightColor", mgl.Vec3{3.0, 3.0, 3.0})
	/*
	   Here we set all the uniforms for the 5/6 types of lights we have. We have to set them manually and index
	   the proper PointLight struct in the array to set each uniform variable. This can be done more code-friendly
	   by defining light types as classes and set their values in there, or by using a more efficient uniform approach
	   by using 'Uniform buffer objects', but that is something we'll discuss in the 'Advanced GLSL' tutorial.
	*/
	// directional light
	objectShader.SetVec3("dirLight.direction", mgl.Vec3{-0.2, -1.0, -0.3})
	objectShader.SetVec3("dirLight.ambient", mgl.Vec3{0.05, 0.05, 0.05})
	objectShader.SetVec3("dirLight.diffuse", mgl.Vec3{0.4, 0.4, 0.4})
	objectShader.SetVec3("dirLight.specular", mgl.Vec3{0.5, 0.5, 0.5})
	// point light 1
	objectShader.SetVec3("pointLights[0].position", pointLightPositions[0])
	objectShader.SetVec3("pointLights[0].ambient", mgl.Vec3{0.05, 0.05, 0.05})
	objectShader.SetVec3("pointLights[0].diffuse", mgl.Vec3{0.8, 0.8, 0.8})
	objectShader.SetVec3("pointLights[0].specular", mgl.Vec3{.0, 1.0, 1.0})
	objectShader.SetFloat("pointLights[0].constant", 1.0)
	objectShader.SetFloat("pointLights[0].linear", 0.09)
	objectShader.SetFloat("pointLights[0].quadratic", 0.032)
	// point light 2
	objectShader.SetVec3("pointLights[1].position", pointLightPositions[1])
	objectShader.SetVec3("pointLights[1].ambient", mgl.Vec3{0.05, 0.05, 0.05})
	objectShader.SetVec3("pointLights[1].diffuse", mgl.Vec3{0.8, 0.8, 0.8})
	objectShader.SetVec3("pointLights[1].specular", mgl.Vec3{1.0, 1.0, 1.0})
	objectShader.SetFloat("pointLights[1].constant", 1.0)
	objectShader.SetFloat("pointLights[1].linear", 0.09)
	objectShader.SetFloat("pointLights[1].quadratic", 0.032)
	// point light 3
	objectShader.SetVec3("pointLights[2].position", pointLightPositions[2])
	objectShader.SetVec3("pointLights[2].ambient", mgl.Vec3{0.05, 0.05, 0.05})
	objectShader.SetVec3("pointLights[2].diffuse", mgl.Vec3{0.8, 0.8, 0.8})
	objectShader.SetVec3("pointLights[2].specular", mgl.Vec3{1.0, 1.0, 1.0})
	objectShader.SetFloat("pointLights[2].constant", 1.0)
	objectShader.SetFloat("pointLights[2].linear", 0.09)
	objectShader.SetFloat("pointLights[2].quadratic", 0.032)
	// point light 4
	objectShader.SetVec3("pointLights[3].position", pointLightPositions[3])
	objectShader.SetVec3("pointLights[3].ambient", mgl.Vec3{0.05, 0.05, 0.05})
	objectShader.SetVec3("pointLights[3].diffuse", mgl.Vec3{0.8, 0.8, 0.8})
	objectShader.SetVec3("pointLights[3].specular", mgl.Vec3{1.0, 1.0, 1.0})
	objectShader.SetFloat("pointLights[3].constant", 1.0)
	objectShader.SetFloat("pointLights[3].linear", 0.09)
	objectShader.SetFloat("pointLights[3].quadratic", 0.032)
	// spotLight
	objectShader.SetVec3("spotLight.position", camera.CameraPos)
	objectShader.SetVec3("spotLight.direction", camera.CameraFront)
	objectShader.SetVec3("spotLight.ambient", mgl.Vec3{0.0, 0.0, 0.0})
	objectShader.SetVec3("spotLight.diffuse", mgl.Vec3{1.0, 1.0, 1.0})
	objectShader.SetVec3("spotLight.specular", mgl.Vec3{1.0, 1.0, 1.0})
	objectShader.SetFloat("spotLight.constant", 1.0)
	objectShader.SetFloat("spotLight.linear", 0.09)
	objectShader.SetFloat("spotLight.quadratic", 0.032)
	objectShader.SetFloat("spotLight.outerCutOff", float32(math.Cos(float64(mgl.DegToRad(15.0)))))
	objectShader.SetFloat("spotLight.cutOff", float32(math.Cos(float64(mgl.DegToRad(12.5)))))

	// material properties
	objectShader.SetFloat("material.shininess", 32.0)

	//Transformation Matrices
	projection := mgl.Perspective(mgl.DegToRad(float32(camera.Fov)), float32(cnf.Width)/float32(cnf.Height), 0.1, 100.0)
	objectShader.SetMat4("projection", projection)

	// camera/view transformation
	view := mgl.Ident4()
	view = view.Mul4(mgl.LookAtV(camera.CameraPos, camera.CameraPos.Add(camera.CameraFront), camera.CameraUp))
	objectShader.SetMat4("view", view)
	//objectShader.SetVec3("lightPos", lightPosition)
	objectShader.SetVec3("viewPos", camera.CameraPos)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, specularMap)

	gl.BindVertexArray(vao)

	for _, cubePosition := range cubePositions {
		model := mgl.Ident4()
		model = model.Mul4(mgl.Translate3D(cubePosition.X(), cubePosition.Y(), cubePosition.Z()))
		objectShader.SetMat4("model", model)

		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

	//also draw the lamp object
	lightCubeShader.Use()
	lightCubeShader.SetMat4("projection", projection)
	lightCubeShader.SetMat4("view", view)
	lightCubeShader.SetVec3("color", mgl.Vec3{1.0, 1.0, 1.0})

	for _, pointLight := range pointLightPositions {
		model := mgl.Ident4()
		model = model.Mul4(mgl.Translate3D(pointLight.X(), pointLight.Y(), pointLight.Z()))
		model = model.Mul4(mgl.Scale3D(0.3, 0.3, 0.3)) // a smaller cube
		lightCubeShader.SetMat4("model", model)

		gl.BindVertexArray(lightVao)
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
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyI) == glfw.Press {
		log.Println("Information dump")
		log.Println("Camera Position: ", camera.CameraPos)
		log.Println("Camera Front: ", camera.CameraFront)
		log.Println("Camera Up: ", camera.CameraUp)
	}

	cameraSpeed := float32(2.5 * deltaTime)
	// Forward
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.CameraPos = camera.CameraPos.Add(camera.CameraFront.Mul(cameraSpeed))
	}

	// Backward
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.CameraPos = camera.CameraPos.Sub(camera.CameraFront.Mul(cameraSpeed))
	}

	// Left
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.CameraPos = camera.CameraPos.Sub(camera.CameraFront.Cross(camera.CameraUp).Normalize().Mul(cameraSpeed))
	}

	// Right
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.CameraPos = camera.CameraPos.Add(camera.CameraFront.Cross(camera.CameraUp).Normalize().Mul(cameraSpeed))
	}
}
