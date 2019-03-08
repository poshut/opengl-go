package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// Texture represents an OpenGL 2D texture
type Texture uint32

// Bind binds the texture to the provided numeric texture unit
func (t Texture) Bind(unit int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(unit))
	gl.BindTexture(gl.TEXTURE_2D, uint32(t))
}

// Unbind removes the binding of the texture from the provided numeric texture unit
func (t Texture) Unbind(unit int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(unit))
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// NewTextureFromData creates a texture from the provided io.Reader
func NewTextureFromData(r io.Reader, mipmap bool) (Texture, error) {
	i, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(i.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), i, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
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

	return Texture(texture), nil

}

// NewTextureFromFile creates a texture from the provided path
func NewTextureFromFile(path string, mipmap bool) (Texture, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return NewTextureFromData(file, mipmap)
}
