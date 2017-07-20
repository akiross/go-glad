package goglad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"image"
	"image/draw"
)

type Texture uint32

func NewTexture() Texture {
	var tex uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &tex)
	return Texture(tex)
}

func (tex Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(tex))
}

func (tex Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Create empty texture of given size
func (tex Texture) Storage2D(w, h int) {
	gl.TextureStorage2D(uint32(tex), 1, gl.RGBA8, int32(w), int32(h))
}

// Copy image into pre-allocated texture (call Storage2D before this)
func (tex Texture) Image2D(img image.Image) {
	// Copy image to RGBA format, if necessary
	var rgba *image.RGBA
	switch img.(type) {
	case *image.RGBA:
		rgba = img.(*image.RGBA)
	default:
		rgba = image.NewRGBA(img.Bounds())
	}
	// Get image dimension
	w, h := rgba.Rect.Size().X, rgba.Rect.Size().Y
	// Check if stride is fine
	if rgba.Stride != 4*w {
		panic("Unsupported image stride")
	}
	// Copy image onto RGBA surface
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// Copy data to OpenGL context
	gl.TextureSubImage2D(uint32(tex), 0, 0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
}

func (tex Texture) SetFilters(magFilter, minFilter int32) {
	gl.TextureParameteri(uint32(tex), gl.TEXTURE_MAG_FILTER, magFilter)
	gl.TextureParameteri(uint32(tex), gl.TEXTURE_MIN_FILTER, minFilter)
}

func (tex Texture) Clear(r, g, b, a byte) {
	rgba := []byte{r, g, b, a}
	gl.ClearTexImage(uint32(tex), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba))
}
