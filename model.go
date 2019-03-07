package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Model represents a model without texture information
type Model struct {
	vao     uint32
	vbos    []uint32
	indices uint32
	size    int32
}

// Delete deletes the model
func (m Model) Delete() {
	gl.DeleteBuffers(int32(len(m.vbos)), &m.vbos[0])
	gl.DeleteBuffers(1, &m.indices)
	gl.DeleteVertexArrays(1, &m.vao)
}

// func CreateModelFromFile(file string) (Model, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return Model{0, nil, 0}, err
// 	}
// 	object, err := obj.NewReader(f).Read()
// 	if err != nil {
// 		return Model{0, nil, 0}, err
// 	}

// 	// add VAO, VBO and index buffer

// }

// CreateModelFromData creates a model from the provided vertex and index data
func CreateModelFromData(vertices []float32, indices []uint32) (Model, error) {
	model := Model{vao: 0, vbos: []uint32{0, 0}, size: int32(len(indices))}
	gl.GenVertexArrays(1, &model.vao)
	gl.BindVertexArray(model.vao)

	gl.GenBuffers(1, &model.vbos[0])
	gl.BindBuffer(gl.ARRAY_BUFFER, model.vbos[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &model.indices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, model.indices)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindVertexArray(0)

	return model, nil
}

// Draw draws the model to the screen. The shader should be already bound.
func (m Model) Draw() {
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.indices)
	gl.EnableVertexAttribArray(0)
	gl.DrawElements(gl.TRIANGLES, m.size, gl.UNSIGNED_INT, nil)

	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
}
