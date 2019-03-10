package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	cursorSpeed = 0.005
	moveSpeed   = 1.0
	fov         = 70.0
	nearPlane   = 0.1
	farPlane    = 1000.0
)

// Camera represents the camera in the world
type Camera struct {
	position mgl32.Vec3
	yaw      float64
	pitch    float64

	lastCursorPosX float64
	lastCursorPosY float64
}

// NewCamera initializes a camera and set OpenGL options
func NewCamera(window *glfw.Window) Camera {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	lastX, lastY := window.GetCursorPos()
	return Camera{mgl32.Vec3{0.0, 0.0, 0.0}, -math.Pi, 0.0, lastX, lastY}
}

// Load Loads the viewMatrix into the provided shader program. THE PROGRAM MUST BE ACTIVE!
func (c *Camera) Load(program *ShaderProgram) {
	directionVector := mgl32.Vec3{float32(math.Cos(c.pitch) * math.Sin(c.yaw)), float32(math.Sin(c.pitch)), float32(math.Cos(c.pitch) * math.Cos(c.yaw))}
	viewMatrix := mgl32.LookAtV(c.position, c.position.Add(directionVector), mgl32.Vec3{0.0, 1.0, 0.0})
	program.LoadUniformMatrix("viewMatrix", viewMatrix)
}

// Update updates the camera from the current input situation
func (c *Camera) Update(window *glfw.Window) {

	newX, newY := window.GetCursorPos()
	if c.lastCursorPosX != 0 || c.lastCursorPosY != 0 {
		c.yaw += cursorSpeed * (c.lastCursorPosX - newX)
		c.pitch += cursorSpeed * (c.lastCursorPosY - newY)
		c.pitch = clampf64(c.pitch, -math.Pi/2+0.1, math.Pi/2-0.1)
	}
	c.lastCursorPosX = newX
	c.lastCursorPosY = newY

	if window.GetKey(glfw.KeyD) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{float32(-math.Cos(c.yaw) * moveSpeed), 0, float32(math.Sin(c.yaw) * moveSpeed)})
	} else if window.GetKey(glfw.KeyA) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{float32(math.Cos(c.yaw) * moveSpeed), 0, float32(-math.Sin(c.yaw) * moveSpeed)})
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{float32(math.Sin(c.yaw) * moveSpeed), 0, float32(math.Cos(c.yaw) * moveSpeed)})
	} else if window.GetKey(glfw.KeyS) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{float32(-math.Sin(c.yaw) * moveSpeed), 0, float32(-math.Cos(c.yaw) * moveSpeed)})
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{0.0, moveSpeed, 0.0})
	} else if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		c.position = c.position.Add(mgl32.Vec3{0.0, -moveSpeed, 0.0})
	}
}
