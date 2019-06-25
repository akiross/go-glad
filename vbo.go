package glad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// VertexBufferObject represents a (vertex) buffer object in OpenGL context
type VertexBufferObject uint32

// NewVertexBufferObject creates a new buffer object
func NewVertexBufferObject() VertexBufferObject {
	var vbo uint32
	gl.CreateBuffers(1, &vbo)
	return VertexBufferObject(vbo)
}

// TODO map buffer

func (vbo VertexBufferObject) Delete() {
}

func (vbo VertexBufferObject) Bind(target uint32) {
	gl.BindBuffer(target, uint32(vbo))
}

func (vbo VertexBufferObject) Unbind(target uint32) {
	gl.BindBuffer(target, 0)
}

func (vbo VertexBufferObject) BufferData32(data []float32, usage uint32) {
	gl.NamedBufferData(uint32(vbo), int(len(data))*4, gl.Ptr(data), usage)
}

func (vbo VertexBufferObject) BufferData64(data []float64, usage uint32) {
	gl.NamedBufferData(uint32(vbo), int(len(data))*8, gl.Ptr(data), usage)
}

/*
func (vbo VertexBufferObject) BufferStorage32(data []float32, usage uint32) {
	gl.NamedBufferStorage(uint32(vbo), int32(len(data))*4, gl.Ptr(data), usage)
}

func (vbo VertexBufferObject) BufferStorage64(data []float64, usage uint32) {
	gl.NamedBufferStorage(uint32(vbo), int32(len(data))*8, gl.Ptr(data), usage)
}
*/

func (vbo VertexBufferObject) BufferSubData32(data []float32, offset int) {
	gl.NamedBufferSubData(uint32(vbo), offset, int(len(data))*4, gl.Ptr(data))
}

func (vbo VertexBufferObject) BufferSubData64(data []float64, offset int) {
	gl.NamedBufferSubData(uint32(vbo), offset, int(len(data))*8, gl.Ptr(data))
}

func (vbo VertexBufferObject) Clear32(data []float32) {
	switch len(data) {
	case 1:
		gl.ClearNamedBufferData(uint32(vbo), gl.R32F, gl.RED, gl.FLOAT, gl.Ptr(data))
	case 2:
		gl.ClearNamedBufferData(uint32(vbo), gl.RG32F, gl.RG, gl.FLOAT, gl.Ptr(data))
	case 3:
		gl.ClearNamedBufferData(uint32(vbo), gl.RGB32F, gl.RGB, gl.FLOAT, gl.Ptr(data))
	case 4:
		gl.ClearNamedBufferData(uint32(vbo), gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.Ptr(data))
	default:
		panic("Clear32 supports only slices of size 1, 2, 3 or 4.")
	}
}

func (vbo VertexBufferObject) CopyTo(dest VertexBufferObject, readOffset, writeOffset int, size int32) {
	gl.CopyNamedBufferSubData(uint32(vbo), uint32(dest), readOffset, writeOffset, int(size))
}
