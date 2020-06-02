package graphics

import (
	"fmt"
	"github.com/PetrusJPrinsloo/learnopengl/config"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
	_ "image/jpeg"
	"log"
	"neilpa.me/go-stbi"
	"os"
	"strings"
	"unsafe"
)

func MakeVao(vertices []float32, indices []uint32) uint32 {
	var (
		vbo uint32
		vao uint32
		ebo uint32
	)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	var offset int = 0
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(unsafe.Sizeof(float32(0))*8), nil)
	gl.EnableVertexAttribArray(0)

	offset += 3
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(unsafe.Sizeof(float32(0))*8), gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)

	offset += 3
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, int32(unsafe.Sizeof(float32(0))*8), gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)

	// note that this is allowed, the call to glVertexAttribPointer registered VBO as the vertex attribute's bound vertex buffer object so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 2)

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to glBindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	gl.BindVertexArray(0)

	return vao
}

func MakeTexture(path string) *uint32{
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// set the texture wrapping parameters
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// set texture filtering parameters
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	width, height := getImageDimension(path)
	data, err := stbi.Load(path)
	if err != nil {
		panic(err)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(data.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return &texture
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return img.Width, img.Height
}

// initGlfw initializes glfw and returns a Window to use.
func InitGlfw(cnf *config.Config) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(cnf.Width, cnf.Height, "Learn OpenGL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL(vertexShaderSource string, fragmentShaderSource string) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
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
