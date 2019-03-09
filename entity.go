package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Entity represents an entity in the game world
type Entity struct {
	position mgl32.Vec3
	rotation mgl32.Vec3
	scale    float32
	model    *Model
}

// Load loads the uniform variables unique to the current entity. THE SHADER PROGRAM MUST BE ACTIVE!
func (e *Entity) Load(shader *ShaderProgram) {
	translation := mgl32.Translate3D(e.position.X(), e.position.Y(), e.position.Z())
	rotationX := mgl32.HomogRotate3DX(e.rotation.X())
	rotationY := mgl32.HomogRotate3DY(e.rotation.Y())
	rotationZ := mgl32.HomogRotate3DY(e.rotation.Z())
	scale := mgl32.Scale3D(e.scale, e.scale, e.scale)
	modelMatrix := translation.Mul4(rotationX).Mul4(rotationY).Mul4(rotationZ).Mul4(scale)
	shader.LoadUniformMatrix("modelMatrix", modelMatrix)
}
