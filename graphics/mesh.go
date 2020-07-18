package graphics

type Mesh struct {
	// public
	Vertices []Vertex
	Indices  []int32
	Textures []Texture

	// private
	vao, vbo, ebo int32
}

// public
func NewMesh(vertices []Vertex, indices []int32, texture []Texture) *Mesh {
	mesh := Mesh{}

	return &mesh
}

func (m *Mesh) Draw(shader Shader) {

}

// private
func (m *Mesh) setupMesh() {

}
