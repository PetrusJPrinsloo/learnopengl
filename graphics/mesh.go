package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strconv"
	"unsafe"
)

type Mesh struct {
	Vertices []Vertex
	Indices  []int32
	Textures []Texture

	vao, vbo, ebo uint32
}

func NewMesh(vertices []Vertex, indices []int32, textures []Texture) *Mesh {
	mesh := Mesh{
		Vertices: vertices,
		Indices:  indices,
		Textures: textures,
	}

	mesh.setupMesh()

	return &mesh
}

func (m *Mesh) Draw(shader *Shader) {
	diffuseNr := 1
	specularNr := 1

	for i, tex := range m.Textures {
		var number string
		name := tex.Type
		if name == "texture_diffuse" {
			number = strconv.Itoa(diffuseNr)
			diffuseNr++
		} else if name == "texture_specular" {
			number = strconv.Itoa(specularNr)
			specularNr++
		}

		shader.SetFloat("material."+name+number, float32(i))
		gl.BindTexture(gl.TEXTURE_2D, tex.Id)
	}

	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.Ptr(0))
	gl.BindVertexArray(0)
}

// Sets up the mesh for opengl from the data
func (m *Mesh) setupMesh() {

	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(m.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(m.Vertices)*int(unsafe.Sizeof(Vertex{})),
		gl.Ptr(m.Vertices),
		gl.STATIC_DRAW,
	)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ARRAY_BUFFER,
		len(m.Indices)*int(unsafe.Sizeof(int32(1))),
		gl.Ptr(m.Indices),
		gl.STATIC_DRAW,
	)

	// Vertex Positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		0,
		3,
		gl.FLOAT,
		false,
		int32(unsafe.Sizeof(Vertex{})),
		gl.PtrOffset(0),
	)

	// Vertex Normals
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1,
		3,
		gl.FLOAT,
		false,
		int32(unsafe.Sizeof(Vertex{})),
		gl.PtrOffset(int(unsafe.Offsetof(Vertex{}.Normal))),
	)

	// Vertex texture Coordinates
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2,
		2,
		gl.FLOAT,
		false,
		int32(unsafe.Sizeof(Vertex{})),
		gl.PtrOffset(int(unsafe.Offsetof(Vertex{}.TexCoords))),
	)

	gl.BindVertexArray(0)
}
