package goglad

import "github.com/go-gl/gl/v4.4-core/gl"

type FramebufferObject uint32

func NewFramebuffer() FramebufferObject {
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	return FramebufferObject(fbo)
}

func (fbo FramebufferObject) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(fbo))
}

func (fbo FramebufferObject) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo FramebufferObject) Texture(att uint32, texture Texture) {
	gl.FramebufferTexture(gl.FRAMEBUFFER, att, uint32(texture), 0)
}
