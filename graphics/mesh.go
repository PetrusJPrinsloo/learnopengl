package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strconv"
	"unsafe"
)

type Mesh struct {
	Id       int
	Vertices []Vertex
	Indices  []uint32
	Textures []Texture
	vao      uint32
	vbo, ebo uint32
}

func NewMesh(v []Vertex, i []uint32, t []Texture) Mesh {
	m := Mesh{
		Vertices: v,
		Indices:  i,
		Textures: t,
	}
	//m.setup()
	return m
}

func (m *Mesh) setup() {
	// size of the Vertex struct
	dummy := m.Vertices[0]
	structSize := int(unsafe.Sizeof(dummy))
	structSize32 := int32(structSize)

	// Create buffers/arrays
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	// Load data into vertex buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	// A great thing about structs is that their memory layout is sequential for all its items.
	// The effect is that we can simply pass a pointer to the struct and it translates perfectly to a gl.m::vec3/2 array which
	// again translates to 3/2 floats which translates to a byte array.
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*structSize, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	// Set the vertex attribute pointers
	// Vertex Positions
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, structSize32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Vertex Normals
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, structSize32, unsafe.Pointer((unsafe.Offsetof(dummy.Normal))))
	gl.EnableVertexAttribArray(1)

	// Vertex Texture Coords
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, structSize32, unsafe.Pointer((unsafe.Offsetof(dummy.TexCoords))))
	gl.EnableVertexAttribArray(2)

	// Vertex Tangent
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, structSize32, unsafe.Pointer(unsafe.Offsetof(dummy.Tangent)))

	// Vertex Bitangent
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, structSize32, unsafe.Pointer(unsafe.Offsetof(dummy.Bitangent)))

	gl.BindVertexArray(0)
}

func (m *Mesh) draw(program uint32) {
	// Bind appropriate textures
	var (
		diffuseNr  uint64
		specularNr uint64
		normalNr   uint64
		heightNr   uint64
		i          uint32
	)
	diffuseNr = 1
	specularNr = 1
	normalNr = 1
	heightNr = 1
	i = 0
	for i = 0; i < uint32(len(m.Textures)); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + i) // Active proper texture unit before binding

		// Retrieve texture number (the N in diffuse_textureN)
		ss := ""
		switch m.Textures[i].TextureType {
		case "texture_diffuse":
			ss = ss + strconv.FormatUint(diffuseNr, 10) // Transfer GLuint to stream
			diffuseNr++
		case "texture_specular":
			ss = ss + strconv.FormatUint(specularNr, 10) // Transfer GLuint to stream
			specularNr++
		case "texture_normal":
			ss = ss + strconv.FormatUint(normalNr, 10) // Transfer GLuint to stream
			normalNr++
		case "texture_height":
			ss = ss + strconv.FormatUint(heightNr, 10) // Transfer GLuint to stream
			heightNr++
		}

		// Now set the sampler to the correct texture unit
		tu := m.Textures[i].TextureType + ss + "\x00"

		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str(tu)), int32(i))
		// And finally bind the texture
		gl.BindTexture(gl.TEXTURE_2D, m.Textures[i].id)
	}

	// Draw mesh
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(m.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	// Always good practice to set everything back to defaults once configured.
	for i = 0; i < uint32(len(m.Textures)); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + i)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
}
