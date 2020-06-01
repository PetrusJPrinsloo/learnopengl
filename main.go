package main

import (
	"github.com/PetrusJPrinsloo/learnopengl/config"
	"github.com/PetrusJPrinsloo/learnopengl/graphics"
	"github.com/PetrusJPrinsloo/learnopengl/shape"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	cnf *config.Config
)

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getShaderFileContents("resources\\shaders\\vertex\\shader.glsl")
	fragmentShaderSource := getShaderFileContents("resources\\shaders\\fragment\\shader.glsl")

	runtime.LockOSThread()

	window := graphics.InitGlfw(cnf)
	defer glfw.Terminate()
	program := graphics.InitOpenGL(vertexShaderSource, fragmentShaderSource)

	vao := graphics.MakeVao(shape.Triangle)
	for !window.ShouldClose() {
		draw(vao, window, program)
	}
}

func getShaderFileContents(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	text := string(content)
	return text
}

// loop over cells and tell them to draw
func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Triangle)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}
