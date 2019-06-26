package glad

import "github.com/go-gl/gl/v4.5-core/gl"

// FramebufferObject represents a framebuffer in the OpenGL context
// A FBO is a collection of buffers that can be used as rendering
// destination for the OpenGL pipeline
// There are different attachment points in a FBO:
// - at least one gl.COLOR_ATTACHMENTn (n=0, 1, 2, ...)
// - gl.DEPTH_ATTACHMENT for the depth buffer
// - gl.STENCIL_ATTACHMENT for the stencil buffer
type FramebufferObject uint32

// NewFramebuffer creates a new FBO
func NewFramebuffer() FramebufferObject {
	var fbo uint32
	gl.CreateFramebuffers(1, &fbo)
	return FramebufferObject(fbo)
}

// Delete the FBO, freeing its name
func (fbo FramebufferObject) Delete() {
	f := uint32(fbo)
	gl.DeleteFramebuffers(1, &f)
}

// Bind the FBO to the framebuffer target, allowing to use it for GL output
func (fbo FramebufferObject) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(fbo))
}

// Unbind the FBO, restoring the default window-system framebuffer
func (fbo FramebufferObject) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// Texture attaches a texture level to the FBO
func (fbo FramebufferObject) Texture(att uint32, texture Texture) {
	gl.NamedFramebufferTexture(uint32(fbo), att, uint32(texture), 0)
}
