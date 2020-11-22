package input

import (
	"github.com/PetrusJPrinsloo/learnopengl/graphics"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
)

func ProcessInput(window *glfw.Window, camera *graphics.Camera, deltaTime float64) {
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
