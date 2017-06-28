package goglad

import (
	"github.com/go-gl/gl/v4.4-core/gl"
	"image"
	"image/draw"
)

type Texture uint32

func NewTexture() Texture {
	var tex uint32
	gl.GenTextures(1, &tex)
	return Texture(tex)
}

func (tex Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(tex))
}

func (tex Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Reasonable defaults, using RGBA format
func (tex Texture) Image2D(img image.Image) {
	// Copy image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	// Get image dimension
	w, h := rgba.Rect.Size().X, rgba.Rect.Size().Y
	// Check if stride is fine
	if rgba.Stride != 4*w {
		panic("Unsupported image stride")
	}
	// Copy image onto RGBA surface
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// Copy data to OpenGL context
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
}

// Create an empty RGBA texture with given w*h
func (tex Texture) Empty2D(w, h int) {
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
}

func (tex Texture) SetFilters(magFilter, minFilter int32) {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, magFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, minFilter)
}
