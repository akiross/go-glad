package goglad

import (
	"github.com/go-gl/gl/v4.4-core/gl"
)

type VertexBufferObject uint32

func NewVertexBufferObject() VertexBufferObject {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	return VertexBufferObject(vbo)
}

func (vbo VertexBufferObject) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(vbo))
}

func (vbo VertexBufferObject) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vbo VertexBufferObject) BufferData32(data []float32, usage uint32) {
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), usage)
}

func (vbo VertexBufferObject) BufferData64(data []float64, usage uint32) {
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*8, gl.Ptr(data), usage)
}

func (vbo VertexBufferObject) BufferSubData32(data []float32, offset int) {
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(data)*4, gl.Ptr(data))
}

func (vbo VertexBufferObject) BufferSubData64(data []float64, offset int) {
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, len(data)*8, gl.Ptr(data))
}
