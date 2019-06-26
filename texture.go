package glad

import (
	"image"
	"image/draw"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// Texture represents a texture in the OpenGL context
// From the OpenGL Guide Book: "Textures may be read by a shader by associating a sampler variable with a texture unit and using GLSL’s built-in functions to fetch texels from the texture’s image. The way in which the texels are fetched depends on a number of parameters that are contained in another object called a sampler object. Sampler objects are bound to sampler units much as texture objects are bound to texture units. For convenience, a texture object may be considered to contain a built-in sampler object of its own that will be used by default to read from it, if no sampler object is bound to the corresponding sampler unit."Textures may be read by a shader by associating a sampler variable with a texture unit and using GLSL’s built-in functions to fetch texels from the texture’s image. The way in which the texels are fetched depends on a number of parameters that are contained in another object called a sampler object. Sampler objects are bound to sampler units much as texture objects are bound to texture units. For convenience, a texture object may be considered to contain a built-in sampler object of its own that will be used by default to read from it, if no sampler object is bound to the corresponding sampler unit."
type Texture uint32

// NewTexture creates a texture of given type with no attached storage
// For example, target could be gl.TEXTURE_2D. The target roughly specifies the
// texture dimensionality, which is tied to the dimensionality of the sampler.
// After creating the texture, you might want to attach storage to it
func NewTexture(target uint32) Texture {
	var tex uint32
	gl.CreateTextures(target, 1, &tex)
	return Texture(tex)
}

// Delete the texture freeing its name, freeing the associated storage
func (tex Texture) Delete() {
	t := uint32(tex)
	gl.DeleteTextures(1, &t)
}

// Bind the texture to the specified texture unit
// This needs to be called before the texture can be accessed by shaders
// There are many texture units in the context that can be used to bind
// different textures: bind a texture to a unit and a sampler to the same unit
// to access the texture data from the shader
func (tex Texture) Bind(unit uint32) {
	gl.BindTextureUnit(unit, uint32(tex))
}

// Unbind the texture from the texture unit
func (tex Texture) Unbind(unit uint32) {
	gl.BindTextureUnit(unit, 0)
}

// Storage allocates storage for an empty texture of given size (cast to int32)
// format can be, for instance, gl.RGBA8
func (tex Texture) Storage(levels int32, internalFmt uint32, size []int) {
	switch len(size) {
	case 1:
		gl.TextureStorage1D(uint32(tex), levels, internalFmt, int32(size[0]))
	case 2:
		gl.TextureStorage2D(uint32(tex), levels, internalFmt, int32(size[0]), int32(size[1]))
	case 3:
		gl.TextureStorage3D(uint32(tex), levels, internalFmt, int32(size[0]), int32(size[1]), int32(size[2]))
	default:
		log.Fatalln("Texture Storage must have size of length 1, 2 or 3")
	}
}

// SubImage replaces a region of the texture with the data
// 1D, 2D or 3D depends on the len of offset and size (they must be equal)
// pixels is considered an offset to buffer start or pointer to host memory
// depending whether a buffer object is or is not bound to the PIXEL_UNPACK_BUFFER
// target
func (tex Texture) SubImage(level int32, offset, size []int, externalFmt, typ uint32, pixels unsafe.Pointer) {
	if len(offset) != len(size) {
		log.Fatalln("Texture SubImage offset and size must have the same length")
	}
	switch len(offset) {
	case 1:
		gl.TextureSubImage1D(uint32(tex), level, int32(offset[0]), int32(size[0]), externalFmt, typ, pixels)
	case 2:
		gl.TextureSubImage2D(uint32(tex), level, int32(offset[0]), int32(offset[1]), int32(size[0]), int32(size[1]), externalFmt, typ, pixels)
	case 3:
		gl.TextureSubImage3D(uint32(tex), level, int32(offset[0]), int32(offset[1]), int32(offset[2]), int32(size[0]), int32(size[1]), int32(size[2]), externalFmt, typ, pixels)
	default:
		log.Fatalln("Texture SubImage offset and size must have length equal to 1, 2 or 3")
	}
}

// Image2D copies image into pre-allocated texture
// Call Storage with 2D size before this
// The passed image will be copied and can be discarded after the call
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

// GetImage copies texture data to host
func (tex Texture) GetImage(level int32, fmt, typ uint32, size int32, pixels unsafe.Pointer) {
	gl.GetTextureImage(uint32(tex), level, fmt, typ, size, pixels)
}

func (tex Texture) SetFilters(magFilter, minFilter int32) {
	gl.TextureParameteri(uint32(tex), gl.TEXTURE_MAG_FILTER, magFilter)
	gl.TextureParameteri(uint32(tex), gl.TEXTURE_MIN_FILTER, minFilter)
}

func (tex Texture) Clear(r, g, b, a byte) {
	rgba := []byte{r, g, b, a}
	gl.ClearTexImage(uint32(tex), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba))
}
