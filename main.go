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

var cnf *config.Config

func main() {
	cnf = config.ReadFile("default.json")
	vertexShaderSource := getTextFileContents("resources\\shaders\\vertex\\shader.glsl")
	fragmentShaderSource := getTextFileContents("resources\\shaders\\fragment\\shader.glsl")

	runtime.LockOSThread()

	window := graphics.InitGlfw(cnf)
	defer glfw.Terminate()
	program := graphics.InitOpenGL(vertexShaderSource, fragmentShaderSource)

	vao := graphics.MakeVao(shape.Rectangle, shape.Indices)
	texture1 := graphics.MakeTexture("resources\\textures\\container.jpg")
	texture2 := graphics.MakeTexture("resources\\textures\\awesomeface.png")

	gl.UseProgram(program)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture1\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("texture2\x00")), 1)

	for !window.ShouldClose() {
		draw(vao, window, program, texture1, texture2)
	}
}

// loop over cells and tell them to draw
func draw(vao uint32, window *glfw.Window, program uint32, texture1 *uint32, texture2 *uint32) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *texture1)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, *texture2)

	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	glfw.PollEvents()
	window.SwapBuffers()
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
