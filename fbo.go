package glad

import "github.com/go-gl/gl/v4.5-core/gl"

type FramebufferObject uint32

func NewFramebuffer() FramebufferObject {
	var fbo uint32
	gl.CreateFramebuffers(1, &fbo)
	return FramebufferObject(fbo)
}

func (fbo FramebufferObject) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, uint32(fbo))
}

func (fbo FramebufferObject) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (fbo FramebufferObject) Texture(att uint32, texture Texture) {
	gl.NamedFramebufferTexture(uint32(fbo), att, uint32(texture), 0)
}
