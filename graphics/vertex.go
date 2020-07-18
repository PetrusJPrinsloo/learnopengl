package graphics

import mgl "github.com/go-gl/mathgl/mgl32"

type Vertex struct {
	Position  mgl.Vec3
	Normal    mgl.Vec3
	TexCoords mgl.Vec3
}
