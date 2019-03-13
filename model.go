package main

import (
	"fmt"
	_ "image/jpeg"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Model represents a model without texture information that has an index buffer
type Model struct {
	vao      uint32
	vbos     []uint32
	indices  uint32
	size     int32
	textures []Texture
}

// Delete deletes the model
func (m *Model) Delete() {
	gl.DeleteBuffers(int32(len(m.vbos)), &m.vbos[0])
	gl.DeleteBuffers(1, &m.indices)
	gl.DeleteVertexArrays(1, &m.vao)
}

// CreateModelFromFile loads an .obj file into a model
func CreateModelFromFile(file string) (Model, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return Model{0, nil, 0, 0, nil}, err
	}
	lines := strings.Split(string(fileData), "\n")
	vertices := []float32{}
	indices := []uint32{}
	textureCoords := []float32{}
	normals := []float32{}

	realTextureCoords := []float32{}
	realNormals := []float32{}

	for _, l := range lines {
		l = strings.Trim(l, "\r")
		if strings.HasPrefix(l, "v ") {
			lineParts := strings.Split(l, " ")
			if len(lineParts) != 4 {
				return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
			}
			for _, n := range lineParts[1:] {
				f, err := strconv.ParseFloat(n, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, err
				}
				vertices = append(vertices, float32(f))

			}
		} else if strings.HasPrefix(l, "vt ") {
			lineParts := strings.Split(l, " ")
			if len(lineParts) != 3 {
				return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
			}
			for _, n := range lineParts[1:] {
				f, err := strconv.ParseFloat(n, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, err
				}
				textureCoords = append(textureCoords, float32(f))
			}
		} else if strings.HasPrefix(l, "vn ") {
			lineParts := strings.Split(l, " ")
			if len(lineParts) != 4 {
				return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
			}
			for _, n := range lineParts[1:] {
				f, err := strconv.ParseFloat(n, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, err
				}
				normals = append(normals, float32(f))
			}
		} else if strings.HasPrefix(l, "f ") {
			if len(realNormals) == 0 && len(realTextureCoords) == 0 {
				realNormals = make([]float32, len(vertices))
				realTextureCoords = make([]float32, 2*len(vertices)/3)
			}
			lineParts := strings.Split(l, " ")
			if len(lineParts) != 4 {
				return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: Does not have four components %s", file, l)
			}

			for _, p := range lineParts[1:] {
				vertexData := strings.Split(p, "/")
				if len(vertexData) != 3 {
					return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
				}
				vertexIndex, err := strconv.ParseUint(vertexData[0], 10, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
				}
				vertexIndex--
				texCoordIndex, err := strconv.ParseInt(vertexData[1], 10, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
				}
				texCoordIndex--
				normalIndex, err := strconv.ParseInt(vertexData[2], 10, 32)
				if err != nil {
					return Model{0, nil, 0, 0, nil}, fmt.Errorf("Invalid line in %s: %s", file, l)
				}
				normalIndex--
				indices = append(indices, uint32(vertexIndex))
				realNormals[vertexIndex*3] = normals[normalIndex*3]
				realNormals[vertexIndex*3+1] = normals[normalIndex*3+1]
				realNormals[vertexIndex*3+2] = normals[normalIndex*3+2]
				realTextureCoords[vertexIndex*2] = textureCoords[texCoordIndex*2]
				realTextureCoords[vertexIndex*2+1] = 1 - textureCoords[texCoordIndex*2+1]
			}
		}
	}
	return CreateModelFromData(vertices, indices, realTextureCoords, realNormals)
}

// CreateModelFromData creates a model from the provided vertex and index data
func CreateModelFromData(vertices []float32, indices []uint32, textureCoords []float32, normals []float32) (Model, error) {
	model := NewModel()
	model.AddBufferAndAttribute3f(vertices, 3, false)
	model.AddBufferAndAttribute3f(textureCoords, 2, false)
	model.AddBufferAndAttribute3f(normals, 3, true)
	model.SetIndexBuffer(indices)
	return model, nil
}

// NewModel creates a model with a VAO without any buffers
func NewModel() Model {
	model := Model{vao: 0, vbos: []uint32{}, size: 0, indices: 0, textures: []Texture{}}
	gl.GenVertexArrays(1, &model.vao)
	return model
}

// AddBufferAndAttribute3f adds a buffer containing the provided data and a corresponding vertex attribute
func (m *Model) AddBufferAndAttribute3f(data []float32, numComponents int32, normalize bool) {
	gl.BindVertexArray(m.vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
	gl.VertexAttribPointer(uint32(len(m.vbos)), numComponents, gl.FLOAT, normalize, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	m.vbos = append(m.vbos, vbo)
}

// SetIndexBuffer sets the index buffer of the model
func (m *Model) SetIndexBuffer(indices []uint32) {
	gl.BindVertexArray(m.vao)
	gl.GenBuffers(1, &m.indices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.indices)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	m.size = int32(len(indices))

}

// Bind binds all VAOs, vertex attributes textures and model uniforms
func (m *Model) Bind(shader *ShaderProgram) {
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.indices)
	for i := range m.vbos {
		gl.EnableVertexAttribArray(uint32(i))
	}
	for i, t := range m.textures {
		t.Bind(i)
	}
	if len(m.textures) == 0 {
		shader.LoadUniformFloat("hasTexture", 0.0)
	} else {
		shader.LoadUniformFloat("hasTexture", 1.0)
	}
}

// Draw draws the model to the screen. The shader should be already bound.
func (m *Model) Draw() {
	gl.DrawElements(gl.TRIANGLES, m.size, gl.UNSIGNED_INT, nil)
}

// Unbind unbinds all model attributes and textures
func (m *Model) Unbind(shader *ShaderProgram) {
	for i, t := range m.textures {
		t.Unbind(i)
	}
	for i := range m.vbos {
		gl.DisableVertexAttribArray(uint32(i))
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

// AddTexture adds a texture to a given model. Mipmaps will be created if mipmap is true
func (m *Model) AddTexture(path string, mipmap bool) error {
	texture, err := NewTextureFromFile(path, mipmap)
	if err != nil {
		return err
	}
	m.textures = append(m.textures, texture)
	return nil
}
