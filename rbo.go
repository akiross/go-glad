package glad

import "github.com/go-gl/gl/v4.5-core/gl"

// RenderbufferObject is, similarly to a texture, a destination when rendering in a FBO.
// RB cannot be used as textures, therefore OpenGL implementation might be optimized.
// Use RBOs when you want to perform off-screen rendering without using the results
// as a texture.
type RenderbufferObject uint32

// RBOParameter represent
type RBOParameter int32

// NewRenderbuffer creates a new renderbuffer object
func NewRenderbuffer() RenderbufferObject {
	var rbo uint32
	gl.CreateRenderbuffers(1, &rbo)
	return RenderbufferObject(rbo)
}

// Bind the RBO to the renderbuffer target
func (rbo RenderbufferObject) Bind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, uint32(rbo))
}

// Unbind the RBO from the renderbuffer target
func (rbo RenderbufferObject) Unbind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
}

// Storage allocates the storage for the RBO
// format can be gl.RGB, RGBA, etc, gl.STENCIL_INDEX or gl.DEPTH_COMPONENT
func (rbo RenderbufferObject) Storage(format uint32, width, height int32) {
	gl.NamedRenderbufferStorage(uint32(rbo), format, width, height)
}

// GetParameter returns the RBO parameter value
func (rbo RenderbufferObject) GetParameter(param uint32) int32 {
	var v int32
	gl.GetNamedRenderbufferParameteriv(uint32(rbo), param, &v)
	return v
}
