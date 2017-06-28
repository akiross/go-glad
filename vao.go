package goglad

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

type VertexArrayObject uint32

func NewVertexArrayObject() VertexArrayObject {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return VertexArrayObject(vao)
}

func (vao VertexArrayObject) Bind() {
	gl.BindVertexArray(uint32(vao))
}

func (vao VertexArrayObject) Unbind() {
	gl.BindVertexArray(0)
}
