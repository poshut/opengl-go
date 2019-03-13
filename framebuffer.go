package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Framebuffer represents an OpenGL framebuffer object
type Framebuffer struct {
	id            uint32
	textures      []Texture
	renderbuffers []uint32
}

// The quad used to display post processing effects
var quadModel Model

// InitializeFramebuffers initializes the quad model. You have to call this before you can use a framebuffer
func InitializeFramebuffers() {

	var quadVertices = []float32{
		-1.0, -1.0,
		1.0, -1.0,
		-1.0, 1.0,
		1.0, 1.0,
	}

	var quadTexCoords = []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	}

	var quadIndices = []uint32{
		0, 1, 2,
		1, 2, 3,
	}
	quadModel = NewModel()
	quadModel.AddBufferAndAttribute3f(quadVertices, 2, false)
	quadModel.AddBufferAndAttribute3f(quadTexCoords, 2, false)
	quadModel.SetIndexBuffer(quadIndices)
}

// NewFramebuffer creates a framebuffer without any attachments
func NewFramebuffer() (Framebuffer, error) {
	fbo := Framebuffer{0, []Texture{}, []uint32{}}
	gl.GenFramebuffers(1, &fbo.id)
	return fbo, nil
}

// AddColorAttachment adds a color attachment that can be used as a texture
func (f *Framebuffer) AddColorAttachment() error {
	texture, err := NewTextureFromData(windowWidth, windowHeight, nil, false)
	if err != nil {
		return err
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	texture.Bind(0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, uint32(gl.COLOR_ATTACHMENT0+len(f.textures)), gl.TEXTURE_2D, uint32(texture), 0)
	texture.Unbind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	f.textures = append(f.textures, texture)

	return nil
}

// AddRenderbufferDepthAndStencil adds a renderbuffer object that is used for the depth and stencil buffer
func (f *Framebuffer) AddRenderbufferDepthAndStencil() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	var id uint32
	gl.GenRenderbuffers(1, &id)
	gl.BindRenderbuffer(gl.RENDERBUFFER, id)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, windowWidth, windowHeight)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, id)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	f.renderbuffers = append(f.renderbuffers, id)
}

// IsComplete checks if enough attachments are present on the framebuffer
func (f *Framebuffer) IsComplete() bool {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	res := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return res == gl.FRAMEBUFFER_COMPLETE
}

// Delete deletes the framebuffer
func (f *Framebuffer) Delete() {
	gl.DeleteFramebuffers(1, &f.id)
}

// Use binds the framebuffer and clears its buffers
func (f *Framebuffer) Use() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Unuse binds the screen framebuffer
func (f *Framebuffer) Unuse() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Draw renders the quad with the added texture attachments using the provided shader program
func (f *Framebuffer) Draw(program *ShaderProgram) {
	if quadModel.vao == 0 {
		panic("fbos not initialized, did you call InitializeFramebuffers()?")
	}
	quadModel.Bind(program)
	for i, t := range f.textures {
		t.Bind(i)
	}
	quadModel.Draw()
	for i, t := range f.textures {
		t.Unbind(i)
	}
	quadModel.Unbind(program)
}
