package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Model represents a model without texture information
type Model struct {
	vao      uint32
	vbos     []uint32
	indices  uint32
	size     int32
	textures []uint32
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
func CreateModelFromData(vertices []float32, indices []uint32, textureCoords []float32) (Model, error) {
	model := Model{vao: 0, vbos: []uint32{0, 0}, size: int32(len(indices)), indices: 0, textures: []uint32{}}
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
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.GenBuffers(1, &model.vbos[1])
	gl.BindBuffer(gl.ARRAY_BUFFER, model.vbos[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(textureCoords)*4, gl.Ptr(textureCoords), gl.STATIC_DRAW)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)
	return model, nil
}

// Draw draws the model to the screen. The shader should be already bound.
func (m Model) Draw() {
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.indices)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	for i, t := range m.textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, t)
	}
	gl.DrawElements(gl.TRIANGLES, m.size, gl.UNSIGNED_INT, nil)
	for i := 0; i < len(m.textures); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
}

// AddTexture adds a texture to a given model. Mipmaps will be created if mipmap is true
func (m *Model) AddTexture(path string, mipmap bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	i, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(i.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), i, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(len(m.textures)))
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Bounds().Size().X), int32(rgba.Bounds().Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	if mipmap {
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	m.textures = append(m.textures, texture)

	return nil
}
