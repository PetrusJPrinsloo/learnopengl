package graphics

import mgl "github.com/go-gl/mathgl/mgl32"

type Vertex struct {
	Position  mgl.Vec3
	Normal    mgl.Vec3
	TexCoords mgl.Vec2
	Tangent   mgl.Vec3
	Bitangent mgl.Vec3
}
