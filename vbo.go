package glad

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// VertexBufferObject represents a (vertex) buffer object in OpenGL context
// This is memory allocated by the OpenGL server used to store data, not limited
// to vertices (see the Bind method for details)
type VertexBufferObject uint32

// NewVertexBufferObject creates a new buffer object
func NewVertexBufferObject() VertexBufferObject {
	var vbo uint32
	gl.CreateBuffers(1, &vbo)
	return VertexBufferObject(vbo)
}

// TODO map buffer

// Delete the VBO freeing its name
func (vbo VertexBufferObject) Delete() {
	v := uint32(vbo)
	gl.DeleteBuffers(1, &v)
}

// Bind the VBO to the specified target
// values for target include:
// - gl.ARRAY_BUFFER when this buffer contains vertex-attribute data to be used
//                   in conjunction with VAO attributes
// - gl.TEXTURE_BUFFER for data bound to texture object: once attached to textures,
//                     the data will be accessible from the shader
// - gl.COPY_READ_BUFFER and gl.COPY_WRITE_BUFFER two buffers that can be used
//                                                to copy data between buffers
// - gl.DRAW_INDIRECT_BUFFER to store parameters when performing indirect drawing
// TODO doc: add and explain other targets
func (vbo VertexBufferObject) Bind(target uint32) {
	gl.BindBuffer(target, uint32(vbo))
}

// Unbind the VBO from the specified target
func (vbo VertexBufferObject) Unbind(target uint32) {
	gl.BindBuffer(target, 0)
}

// BufferStorage allocates a new immutable data store for the VBO
// The storage cannot change in size: to do so, the VBO must be deleted first
// and created again with a new size. Data can be modified using BufferSubData
func (vbo VertexBufferObject) BufferStorage(data []float32, flags uint32) {
	gl.NamedBufferStorage(uint32(vbo), len(data)*4, gl.Ptr(data), flags)
}

// BufferStorage32 allocates the storage for the VBO and copies float32 data in it
func (vbo VertexBufferObject) BufferStorage32(data []float32, flags uint32) {
	gl.NamedBufferStorage(uint32(vbo), len(data)*4, gl.Ptr(data), flags)
}

// BufferStorage64 allocates the storage for the VBO and copies float64 data in it
func (vbo VertexBufferObject) BufferStorage64(data []float64, flags uint32) {
	gl.NamedBufferStorage(uint32(vbo), len(data)*8, gl.Ptr(data), flags)
}

// BufferData32 allocates a new data store for float32 data in the VBO
// usage can be gl.<frequency>_<nature> with
// - frequency: STREAM, STATIC, DYNAMIC
// - nature: DRAW, READ, COPY
// Pre-existing storage will be deleted, therefore size might change
func (vbo VertexBufferObject) BufferData32(data []float32, usage uint32) {
	gl.NamedBufferData(uint32(vbo), len(data)*4, gl.Ptr(data), usage)
}

// BufferData64 is the same of BufferData32 but with float32
func (vbo VertexBufferObject) BufferData64(data []float64, usage uint32) {
	gl.NamedBufferData(uint32(vbo), len(data)*8, gl.Ptr(data), usage)
}

// BufferSubData32 replaces part of the buffer content with new float32 data
func (vbo VertexBufferObject) BufferSubData32(data []float32, offset int) {
	gl.NamedBufferSubData(uint32(vbo), offset, len(data)*4, gl.Ptr(data))
}

// BufferSubData64 replaces part of the buffer content with new float64 data
func (vbo VertexBufferObject) BufferSubData64(data []float64, offset int) {
	gl.NamedBufferSubData(uint32(vbo), offset, len(data)*8, gl.Ptr(data))
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

// TODO function to get back the data glGetNamedBufferSubData
// TODO glMapNamedBuffer and glUnmapNamedBuffer
