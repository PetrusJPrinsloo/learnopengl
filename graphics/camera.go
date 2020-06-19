package graphics

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	CameraPos   mgl.Vec3
	CameraFront mgl.Vec3
	CameraUp    mgl.Vec3

	Yaw         float64
	Pitch       float64
	LastX       float64
	LastY       float64
	Fov         float64
	Sensitivity float64

	FirstMouse bool
}

func GetCamera() Camera {
	c := Camera{}
	c.CameraPos = mgl.Vec3{0, 0.0, 3}
	c.CameraFront = mgl.Vec3{0.0, 0.0, -1.0}
	c.CameraUp = mgl.Vec3{0, 1, 0}

	c.Yaw = -90.0
	c.Pitch = 0.0
	c.LastX = 0.0
	c.LastY = 0.0
	c.Fov = 45.0

	c.Sensitivity = 0.2
	c.FirstMouse = true

	return c
}

func (c *Camera) MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	if c.FirstMouse {
		c.LastX = xpos
		c.LastY = ypos
		c.FirstMouse = false
	}

	xoffset := xpos - c.LastX
	yoffset := c.LastY - ypos
	c.LastX = xpos
	c.LastY = ypos

	xoffset *= c.Sensitivity
	yoffset *= c.Sensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	if c.Pitch > 89.0 {
		c.Pitch = 89.0
	}
	if c.Pitch < -89.0 {
		c.Pitch = -89.0
	}

	front := mgl.Vec3{
		float32(math.Cos(float64(mgl.DegToRad(float32(c.Yaw)))) * math.Cos(float64(mgl.DegToRad(float32(c.Pitch))))),
		float32(math.Sin(float64(mgl.DegToRad(float32(c.Pitch))))),
		float32(math.Sin(float64(mgl.DegToRad(float32(c.Yaw)))) * math.Cos(float64(mgl.DegToRad(float32(c.Pitch))))),
	}

	c.CameraFront = front.Normalize()
}

func (c *Camera) ScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	c.Fov -= yoff
	if c.Fov < 1.0 {
		c.Fov = 1.0
	}
	if c.Fov > 45.0 {
		c.Fov = 45.0
	}
}
