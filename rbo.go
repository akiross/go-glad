package goglad

import "github.com/go-gl/gl/v4.4-core/gl"

// Renderbuffer is, similarly to a texture, a destination when rendering in a FBO.
// RB cannot be used as textures, therefore OpenGL implementation might be optimized.
// Use RBOs when you want to perform off-screen rendering without using the results
// as a texture.

type RenderbufferObject uint32

func NewRenderbuffer() RenderbufferObject {
	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	return RenderbufferObject(rbo)
}

func (rbo RenderbufferObject) Bind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, uint32(rbo))
}

func (rbo RenderbufferObject) Unbind() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
}
