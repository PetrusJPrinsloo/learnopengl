package main

import (
	"fmt"
	"github.com/PetrusJPrinsloo/learnopengl/config"
	"github.com/PetrusJPrinsloo/learnopengl/graphics"
	"github.com/PetrusJPrinsloo/learnopengl/shape"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v2"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var cnf *config.Config

var cubes = []mgl.Vec3{
	{0, 0, 0},
	{1, 0, 0},
	{2, 0, 0},
	{3, 0, 0},
	{4, 0, 0},
}

var camera = graphics.GetCamera()

var deltaTime = 0.0
var lastFrame = 0.0

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources/shaders/vertex/colors.glsl")
	fragmentShaderSource := getTextFileContents("resources/shaders/fragment/colors.glsl")
	vertexShaderSourceLight := getTextFileContents("resources/shaders/vertex/light_cube.glsl")
	fragmentShaderSourceLight := getTextFileContents("resources/shaders/fragment/light_cube.glsl")

	camera.LastX = float64(cnf.Width) / 2.0
	camera.LastY = float64(cnf.Height) / 2.0

	runtime.LockOSThread()

	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()
	GLFW, err := graphics.InitGlfw(io, cnf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}

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

	GLFW.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	GLFW.Window.SetCursorPosCallback(camera.MouseCallback)
	GLFW.Window.SetScrollCallback(camera.ScrollCallback)

	objectShader.SetInt("material.specular", 1)
	objectShader.SetFloat("light.constant", 1.0)
	objectShader.SetFloat("light.linear", 0.09)
	objectShader.SetFloat("light.quadratic", 0.032)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	renderer, err := graphics.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	defer GLFW.Dispose()
	PrintMemUsage()
	Run(GLFW, renderer, vao, lightVao, GLFW.Window, &objectShader, &lightShader) //, diffuseMap, specularMap)
}

// draw function called from application loop
//func draw(vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader, texture uint32, specularMap uint32) {
func draw(vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader) {
	// per-frame time logic
	// --------------------
	currentFrame := glfw.GetTime()
	deltaTime = currentFrame - lastFrame
	lastFrame = currentFrame

	// input
	// -----
	processInput(window)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//Transformation Matrices
	projection := mgl.Perspective(mgl.DegToRad(float32(camera.Fov)), float32(cnf.Width)/float32(cnf.Height), 0.1, 100.0)

	// camera/view transformation
	view := mgl.Ident4()
	view = view.Mul4(mgl.LookAtV(camera.CameraPos, camera.CameraPos.Add(camera.CameraFront), camera.CameraUp))

	gl.BindVertexArray(vao)

	//also draw the lamp object
	lightCubeShader.Use()
	lightCubeShader.SetMat4("projection", projection)
	lightCubeShader.SetMat4("view", view)
	lightCubeShader.SetVec3("color", mgl.Vec3{1.0, 1.0, 1.0})

	for _, cube := range cubes {
		model := mgl.Ident4()
		model = model.Mul4(mgl.Translate3D(cube.X(), cube.Y(), cube.Z()))
		model = model.Mul4(mgl.Scale3D(1, 1, 1))
		lightCubeShader.SetMat4("model", model)

		gl.BindVertexArray(lightVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

	// Maintenance
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

// Run implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
//func Run(p graphics.Platform, r graphics.Renderer, vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader, texture uint32, specularMap uint32) {
func Run(p graphics.Platform, r graphics.Renderer, vao uint32, lightVao uint32, window *glfw.Window, objectShader *graphics.Shader, lightCubeShader *graphics.Shader) {
	imgui.CurrentIO().SetClipboard(graphics.Clipboard{Platform: p})

	showDemoWindow := false
	showGoDemoWindow := false
	clearColor := [3]float32{0.0, 0.0, 0.0}
	showDebug := false

	//showAnotherWindow := false
	var memoryReadings []float32
	memoryReadings = make([]float32, 2000)

	for !p.ShouldStop() {

		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()
		flags := 0

		// 1. Show a simple window.
		// Tip: if we don't call imgui.Begin()/imgui.End() the widgets automatically appears in a window called "Debug".
		{
			imgui.SetNextWindowPosV(imgui.Vec2{X: 0, Y: 0}, imgui.ConditionFirstUseEver, imgui.Vec2{})
			imgui.SetNextWindowBgAlpha(0.35)
			flags |= imgui.WindowFlagsNoTitleBar
			flags |= imgui.WindowFlagsNoScrollbar
			flags |= imgui.WindowFlagsNoResize
			flags |= imgui.WindowFlagsNoCollapse
			flags |= imgui.WindowFlagsNoMove
			imgui.BeginV("Another window", &showDebug, flags)

			// To display these, you'll need to register a compatible font
			imgui.Text("Telemetry") // Display some text

			//imgui.Checkbox("Demo Window", &showDemoWindow) // Edit bools storing our window open/close state
			//imgui.Checkbox("Go Demo Window", &showGoDemoWindow)

			imgui.PlotLinesV(fmt.Sprintf("Value count %.0fMB", memoryReadings[1999]), memoryReadings, 0, "", math.MaxFloat32, math.MaxFloat32, imgui.Vec2{350, 100})
			//imgui.Checkbox("Another Window", &showAnotherWindow)

			//if imgui.Button("Button") { // Buttons return true when clicked (most widgets return true when edited/activated)
			//	counter++
			//}
			//imgui.SameLine()
			//imgui.Text(fmt.Sprintf("counter = %d", counter))

			imgui.Text(fmt.Sprintf("Application average %.3f ms/frame (%.1f FPS)",
				graphics.MillisPerSecond/imgui.CurrentIO().Framerate(), imgui.CurrentIO().Framerate()))

			imgui.End()
		}

		//// 2. Show another simple window. In most cases you will use an explicit Begin/End pair to name your windows.
		//if showAnotherWindow {
		//	// Pass a pointer to our bool variable (the window will have a closing button that will clear the bool when clicked)
		//	imgui.BeginV("Another window", &showAnotherWindow, 0)
		//	imgui.Text("Hello from another window!")
		//	if imgui.Button("Close Me") {
		//		showAnotherWindow = false
		//	}
		//	imgui.End()
		//}

		//// 3. Show the ImGui demo window. Most of the sample code is in imgui.ShowDemoWindow().
		//// Read its code to learn more about Dear ImGui!
		if showDemoWindow {
			// Normally user code doesn't need/want to call this because positions are saved in .ini file anyway.
			// Here we just want to make the demo initial state a bit more friendly!
			const demoX = 650
			const demoY = 20
			imgui.SetNextWindowPosV(imgui.Vec2{X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})

			imgui.ShowDemoWindow(&showDemoWindow)
		}
		if showGoDemoWindow {
			graphics.Show(&showGoDemoWindow)
		}

		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		draw(vao, lightVao, window, objectShader, lightCubeShader) //, texture, specularMap)

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

		memoryReadings = append(memoryReadings, getMemoryAllocated())
		memoryReadings = memoryReadings[1:]

		// sleep to avoid 100% CPU usage for this demo
		<-time.After(graphics.SleepDuration)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func getMemoryAllocated() float32 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float32(bToMb(m.Alloc))
}
